package rehearsal

import (
	"fmt"
	"sort"

	"showroom/internal/gesture"
	"showroom/internal/model"
)

type Step struct {
	Frame   int
	Signal  gesture.Signal
	Action  gesture.Action
	Success bool
	Detail  string
}

type Script struct {
	Name  string
	Steps []Step
}

type Result struct {
	ScriptName string
	Steps      []Step
	Passed     bool
	FinalMode  model.SceneMode
	Forms      []string
}

func DefaultScript() Script {
	return Script{Name: "welcome-handshake", Steps: []Step{
		{Frame: 1, Signal: gesture.Signal{Kind: "open_palm", Strength: 80, Frame: 1}},
		{Frame: 2, Signal: gesture.Signal{Kind: "point", Strength: 90, Frame: 2}},
		{Frame: 3, Signal: gesture.Signal{Kind: "fist", Strength: 95, Frame: 3}},
		{Frame: 4, Signal: gesture.Signal{Kind: "wave", Strength: 90, Frame: 4}},
	}}
}

func ValidateScript(script Script) error {
	if script.Name == "" {
		return fmt.Errorf("script name is required")
	}
	if len(script.Steps) < 4 {
		return fmt.Errorf("script requires four steps")
	}
	last := 0
	for index, step := range script.Steps {
		if step.Frame <= last || step.Signal.Frame != step.Frame {
			return fmt.Errorf("step %d has invalid frame", index)
		}
		if !gesture.Supported(step.Signal.Kind) {
			return fmt.Errorf("step %d uses unsupported gesture", index)
		}
		last = step.Frame
	}
	return nil
}

func Run(script Script) (Result, error) {
	if err := ValidateScript(script); err != nil {
		return Result{}, err
	}
	result := Result{ScriptName: script.Name, Passed: true, FinalMode: model.ModeIdle}
	for index, step := range script.Steps {
		decision := gesture.Classify(step.Signal)
		step.Action = decision.Action
		step.Success = decision.Accepted
		step.Detail = decision.Reason
		if !step.Success {
			result.Passed = false
		}
		switch decision.Action {
		case gesture.ActionWelcome:
			result.FinalMode = model.ModeWelcome
			result.Forms = append(result.Forms, "drift")
		case gesture.ActionTour:
			result.FinalMode = model.ModeTour
			result.Forms = append(result.Forms, "drift")
		case gesture.ActionHeart:
			result.Forms = append(result.Forms, "heart")
		case gesture.ActionRestore:
			result.Forms = append(result.Forms, "drift")
		default:
			result.Forms = append(result.Forms, "unknown")
		}
		result.Steps = append(result.Steps, step)
		if index == len(script.Steps)-1 && decision.Action == gesture.ActionRestore {
			result.FinalMode = model.ModeIdle
		}
	}
	return result, nil
}

func (r Result) SuccessfulSteps() int {
	count := 0
	for _, step := range r.Steps {
		if step.Success {
			count++
		}
	}
	return count
}

func (r Result) OrderedFrames() []int {
	frames := make([]int, 0, len(r.Steps))
	for _, step := range r.Steps {
		frames = append(frames, step.Frame)
	}
	sort.Ints(frames)
	return frames
}

func (r Result) Summary() string {
	return fmt.Sprintf("%s:%d/%d:%s", r.ScriptName, r.SuccessfulSteps(), len(r.Steps), r.FinalMode)
}
