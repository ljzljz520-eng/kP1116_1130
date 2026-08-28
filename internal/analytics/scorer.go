package analytics

import (
	"fmt"
	"math"
	"strings"

	"showroom/internal/model"
)

type Score struct {
	Engagement  int
	Clarity     int
	Reliability int
	Notes       []string
}

func ScoreReport(report Report) Score {
	score := Score{}
	score.Engagement = clampInt(report.AcceptedGestures*20, 0, 100)
	score.Clarity = clampInt(len(report.ModesVisited)*25, 0, 100)
	score.Reliability = clampInt(int(math.Round(report.CompletionRatio()*100)), 0, 100)
	if report.DominantGesture == "fist" {
		score.Notes = append(score.Notes, "heart interaction detected")
	}
	if report.Completed {
		score.Notes = append(score.Notes, "closing handoff completed")
	}
	if report.RejectedGestures > report.AcceptedGestures {
		score.Notes = append(score.Notes, "gesture confidence needs tuning")
	}
	return score
}

func clampInt(value, low, high int) int {
	if value < low {
		return low
	}
	if value > high {
		return high
	}
	return value
}

func ValidateReport(report Report) error {
	if strings.TrimSpace(report.SessionID) == "" {
		return fmt.Errorf("report session is required")
	}
	if report.TotalGestures < 0 || report.AcceptedGestures < 0 || report.RejectedGestures < 0 {
		return fmt.Errorf("report counts cannot be negative")
	}
	if report.AcceptedGestures+report.RejectedGestures != report.TotalGestures {
		return fmt.Errorf("report counts do not balance")
	}
	for _, mode := range report.ModesVisited {
		if mode != model.ModeWelcome && mode != model.ModeTour && mode != model.ModeClosing && mode != model.ModeIdle {
			return fmt.Errorf("report contains unknown mode")
		}
	}
	return nil
}
