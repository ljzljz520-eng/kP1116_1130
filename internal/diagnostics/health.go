package diagnostics

import (
	"context"
	"fmt"
	"sort"

	"showroom/internal/config"
	"showroom/internal/persistence"
)

type HealthReport struct {
	Ready       bool            `json:"ready"`
	Checks      map[string]bool `json:"checks"`
	RowCounts   map[string]int  `json:"row_counts"`
	Address     string          `json:"address"`
	GalleryName string          `json:"gallery_name"`
}

func Inspect(ctx context.Context, store *persistence.Store, cfg config.Config) (HealthReport, error) {
	report := HealthReport{Checks: make(map[string]bool), RowCounts: make(map[string]int), Address: cfg.Address, GalleryName: cfg.GalleryName}
	if err := cfg.Validate(); err != nil {
		report.Checks["config"] = false
		return report, err
	}
	report.Checks["config"] = true
	if err := store.Ping(ctx); err != nil {
		report.Checks["database"] = false
		return report, err
	}
	report.Checks["database"] = true
	for _, table := range []string{"phrases", "sessions", "gesture_events", "display_states", "audit_entries"} {
		count, err := store.Count(ctx, table)
		if err != nil {
			return report, fmt.Errorf("inspect %s: %w", table, err)
		}
		report.RowCounts[table] = count
	}
	report.Ready = true
	return report, nil
}

func (r HealthReport) FailedChecks() []string {
	items := make([]string, 0)
	for name, ok := range r.Checks {
		if !ok {
			items = append(items, name)
		}
	}
	sort.Strings(items)
	return items
}

func (r HealthReport) Summary() string {
	if r.Ready {
		return "showroom ready"
	}
	return fmt.Sprintf("showroom unavailable: %v", r.FailedChecks())
}
