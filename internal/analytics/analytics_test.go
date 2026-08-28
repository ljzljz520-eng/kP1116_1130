package analytics

import (
	"testing"

	"showroom/internal/model"
)

func TestBuildAndScore(t *testing.T) {
	report := Build("session", []model.GestureEvent{{Kind: "fist", Accepted: true}, {Kind: "wave", Accepted: false}}, []model.DisplayState{{Mode: model.ModeWelcome}, {Mode: model.ModeClosing}})
	if err := ValidateReport(report); err != nil {
		t.Fatal(err)
	}
	if !report.Completed || report.DominantGesture == "" {
		t.Fatalf("unexpected report: %#v", report)
	}
	score := ScoreReport(report)
	if score.Engagement != 20 || score.Reliability != 50 {
		t.Fatalf("unexpected score: %#v", score)
	}
}
