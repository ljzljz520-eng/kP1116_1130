package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"showroom/internal/model"
)

func (s *Store) SavePhrase(ctx context.Context, phrase model.Phrase) error {
	if err := phrase.Validate(); err != nil {
		return err
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO phrases(id,text,mode,priority,enabled,created_at) VALUES(?,?,?,?,?,?)
		ON CONFLICT(id) DO UPDATE SET text=excluded.text, mode=excluded.mode, priority=excluded.priority, enabled=excluded.enabled`,
		phrase.ID, phrase.Text, string(phrase.Mode), phrase.Priority, boolInt(phrase.Enabled), phrase.CreatedAt)
	if err != nil {
		return fmt.Errorf("save phrase: %w", err)
	}
	return nil
}

func (s *Store) ListPhrases(ctx context.Context, mode model.SceneMode) ([]model.Phrase, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id,text,mode,priority,enabled,created_at FROM phrases WHERE mode=? AND enabled=1 ORDER BY priority DESC,id`, mode)
	if err != nil {
		return nil, fmt.Errorf("list phrases: %w", err)
	}
	defer rows.Close()
	var result []model.Phrase
	for rows.Next() {
		var phrase model.Phrase
		var enabled int
		if err := rows.Scan(&phrase.ID, &phrase.Text, &phrase.Mode, &phrase.Priority, &enabled, &phrase.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan phrase: %w", err)
		}
		phrase.Enabled = enabled == 1
		result = append(result, phrase)
	}
	return result, rows.Err()
}

func (s *Store) SaveSession(ctx context.Context, session model.Session) error {
	if err := session.Validate(); err != nil {
		return err
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO sessions(id,visitor_name,mode,started_at,ended_at,active) VALUES(?,?,?,?,?,?)
		ON CONFLICT(id) DO UPDATE SET visitor_name=excluded.visitor_name,mode=excluded.mode,ended_at=excluded.ended_at,active=excluded.active`,
		session.ID, session.VisitorName, string(session.Mode), session.StartedAt, session.EndedAt, boolInt(session.Active))
	if err != nil {
		return fmt.Errorf("save session: %w", err)
	}
	return nil
}

func (s *Store) FindSession(ctx context.Context, id string) (model.Session, error) {
	var session model.Session
	var active int
	err := s.db.QueryRowContext(ctx, `SELECT id,visitor_name,mode,started_at,ended_at,active FROM sessions WHERE id=?`, id).
		Scan(&session.ID, &session.VisitorName, &session.Mode, &session.StartedAt, &session.EndedAt, &active)
	if err != nil {
		return model.Session{}, fmt.Errorf("find session: %w", err)
	}
	session.Active = active == 1
	return session, nil
}

func (s *Store) SaveGesture(ctx context.Context, event model.GestureEvent) error {
	if err := event.Validate(); err != nil {
		return err
	}
	_, err := s.db.ExecContext(ctx, `INSERT OR REPLACE INTO gesture_events(id,session_id,kind,strength,sequence_no,accepted,created_at) VALUES(?,?,?,?,?,?,?)`,
		event.ID, event.SessionID, event.Kind, event.Strength, event.Sequence, boolInt(event.Accepted), event.CreatedAt)
	return err
}

func (s *Store) ListGestures(ctx context.Context, sessionID string) ([]model.GestureEvent, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id,session_id,kind,strength,sequence_no,accepted,created_at FROM gesture_events WHERE session_id=? ORDER BY sequence_no`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []model.GestureEvent
	for rows.Next() {
		var event model.GestureEvent
		var accepted int
		if err := rows.Scan(&event.ID, &event.SessionID, &event.Kind, &event.Strength, &event.Sequence, &accepted, &event.CreatedAt); err != nil {
			return nil, err
		}
		event.Accepted = accepted == 1
		result = append(result, event)
	}
	return result, rows.Err()
}

func (s *Store) SaveDisplayState(ctx context.Context, state model.DisplayState) error {
	if err := state.Validate(); err != nil {
		return err
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO display_states(session_id,mode,phrase_id,phrase_text,particle_form,revision,updated_at) VALUES(?,?,?,?,?,?,?)
		ON CONFLICT(session_id) DO UPDATE SET mode=excluded.mode,phrase_id=excluded.phrase_id,phrase_text=excluded.phrase_text,particle_form=excluded.particle_form,revision=excluded.revision,updated_at=excluded.updated_at`,
		state.SessionID, string(state.Mode), state.PhraseID, state.PhraseText, state.ParticleForm, state.Revision, state.UpdatedAt)
	return err
}

func (s *Store) FindDisplayState(ctx context.Context, sessionID string) (model.DisplayState, error) {
	var state model.DisplayState
	err := s.db.QueryRowContext(ctx, `SELECT session_id,mode,phrase_id,phrase_text,particle_form,revision,updated_at FROM display_states WHERE session_id=?`, sessionID).
		Scan(&state.SessionID, &state.Mode, &state.PhraseID, &state.PhraseText, &state.ParticleForm, &state.Revision, &state.UpdatedAt)
	if err != nil {
		return model.DisplayState{}, fmt.Errorf("find display state: %w", err)
	}
	return state, nil
}

func (s *Store) SaveAudit(ctx context.Context, entry model.AuditEntry) error {
	_, err := s.db.ExecContext(ctx, `INSERT OR REPLACE INTO audit_entries(id,session_id,action,detail,created_at) VALUES(?,?,?,?,?)`, entry.ID, entry.SessionID, entry.Action, entry.Detail, entry.CreatedAt)
	return err
}

func (s *Store) Count(ctx context.Context, table string) (int, error) {
	allowed := map[string]bool{"phrases": true, "sessions": true, "gesture_events": true, "display_states": true, "audit_entries": true}
	if !allowed[table] {
		return 0, fmt.Errorf("unsupported table %q", table)
	}
	var count int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+table).Scan(&count)
	return count, err
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

var _ = sql.ErrNoRows
