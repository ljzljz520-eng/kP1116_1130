package scheduler

import (
	"fmt"
	"strings"
)

type Frame struct {
	Sequence int
	Cue      Cue
	Label    string
	Final    bool
}

func MakeFrame(sequence int, cue Cue, final bool) (Frame, error) {
	if sequence < 1 {
		return Frame{}, fmt.Errorf("frame sequence must be positive")
	}
	if cue.Name == "" || cue.Action == "" {
		return Frame{}, fmt.Errorf("frame cue is incomplete")
	}
	label := strings.ToUpper(cue.Name)
	if final {
		label = "FINAL:" + label
	}
	return Frame{Sequence: sequence, Cue: cue, Label: label, Final: final}, nil
}

func IsDisplayAction(action string) bool {
	switch action {
	case "welcome", "tour", "closing", "drift", "heart", "restore":
		return true
	default:
		return false
	}
}

func ValidateFrame(frame Frame) error {
	if frame.Sequence < 1 {
		return fmt.Errorf("frame sequence is invalid")
	}
	if !IsDisplayAction(frame.Cue.Action) {
		return fmt.Errorf("frame action %s is not displayable", frame.Cue.Action)
	}
	if frame.Final && !strings.HasPrefix(frame.Label, "FINAL:") {
		return fmt.Errorf("final frame label is malformed")
	}
	return nil
}

func Timeline(plan *Plan) ([]Frame, error) {
	if plan == nil {
		return nil, fmt.Errorf("plan is nil")
	}
	plan.Reset()
	frames := make([]Frame, 0, len(plan.cues))
	sequence := 1
	for {
		cue, ok := plan.Next()
		if !ok {
			break
		}
		frame, err := MakeFrame(sequence, cue, plan.Remaining() == 0)
		if err != nil {
			return nil, err
		}
		frames = append(frames, frame)
		sequence++
	}
	return frames, nil
}
