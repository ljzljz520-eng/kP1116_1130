package session

import (
	"fmt"
	"strings"

	"showroom/internal/model"
)

type Gate struct {
	MinimumStrength int
	AllowedModes    map[model.SceneMode]bool
}

func DefaultGate() Gate {
	return Gate{MinimumStrength: 35, AllowedModes: map[model.SceneMode]bool{model.ModeWelcome: true, model.ModeTour: true, model.ModeIdle: true}}
}

func (g Gate) Check(session model.Session, signalKind string, strength int) error {
	if !session.Active {
		return fmt.Errorf("session is closed")
	}
	if !g.AllowedModes[session.Mode] {
		return fmt.Errorf("mode %s cannot accept gestures", session.Mode)
	}
	if strings.TrimSpace(signalKind) == "" {
		return fmt.Errorf("gesture kind is required")
	}
	if strength < g.MinimumStrength {
		return fmt.Errorf("gesture strength is below gate")
	}
	return nil
}

func Transition(session model.Session, mode model.SceneMode) (model.Session, error) {
	if !session.Active {
		return model.Session{}, fmt.Errorf("closed session cannot transition")
	}
	if mode == model.ModeClosing {
		return model.Session{}, fmt.Errorf("closing requires close operation")
	}
	if mode != model.ModeWelcome && mode != model.ModeTour && mode != model.ModeIdle {
		return model.Session{}, fmt.Errorf("unsupported transition")
	}
	session.Mode = mode
	return session, nil
}
