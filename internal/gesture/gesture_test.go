package gesture

import "testing"

func TestParseAndClassify(t *testing.T) {
	signal, err := Parse("fist:88", 4)
	if err != nil {
		t.Fatal(err)
	}
	decision := Classify(signal)
	if !decision.Accepted || decision.Action != ActionHeart {
		t.Fatalf("unexpected decision: %#v", decision)
	}
	if _, err := Parse("fist:bad", 4); err == nil {
		t.Fatal("expected invalid strength")
	}
}

func TestWeakAndUnsupportedGestures(t *testing.T) {
	weak, err := Parse("wave:10", 1)
	if err != nil {
		t.Fatal(err)
	}
	if decision := Classify(weak); decision.Accepted {
		t.Fatal("weak signal should be ignored")
	}
	if decision := Classify(Signal{Kind: "unknown", Strength: 90, Frame: 2}); decision.Action != ActionIgnore {
		t.Fatal("unknown signal should be ignored")
	}
}
