package diagnostics

import (
	"context"
	"fmt"
	"strings"

	"showroom/internal/model"
)

type EventSummary struct {
	SessionID string
	Total     int
	Accepted  int
	Latest    string
}

func SummarizeEvents(sessionID string, events []model.GestureEvent) EventSummary {
	summary := EventSummary{SessionID: sessionID, Total: len(events)}
	for _, event := range events {
		if event.Accepted {
			summary.Accepted++
		}
		if event.Sequence > 0 && (summary.Latest == "" || event.Sequence >= latestSequence(events, summary.Latest)) {
			summary.Latest = event.Kind
		}
	}
	return summary
}

func latestSequence(events []model.GestureEvent, kind string) int {
	sequence := 0
	for _, event := range events {
		if event.Kind == kind && event.Sequence > sequence {
			sequence = event.Sequence
		}
	}
	return sequence
}

func ValidateSummary(summary EventSummary) error {
	if strings.TrimSpace(summary.SessionID) == "" {
		return fmt.Errorf("summary session is required")
	}
	if summary.Accepted < 0 || summary.Accepted > summary.Total {
		return fmt.Errorf("summary accepted count is invalid")
	}
	return nil
}

func CheckContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
