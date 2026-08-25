package welcome

import (
	"context"
	"path/filepath"
	"testing"

	"showroom/internal/model"
	"showroom/internal/persistence"
)

func newWelcomeTestService(t *testing.T) (*Service, *persistence.Store) {
	t.Helper()
	store, err := persistence.Open(filepath.Join(t.TempDir(), "welcome.db"))
	if err != nil {
		t.Fatal(err)
	}
	service := NewService(store, NewPhraseBook(), func() string { return "test-time" })
	if err := service.SeedDefaults(context.Background()); err != nil {
		store.Close()
		t.Fatal(err)
	}
	return service, store
}

func TestWelcomeUsesLatestPhrase(t *testing.T) {
	service, store := newWelcomeTestService(t)
	defer store.Close()
	ctx := context.Background()
	if err := service.ConfirmPhrase(ctx, model.Phrase{ID: "first", Text: "First Confirmed", Mode: model.ModeWelcome, Priority: 50, Enabled: true, CreatedAt: "a"}); err != nil {
		t.Fatal(err)
	}
	if err := service.ConfirmPhrase(ctx, model.Phrase{ID: "second", Text: "Latest Confirmed", Mode: model.ModeWelcome, Priority: 51, Enabled: true, CreatedAt: "b"}); err != nil {
		t.Fatal(err)
	}
	ending, err := service.EndDisplay(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if ending.Text != "Latest Confirmed" {
		t.Fatalf("ending phrase = %q, want latest confirmation", ending.Text)
	}
}

func TestPhraseNormalization(t *testing.T) {
	if got := NormalizePhrase("  Gallery   opens  "); got != "Gallery opens" {
		t.Fatalf("normalized phrase = %q", got)
	}
}

func TestSelectorChoosesEnabledPhrase(t *testing.T) {
	service, store := newWelcomeTestService(t)
	defer store.Close()
	selector := NewSelector(service)
	phrase, err := selector.Choose(context.Background(), model.ModeWelcome, "welcome-default")
	if err != nil || phrase.Text != "Welcome to the Gallery" {
		t.Fatalf("unexpected phrase: %#v %v", phrase, err)
	}
}
