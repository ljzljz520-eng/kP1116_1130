package model

import "testing"

func TestPhraseAndSessionValidation(t *testing.T) {
	phrase := Phrase{ID: "p1", Text: "Hello", Mode: ModeWelcome, Priority: 1}
	if err := phrase.Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (Phrase{Text: "missing"}).Validate(); err == nil {
		t.Fatal("expected missing phrase id")
	}
	session := Session{ID: "s1", Mode: ModeIdle, Active: true}
	if err := session.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestTimelineLatestAndAccepted(t *testing.T) {
	timeline := NewTimeline(
		GestureEvent{Sequence: 2, Kind: "wave", Accepted: true},
		GestureEvent{Sequence: 1, Kind: "fist", Accepted: false},
	)
	latest, ok := timeline.Latest()
	if !ok || latest.Kind != "wave" || timeline.AcceptedCount() != 1 {
		t.Fatalf("unexpected timeline result: %#v %v", latest, ok)
	}
	if !timeline.HasKind("wave") {
		t.Fatal("expected wave")
	}
}
