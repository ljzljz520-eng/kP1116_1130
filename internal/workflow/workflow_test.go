package workflow

import (
	"context"
	"path/filepath"
	"testing"

	"showroom/internal/audit"
	"showroom/internal/display"
	"showroom/internal/gesture"
	"showroom/internal/model"
	"showroom/internal/particles"
	"showroom/internal/persistence"
	"showroom/internal/welcome"
)

func newWorkflow(t *testing.T) (*Orchestrator, *persistence.Store) {
	t.Helper()
	store, err := persistence.Open(filepath.Join(t.TempDir(), "workflow.db"))
	if err != nil {
		t.Fatal(err)
	}
	now := func() string { return "workflow-time" }
	service := welcome.NewService(store, welcome.NewPhraseBook(), now)
	if err := service.SeedDefaults(context.Background()); err != nil {
		store.Close()
		t.Fatal(err)
	}
	flow := New(store, service, display.NewController(store, particles.NewEmitter(8), now), audit.NewLogger(store, now))
	return flow, store
}

func TestPrimaryWorkflow(t *testing.T) {
	flow, store := newWorkflow(t)
	defer store.Close()
	ctx := context.Background()
	if err := flow.StartSession(ctx, model.Session{ID: "s1", VisitorName: "staff", Mode: model.ModeWelcome, StartedAt: "start", Active: true}); err != nil {
		t.Fatal(err)
	}
	decision, state, err := flow.HandleSignal(ctx, "s1", gesture.Signal{Kind: "open_palm", Strength: 80, Frame: 1})
	if err != nil || !decision.Accepted || state.PhraseText != "Welcome to the Gallery" {
		t.Fatalf("unexpected welcome workflow: %#v %#v %v", decision, state, err)
	}
}

func TestSecondaryWorkflow(t *testing.T) {
	flow, store := newWorkflow(t)
	defer store.Close()
	ctx := context.Background()
	if _, err := flow.ConfirmCustom(ctx, "s2", model.ModeTour, "tour-custom", "Custom Tour"); err != nil {
		t.Fatal(err)
	}
	_, state, err := flow.HandleSignal(ctx, "s2", gesture.Signal{Kind: "fist", Strength: 95, Frame: 2})
	if err != nil || state.ParticleForm != "heart" {
		t.Fatalf("unexpected heart workflow: %#v %v", state, err)
	}
	_, state, err = flow.HandleSignal(ctx, "s2", gesture.Signal{Kind: "wave", Strength: 95, Frame: 3})
	if err != nil || state.ParticleForm != "drift" {
		t.Fatalf("unexpected restore workflow: %#v %v", state, err)
	}
}

func TestTertiaryWorkflow(t *testing.T) {
	flow, store := newWorkflow(t)
	defer store.Close()
	ctx := context.Background()
	if _, err := flow.ConfirmCustom(ctx, "s3", model.ModeWelcome, "custom-a", "First Phrase"); err != nil {
		t.Fatal(err)
	}
	if _, err := flow.ConfirmCustom(ctx, "s3", model.ModeWelcome, "custom-b", "Final Phrase"); err != nil {
		t.Fatal(err)
	}
	state, err := flow.EndDisplay(ctx, "s3")
	if err != nil {
		t.Fatal(err)
	}
	if state.Mode != model.ModeClosing {
		t.Fatalf("expected closing mode: %#v", state)
	}
}

func TestWorkflowTick(t *testing.T) {
	flow, store := newWorkflow(t)
	defer store.Close()
	if _, err := flow.Tick(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
}
