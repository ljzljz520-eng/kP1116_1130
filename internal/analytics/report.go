package analytics

import (
	"fmt"
	"sort"
	"strings"

	"showroom/internal/model"
)

type Report struct {
	SessionID        string
	TotalGestures    int
	AcceptedGestures int
	RejectedGestures int
	DominantGesture  string
	ModesVisited     []model.SceneMode
	Completed        bool
}

func Build(sessionID string, events []model.GestureEvent, states []model.DisplayState) Report {
	report := Report{SessionID: sessionID, TotalGestures: len(events)}
	counts := make(map[string]int)
	modes := make(map[model.SceneMode]bool)
	for _, event := range events {
		if event.Accepted {
			report.AcceptedGestures++
		} else {
			report.RejectedGestures++
		}
		counts[event.Kind]++
	}
	for _, state := range states {
		if state.Mode != "" {
			modes[state.Mode] = true
		}
		if state.Mode == model.ModeClosing {
			report.Completed = true
		}
	}
	for mode := range modes {
		report.ModesVisited = append(report.ModesVisited, mode)
	}
	sort.Slice(report.ModesVisited, func(i, j int) bool { return report.ModesVisited[i] < report.ModesVisited[j] })
	report.DominantGesture = dominant(counts)
	return report
}

func dominant(counts map[string]int) string {
	best := ""
	bestCount := 0
	for kind, count := range counts {
		if count > bestCount || (count == bestCount && kind < best) {
			best = kind
			bestCount = count
		}
	}
	return best
}

func (r Report) CompletionRatio() float64 {
	if r.TotalGestures == 0 {
		return 0
	}
	return float64(r.AcceptedGestures) / float64(r.TotalGestures)
}

func (r Report) Summary() string {
	status := "incomplete"
	if r.Completed {
		status = "completed"
	}
	return fmt.Sprintf("%s:%s:%d/%d", strings.TrimSpace(r.SessionID), status, r.AcceptedGestures, r.TotalGestures)
}
