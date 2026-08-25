package model

import (
	"errors"
	"strings"
)

type SceneMode string

const (
	ModeWelcome SceneMode = "welcome"
	ModeTour    SceneMode = "tour"
	ModeClosing SceneMode = "closing"
	ModeIdle    SceneMode = "idle"
)

type Phrase struct {
	ID        string
	Text      string
	Mode      SceneMode
	Priority  int
	Enabled   bool
	CreatedAt string
}

type Session struct {
	ID          string
	VisitorName string
	Mode        SceneMode
	StartedAt   string
	EndedAt     string
	Active      bool
}

type GestureEvent struct {
	ID        string
	SessionID string
	Kind      string
	Strength  int
	Sequence  int
	Accepted  bool
	CreatedAt string
}

type DisplayState struct {
	SessionID    string
	Mode         SceneMode
	PhraseID     string
	PhraseText   string
	ParticleForm string
	Revision     int
	UpdatedAt    string
}

type AuditEntry struct {
	ID        string
	SessionID string
	Action    string
	Detail    string
	CreatedAt string
}

func (p Phrase) Validate() error {
	if strings.TrimSpace(p.ID) == "" {
		return errors.New("phrase id is required")
	}
	if strings.TrimSpace(p.Text) == "" {
		return errors.New("phrase text is required")
	}
	if p.Mode != ModeWelcome && p.Mode != ModeTour && p.Mode != ModeClosing && p.Mode != ModeIdle {
		return errors.New("unknown phrase mode")
	}
	if p.Priority < 0 {
		return errors.New("phrase priority cannot be negative")
	}
	return nil
}

func (s Session) Validate() error {
	if strings.TrimSpace(s.ID) == "" {
		return errors.New("session id is required")
	}
	if s.Mode != ModeWelcome && s.Mode != ModeTour && s.Mode != ModeClosing && s.Mode != ModeIdle {
		return errors.New("unknown session mode")
	}
	if s.Active && strings.TrimSpace(s.EndedAt) != "" {
		return errors.New("active session cannot have end time")
	}
	return nil
}

func (g GestureEvent) Validate() error {
	if strings.TrimSpace(g.ID) == "" || strings.TrimSpace(g.SessionID) == "" {
		return errors.New("gesture identity is required")
	}
	if g.Strength < 0 || g.Strength > 100 {
		return errors.New("gesture strength must be between 0 and 100")
	}
	if g.Sequence < 1 {
		return errors.New("gesture sequence must be positive")
	}
	return nil
}

func (d DisplayState) Validate() error {
	if strings.TrimSpace(d.SessionID) == "" {
		return errors.New("display session is required")
	}
	if d.Revision < 1 {
		return errors.New("display revision must be positive")
	}
	if d.ParticleForm == "" {
		return errors.New("particle form is required")
	}
	return nil
}
