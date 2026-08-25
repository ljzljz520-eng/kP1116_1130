package display

import (
	"context"
	"path/filepath"
	"testing"

	"showroom/internal/model"
	"showroom/internal/particles"
	"showroom/internal/persistence"
)

func TestControllerPersistsFrame(t *testing.T) {
	store, err := persistence.Open(filepath.Join(t.TempDir(), "display.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	controller := NewController(store, particles.NewEmitter(4), func() string { return "frame-time" })
	state, err := controller.ShowPhrase(context.Background(), "session-1", model.Phrase{ID: "welcome", Text: "Hello", Mode: model.ModeWelcome})
	if err != nil || state.Revision != 1 {
		t.Fatalf("unexpected state: %#v %v", state, err)
	}
	state, err = controller.ApplyGesture(context.Background(), "session-1", "heart")
	if err != nil || state.ParticleForm != "heart" {
		t.Fatalf("unexpected heart state: %#v %v", state, err)
	}
	frame := BuildFrame(state, particles.NewEmitter(4).Current(), "Gallery")
	if len(PlainSummary(frame)) == 0 {
		t.Fatal("expected summary")
	}
}
