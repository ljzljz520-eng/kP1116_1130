package persistence

import (
	"context"
	"fmt"
	"sort"

	"showroom/internal/model"
)

type SessionBundle struct {
	Session  model.Session
	Gestures []model.GestureEvent
	State    model.DisplayState
}

func (s *Store) LoadSessionBundle(ctx context.Context, sessionID string) (SessionBundle, error) {
	session, err := s.FindSession(ctx, sessionID)
	if err != nil {
		return SessionBundle{}, err
	}
	gestures, err := s.ListGestures(ctx, sessionID)
	if err != nil {
		return SessionBundle{}, fmt.Errorf("load gestures: %w", err)
	}
	state, err := s.FindDisplayState(ctx, sessionID)
	if err != nil {
		return SessionBundle{}, fmt.Errorf("load display state: %w", err)
	}
	return SessionBundle{Session: session, Gestures: gestures, State: state}, nil
}

func (s *Store) SaveBundle(ctx context.Context, bundle SessionBundle) error {
	if err := bundle.Session.Validate(); err != nil {
		return err
	}
	if err := bundle.State.Validate(); err != nil {
		return err
	}
	if bundle.State.SessionID != bundle.Session.ID {
		return fmt.Errorf("bundle session ids do not match")
	}
	if err := s.SaveSession(ctx, bundle.Session); err != nil {
		return err
	}
	for _, event := range bundle.Gestures {
		if err := s.SaveGesture(ctx, event); err != nil {
			return err
		}
	}
	return s.SaveDisplayState(ctx, bundle.State)
}

func (s *Store) DisablePhrase(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, `UPDATE phrases SET enabled=0 WHERE id=?`, id)
	if err != nil {
		return err
	}
	if changed, err := result.RowsAffected(); err != nil {
		return err
	} else if changed == 0 {
		return fmt.Errorf("phrase %s does not exist", id)
	}
	return nil
}

func (s *Store) LatestGesture(ctx context.Context, sessionID string) (model.GestureEvent, error) {
	var event model.GestureEvent
	var accepted int
	err := s.db.QueryRowContext(ctx, `SELECT id,session_id,kind,strength,sequence_no,accepted,created_at FROM gesture_events WHERE session_id=? ORDER BY sequence_no DESC LIMIT 1`, sessionID).
		Scan(&event.ID, &event.SessionID, &event.Kind, &event.Strength, &event.Sequence, &accepted, &event.CreatedAt)
	if err != nil {
		return model.GestureEvent{}, err
	}
	event.Accepted = accepted == 1
	return event, nil
}

func SortGestures(events []model.GestureEvent) []model.GestureEvent {
	items := append([]model.GestureEvent(nil), events...)
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Sequence == items[j].Sequence {
			return items[i].ID < items[j].ID
		}
		return items[i].Sequence < items[j].Sequence
	})
	return items
}

func (s *Store) AuditTrail(ctx context.Context, sessionID string) ([]model.AuditEntry, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id,session_id,action,detail,created_at FROM audit_entries WHERE session_id=? ORDER BY created_at,id`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	entries := make([]model.AuditEntry, 0)
	for rows.Next() {
		var entry model.AuditEntry
		if err := rows.Scan(&entry.ID, &entry.SessionID, &entry.Action, &entry.Detail, &entry.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}
