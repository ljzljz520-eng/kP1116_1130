package display

import "testing"

func TestPanelStateTransitions(t *testing.T) {
	panel := NewPanelState(true)
	if panel.IsVisible() || panel.Revision != 1 {
		t.Fatalf("unexpected initial panel: %#v", panel)
	}
	if err := panel.Apply(PanelShow); err != nil {
		t.Fatal(err)
	}
	if !panel.IsVisible() {
		t.Fatal("panel should be visible")
	}
	if err := panel.Apply(PanelFocusHeart); err != nil {
		t.Fatal(err)
	}
	if panel.FocusedAction() != "heart" || !panel.Allows("toggle") {
		t.Fatalf("unexpected panel focus: %#v", panel)
	}
	if err := panel.Apply("not-an-action"); err == nil {
		t.Fatal("expected unsupported panel action")
	}
}

func TestPanelSnapshotRoundTrip(t *testing.T) {
	panel := NewPanelState(false)
	if err := panel.SetFocus("tour"); err != nil {
		t.Fatal(err)
	}
	snapshot := panel.Snapshot()
	data, err := snapshot.Encode()
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodePanelSnapshot(data)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.Summary() != snapshot.Summary() || len(decoded.ActionLabels) != len(PanelActions()) {
		t.Fatalf("snapshot changed during round trip: %#v %#v", snapshot, decoded)
	}
}

func TestPanelConfiguration(t *testing.T) {
	controller := NewController(nil, nil, func() string { return "fixture" })
	if controller.Panel().Visible != true {
		t.Fatal("default controller panel should be visible")
	}
	configured := controller.ConfigurePanel(true)
	if configured.Visible || configured.Mode != PanelHidden {
		t.Fatalf("unexpected hidden panel: %#v", configured)
	}
	updated, err := controller.ApplyPanelAction(PanelToggle)
	if err != nil || !updated.Visible {
		t.Fatalf("unexpected toggled panel: %#v %v", updated, err)
	}
}
