package persistence

import (
	"context"
	"path/filepath"
	"testing"

	"showroom/internal/model"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "showroom.db")
	ctx := context.Background()
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SavePhrase(ctx, model.Phrase{ID: "phrase-1", Text: "Persisted Welcome", Mode: model.ModeWelcome, Priority: 3, Enabled: true, CreatedAt: "t1"}); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveSession(ctx, model.Session{ID: "session-1", VisitorName: "staff", Mode: model.ModeWelcome, StartedAt: "t1", Active: true}); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveGesture(ctx, model.GestureEvent{ID: "gesture-1", SessionID: "session-1", Kind: "fist", Strength: 90, Sequence: 1, Accepted: true, CreatedAt: "t1"}); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveDisplayState(ctx, model.DisplayState{SessionID: "session-1", Mode: model.ModeWelcome, PhraseID: "phrase-1", PhraseText: "Persisted Welcome", ParticleForm: "drift", Revision: 1, UpdatedAt: "t1"}); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveAudit(ctx, model.AuditEntry{ID: "audit-1", SessionID: "session-1", Action: "seed", Detail: "fixture", CreatedAt: "t1"}); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	store, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	session, err := store.FindSession(ctx, "session-1")
	if err != nil || session.VisitorName != "staff" {
		t.Fatalf("session was not restored: %#v %v", session, err)
	}
	state, err := store.FindDisplayState(ctx, "session-1")
	if err != nil || state.PhraseText != "Persisted Welcome" {
		t.Fatalf("state was not restored: %#v %v", state, err)
	}
	for _, table := range []string{"phrases", "sessions", "gesture_events", "display_states", "audit_entries"} {
		count, err := store.Count(ctx, table)
		if err != nil || count != 1 {
			t.Fatalf("unexpected %s count: %d %v", table, count, err)
		}
	}
}

func TestPhraseOrdering(t *testing.T) {
	store, err := Open(filepath.Join(t.TempDir(), "ordering.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	ctx := context.Background()
	for _, phrase := range []model.Phrase{{ID: "low", Text: "Low", Mode: model.ModeWelcome, Priority: 1, Enabled: true, CreatedAt: "t"}, {ID: "high", Text: "High", Mode: model.ModeWelcome, Priority: 5, Enabled: true, CreatedAt: "t"}} {
		if err := store.SavePhrase(ctx, phrase); err != nil {
			t.Fatal(err)
		}
	}
	phrases, err := store.ListPhrases(ctx, model.ModeWelcome)
	if err != nil || len(phrases) != 2 || phrases[0].ID != "high" {
		t.Fatalf("unexpected ordering: %#v %v", phrases, err)
	}
}
