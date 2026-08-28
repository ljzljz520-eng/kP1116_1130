package display

import (
	"encoding/json"
	"fmt"
	"strings"
)

type PanelMode string

const (
	PanelVisible PanelMode = "visible"
	PanelHidden  PanelMode = "hidden"
)

type PanelAction string

const (
	PanelShow         PanelAction = "show"
	PanelHide         PanelAction = "hide"
	PanelToggle       PanelAction = "toggle"
	PanelFocusWelcome PanelAction = "focus_welcome"
	PanelFocusTour    PanelAction = "focus_tour"
	PanelFocusHeart   PanelAction = "focus_heart"
	PanelFocusRestore PanelAction = "focus_restore"
	PanelClearFocus   PanelAction = "clear_focus"
)

type PanelState struct {
	Mode     PanelMode
	Focus    string
	Revision int
}

type PanelSnapshot struct {
	Mode         PanelMode `json:"mode"`
	Visible      bool      `json:"visible"`
	Focus        string    `json:"focus"`
	Revision     int       `json:"revision"`
	ActionLabels []string  `json:"actions"`
}

func NewPanelState(hidden bool) PanelState {
	state := PanelState{Mode: PanelVisible, Revision: 1}
	if hidden {
		state.Mode = PanelHidden
	}
	return state
}

func (p PanelState) Validate() error {
	if p.Mode != PanelVisible && p.Mode != PanelHidden {
		return fmt.Errorf("panel mode %q is invalid", p.Mode)
	}
	if p.Revision < 1 {
		return fmt.Errorf("panel revision must be positive")
	}
	if p.Focus != "" && !validPanelFocus(p.Focus) {
		return fmt.Errorf("panel focus %q is invalid", p.Focus)
	}
	return nil
}

func (p *PanelState) Apply(action PanelAction) error {
	if p == nil {
		return fmt.Errorf("panel state is unavailable")
	}
	action = normalizePanelAction(action)
	if !isPanelAction(action) {
		return fmt.Errorf("panel action %q is unsupported", action)
	}
	switch action {
	case PanelShow:
		p.Mode = PanelVisible
	case PanelHide:
		p.Mode = PanelHidden
	case PanelToggle:
		if p.Mode == PanelHidden {
			p.Mode = PanelVisible
		} else {
			p.Mode = PanelHidden
		}
	case PanelFocusWelcome:
		p.Focus = "welcome"
	case PanelFocusTour:
		p.Focus = "tour"
	case PanelFocusHeart:
		p.Focus = "heart"
	case PanelFocusRestore:
		p.Focus = "restore"
	case PanelClearFocus:
		p.Focus = ""
	}
	p.Revision++
	return nil
}

func (p *PanelState) SetHidden(hidden bool) {
	if p == nil {
		return
	}
	if hidden {
		p.Mode = PanelHidden
	} else {
		p.Mode = PanelVisible
	}
	p.Revision++
}

func (p *PanelState) Toggle() {
	if p == nil {
		return
	}
	_ = p.Apply(PanelToggle)
}

func (p *PanelState) SetFocus(focus string) error {
	if p == nil {
		return fmt.Errorf("panel state is unavailable")
	}
	focus = strings.ToLower(strings.TrimSpace(focus))
	if focus != "" && !validPanelFocus(focus) {
		return fmt.Errorf("panel focus %q is invalid", focus)
	}
	if p.Focus != focus {
		p.Focus = focus
		p.Revision++
	}
	return nil
}

func (p PanelState) IsHidden() bool {
	return p.Mode == PanelHidden
}

func (p PanelState) IsVisible() bool {
	return p.Mode == PanelVisible
}

func (p PanelState) FocusedAction() string {
	if p.Focus == "" {
		return "none"
	}
	return p.Focus
}

func (p PanelState) Allows(action string) bool {
	action = strings.ToLower(strings.TrimSpace(action))
	if action == "" {
		return false
	}
	if action == "show" || action == "hide" || action == "toggle" {
		return true
	}
	return validPanelFocus(action)
}

func (p PanelState) Snapshot() PanelSnapshot {
	return PanelSnapshot{
		Mode:         p.Mode,
		Visible:      p.IsVisible(),
		Focus:        p.FocusedAction(),
		Revision:     p.Revision,
		ActionLabels: PanelActionLabels(),
	}
}

func PanelActions() []PanelAction {
	return []PanelAction{PanelShow, PanelHide, PanelToggle, PanelFocusWelcome, PanelFocusTour, PanelFocusHeart, PanelFocusRestore, PanelClearFocus}
}

func PanelActionLabels() []string {
	labels := make([]string, 0, len(PanelActions()))
	for _, action := range PanelActions() {
		labels = append(labels, panelActionLabel(action))
	}
	return labels
}

func (s PanelSnapshot) Validate() error {
	if s.Mode != PanelVisible && s.Mode != PanelHidden {
		return fmt.Errorf("snapshot mode %q is invalid", s.Mode)
	}
	if s.Visible != (s.Mode == PanelVisible) {
		return fmt.Errorf("snapshot visibility does not match mode")
	}
	if s.Revision < 1 {
		return fmt.Errorf("snapshot revision must be positive")
	}
	if s.Focus != "none" && !validPanelFocus(s.Focus) {
		return fmt.Errorf("snapshot focus %q is invalid", s.Focus)
	}
	return nil
}

func (s PanelSnapshot) Encode() ([]byte, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(s)
}

func DecodePanelSnapshot(data []byte) (PanelSnapshot, error) {
	var snapshot PanelSnapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return PanelSnapshot{}, fmt.Errorf("decode panel snapshot: %w", err)
	}
	if err := snapshot.Validate(); err != nil {
		return PanelSnapshot{}, err
	}
	return snapshot, nil
}

func (s PanelSnapshot) Summary() string {
	return fmt.Sprintf("%s:%s:%d", s.Mode, s.Focus, s.Revision)
}

func normalizePanelAction(action PanelAction) PanelAction {
	normalized := strings.ToLower(strings.TrimSpace(string(action)))
	return PanelAction(normalized)
}

func isPanelAction(action PanelAction) bool {
	for _, candidate := range PanelActions() {
		if action == candidate {
			return true
		}
	}
	return false
}

func validPanelFocus(focus string) bool {
	switch strings.ToLower(strings.TrimSpace(focus)) {
	case "welcome", "tour", "heart", "restore":
		return true
	default:
		return false
	}
}

func panelActionLabel(action PanelAction) string {
	switch action {
	case PanelShow:
		return "Show controls"
	case PanelHide:
		return "Hide controls"
	case PanelToggle:
		return "Toggle controls"
	case PanelFocusWelcome:
		return "Welcome phrase"
	case PanelFocusTour:
		return "Tour prompt"
	case PanelFocusHeart:
		return "Heart gesture"
	case PanelFocusRestore:
		return "Restore gesture"
	case PanelClearFocus:
		return "Clear focus"
	default:
		return "Unknown action"
	}
}
