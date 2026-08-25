package session

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"showroom/internal/model"
	"showroom/internal/persistence"
)

type Manager struct {
	store  *persistence.Store
	now    func() string
	mu     sync.RWMutex
	active map[string]model.Session
	closed map[string]model.Session
}

func NewManager(store *persistence.Store, now func() string) *Manager {
	if now == nil {
		now = func() string { return "fixture-now" }
	}
	return &Manager{store: store, now: now, active: make(map[string]model.Session), closed: make(map[string]model.Session)}
}

func (m *Manager) Open(ctx context.Context, id, visitor string, mode model.SceneMode) (model.Session, error) {
	id = strings.TrimSpace(id)
	visitor = strings.TrimSpace(visitor)
	if id == "" {
		return model.Session{}, fmt.Errorf("session id is required")
	}
	if visitor == "" {
		visitor = "guest"
	}
	if mode == model.ModeClosing {
		return model.Session{}, fmt.Errorf("cannot open a closing session")
	}
	m.mu.Lock()
	if existing, ok := m.active[id]; ok {
		m.mu.Unlock()
		return existing, nil
	}
	session := model.Session{ID: id, VisitorName: visitor, Mode: mode, StartedAt: m.now(), Active: true}
	m.active[id] = session
	m.mu.Unlock()
	if err := m.store.SaveSession(ctx, session); err != nil {
		m.mu.Lock()
		delete(m.active, id)
		m.mu.Unlock()
		return model.Session{}, err
	}
	return session, nil
}

func (m *Manager) Close(ctx context.Context, id string) (model.Session, error) {
	m.mu.Lock()
	session, ok := m.active[id]
	if !ok {
		m.mu.Unlock()
		return model.Session{}, fmt.Errorf("session %s is not active", id)
	}
	session.Active = false
	session.EndedAt = m.now()
	session.Mode = model.ModeClosing
	delete(m.active, id)
	m.closed[id] = session
	m.mu.Unlock()
	if err := m.store.SaveSession(ctx, session); err != nil {
		return model.Session{}, err
	}
	return session, nil
}

func (m *Manager) Find(ctx context.Context, id string) (model.Session, error) {
	m.mu.RLock()
	if value, ok := m.active[id]; ok {
		m.mu.RUnlock()
		return value, nil
	}
	if value, ok := m.closed[id]; ok {
		m.mu.RUnlock()
		return value, nil
	}
	m.mu.RUnlock()
	return m.store.FindSession(ctx, id)
}

func (m *Manager) ActiveCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.active)
}

func (m *Manager) ClosedCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.closed)
}
