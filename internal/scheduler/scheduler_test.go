package scheduler

import "testing"

func TestDefaultPlanTimeline(t *testing.T) {
	plan := DefaultPlan()
	frames, err := Timeline(plan)
	if err != nil || len(frames) != 4 {
		t.Fatalf("unexpected timeline: %#v %v", frames, err)
	}
	if !frames[len(frames)-1].Final || frames[0].Label != "AMBIENT" {
		t.Fatalf("unexpected frame labels: %#v", frames)
	}
	if err := ValidateFrame(frames[2]); err != nil {
		t.Fatal(err)
	}
}
