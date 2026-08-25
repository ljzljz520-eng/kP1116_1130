package particles

import "testing"

func TestEmitterDriftsAndFormsHeart(t *testing.T) {
	emitter := NewEmitter(12)
	initial := emitter.Current()
	advanced := emitter.Advance(2)
	if advanced.Frame != initial.Frame+2 || advanced.Count != 12 {
		t.Fatalf("unexpected advance: %#v", advanced)
	}
	heart := emitter.SetForm("heart")
	if heart.Form != "heart" {
		t.Fatal("expected heart form")
	}
	reset := emitter.Reset()
	if reset.Form != "drift" || reset.Frame != 0 {
		t.Fatal("expected drift reset")
	}
}
