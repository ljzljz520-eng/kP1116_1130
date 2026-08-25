package workflow

import (
	"context"
	"fmt"
	"strings"

	"showroom/internal/audit"
	"showroom/internal/display"
	"showroom/internal/gesture"
	"showroom/internal/model"
	"showroom/internal/persistence"
	"showroom/internal/scheduler"
	"showroom/internal/session"
	"showroom/internal/welcome"
)

type Orchestrator struct {
	store      *persistence.Store
	welcome    *welcome.Service
	selector   *welcome.Selector
	display    *display.Controller
	audit      *audit.Logger
	sessions   *session.Manager
	plan       *scheduler.Plan
	sequence   int
	defaultEnd model.Phrase
}

func New(store *persistence.Store, welcomeService *welcome.Service, controller *display.Controller, logger *audit.Logger) *Orchestrator {
	return &Orchestrator{store: store, welcome: welcomeService, selector: welcome.NewSelector(welcomeService), display: controller, audit: logger, sessions: session.NewManager(store, nil), plan: scheduler.DefaultPlan()}
}

func (o *Orchestrator) StartSession(ctx context.Context, session model.Session) error {
	if err := session.Validate(); err != nil {
		return err
	}
	if _, err := o.sessions.Open(ctx, session.ID, session.VisitorName, session.Mode); err != nil {
		return err
	}
	_, err := o.audit.Record(ctx, session.ID, "session_started", session.VisitorName)
	return err
}

func (o *Orchestrator) HandleSignal(ctx context.Context, sessionID string, signal gesture.Signal) (gesture.Decision, model.DisplayState, error) {
	decision := gesture.Classify(signal)
	event, err := gesture.ToEvent(signal, sessionID, decision.Accepted)
	if err != nil {
		return decision, model.DisplayState{}, err
	}
	if err := o.store.SaveGesture(ctx, event); err != nil {
		return decision, model.DisplayState{}, err
	}
	if !decision.Accepted {
		return decision, o.display.State(), nil
	}
	o.sequence++
	state, err := o.applyDecision(ctx, sessionID, decision)
	if err != nil {
		return decision, model.DisplayState{}, err
	}
	_, err = o.audit.Record(ctx, sessionID, string(decision.Action), decision.Reason)
	return decision, state, err
}

func (o *Orchestrator) applyDecision(ctx context.Context, sessionID string, decision gesture.Decision) (model.DisplayState, error) {
	switch decision.Action {
	case gesture.ActionHeart:
		return o.display.ApplyGesture(ctx, sessionID, "heart")
	case gesture.ActionRestore:
		return o.display.ApplyGesture(ctx, sessionID, "restore")
	case gesture.ActionWelcome:
		phrase, err := o.selector.Choose(ctx, model.ModeWelcome, "welcome-default")
		if err != nil {
			return model.DisplayState{}, err
		}
		return o.display.ShowPhrase(ctx, sessionID, phrase)
	case gesture.ActionTour:
		phrase, err := o.selector.Choose(ctx, model.ModeTour, "tour-default")
		if err != nil {
			return model.DisplayState{}, err
		}
		return o.display.ShowPhrase(ctx, sessionID, phrase)
	default:
		return model.DisplayState{}, fmt.Errorf("unsupported action %q", decision.Action)
	}
}

func (o *Orchestrator) ConfirmCustom(ctx context.Context, sessionID string, mode model.SceneMode, id, text string) (model.DisplayState, error) {
	phrase, err := o.selector.ChooseText(ctx, mode, id, text)
	if err != nil {
		return model.DisplayState{}, err
	}
	return o.display.ShowPhrase(ctx, sessionID, phrase)
}

func (o *Orchestrator) EndDisplay(ctx context.Context, sessionID string) (model.DisplayState, error) {
	phrase, err := o.welcome.EndDisplay(ctx)
	if err != nil {
		return model.DisplayState{}, err
	}
	o.defaultEnd = phrase
	state, err := o.display.ShowPhrase(ctx, sessionID, phrase)
	if err != nil {
		return model.DisplayState{}, err
	}
	_, err = o.audit.Record(ctx, sessionID, "display_ended", phrase.Text)
	return state, err
}

func (o *Orchestrator) LastEndPhrase() model.Phrase {
	return o.defaultEnd
}

func (o *Orchestrator) Tick(ctx context.Context, step int) (model.DisplayState, error) {
	return o.display.Tick(ctx, step)
}

func (o *Orchestrator) ApplyPanelAction(ctx context.Context, sessionID string, action display.PanelAction) (display.PanelSnapshot, error) {
	snapshot, err := o.display.ApplyPanelAction(action)
	if err != nil {
		return display.PanelSnapshot{}, err
	}
	if strings.TrimSpace(sessionID) != "" {
		_, err = o.audit.Record(ctx, sessionID, "panel_"+string(action), snapshot.Summary())
	}
	return snapshot, err
}

func (o *Orchestrator) Panel() display.PanelSnapshot {
	return o.display.Panel()
}

func (o *Orchestrator) NextCue() (scheduler.Cue, bool) {
	return o.plan.Next()
}
