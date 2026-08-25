package display

import (
	"context"
	"fmt"

	"showroom/internal/model"
	"showroom/internal/particles"
	"showroom/internal/persistence"
)

type Controller struct {
	store   *persistence.Store
	emitter *particles.Emitter
	state   model.DisplayState
	panel   PanelState
	now     func() string
}

func NewController(store *persistence.Store, emitter *particles.Emitter, now func() string) *Controller {
	return NewControllerWithPanel(store, emitter, now, false)
}

func NewControllerWithPanel(store *persistence.Store, emitter *particles.Emitter, now func() string, controlsHidden bool) *Controller {
	if now == nil {
		now = func() string { return "fixture-now" }
	}
	return &Controller{store: store, emitter: emitter, panel: NewPanelState(controlsHidden), now: now}
}

func (c *Controller) ShowPhrase(ctx context.Context, sessionID string, phrase model.Phrase) (model.DisplayState, error) {
	if phrase.Mode != model.ModeWelcome && phrase.Mode != model.ModeTour && phrase.Mode != model.ModeClosing {
		return model.DisplayState{}, fmt.Errorf("phrase mode cannot be shown")
	}
	c.state.SessionID = sessionID
	c.state.Mode = phrase.Mode
	c.state.PhraseID = phrase.ID
	c.state.PhraseText = phrase.Text
	c.state.ParticleForm = "drift"
	c.state.Revision++
	c.state.UpdatedAt = c.now()
	if err := c.store.SaveDisplayState(ctx, c.state); err != nil {
		return model.DisplayState{}, err
	}
	return c.state, nil
}

func (c *Controller) ApplyGesture(ctx context.Context, sessionID string, action string) (model.DisplayState, error) {
	form := particles.FormForAction(action)
	if action == "heart" {
		c.state.ParticleForm = "heart"
	} else if action == "restore" {
		c.state.ParticleForm = "drift"
	} else {
		c.state.ParticleForm = form
	}
	c.state.SessionID = sessionID
	c.state.Revision++
	c.state.UpdatedAt = c.now()
	if c.state.Mode == "" {
		c.state.Mode = model.ModeIdle
	}
	if err := c.store.SaveDisplayState(ctx, c.state); err != nil {
		return model.DisplayState{}, err
	}
	return c.state, nil
}

func (c *Controller) Tick(ctx context.Context, step int) (model.DisplayState, error) {
	snapshot := c.emitter.Advance(step)
	if snapshot.Form != c.state.ParticleForm && c.state.ParticleForm != "heart" {
		c.state.ParticleForm = snapshot.Form
	}
	c.state.Revision++
	c.state.UpdatedAt = c.now()
	if c.state.SessionID == "" {
		return c.state, nil
	}
	if err := c.store.SaveDisplayState(ctx, c.state); err != nil {
		return model.DisplayState{}, err
	}
	return c.state, nil
}

func (c *Controller) State() model.DisplayState {
	return c.state
}

func (c *Controller) Panel() PanelSnapshot {
	return c.panel.Snapshot()
}

func (c *Controller) ApplyPanelAction(action PanelAction) (PanelSnapshot, error) {
	if err := c.panel.Apply(action); err != nil {
		return PanelSnapshot{}, err
	}
	return c.panel.Snapshot(), nil
}

func (c *Controller) ConfigurePanel(hidden bool) PanelSnapshot {
	c.panel.SetHidden(hidden)
	return c.panel.Snapshot()
}
