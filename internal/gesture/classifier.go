package gesture

import (
	"fmt"

	"showroom/internal/model"
)

type Action string

const (
	ActionWelcome Action = "welcome"
	ActionTour    Action = "tour"
	ActionHeart   Action = "heart"
	ActionRestore Action = "restore"
	ActionIgnore  Action = "ignore"
)

type Decision struct {
	Action   Action
	Accepted bool
	Reason   string
}

func Classify(signal Signal) Decision {
	if !Supported(signal.Kind) {
		return Decision{Action: ActionIgnore, Reason: "unsupported gesture"}
	}
	if signal.Strength < 35 {
		return Decision{Action: ActionIgnore, Reason: "signal too weak"}
	}
	switch signal.Kind {
	case "open_palm":
		return Decision{Action: ActionWelcome, Accepted: true, Reason: "open palm selects welcome"}
	case "point":
		return Decision{Action: ActionTour, Accepted: true, Reason: "point selects tour"}
	case "fist":
		return Decision{Action: ActionHeart, Accepted: true, Reason: "fist forms heart"}
	case "wave":
		return Decision{Action: ActionRestore, Accepted: true, Reason: "wave restores background"}
	default:
		return Decision{Action: ActionIgnore, Reason: "no mapped action"}
	}
}

func ToEvent(signal Signal, sessionID string, accepted bool) (model.GestureEvent, error) {
	event := model.GestureEvent{ID: fmt.Sprintf("gesture-%d", signal.Frame), SessionID: sessionID, Kind: signal.Kind, Strength: signal.Strength, Sequence: signal.Frame, Accepted: accepted, CreatedAt: fmt.Sprintf("frame-%d", signal.Frame)}
	if err := event.Validate(); err != nil {
		return model.GestureEvent{}, err
	}
	return event, nil
}
