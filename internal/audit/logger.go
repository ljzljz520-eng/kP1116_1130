package audit

import (
	"context"
	"fmt"
	"sync"

	"showroom/internal/model"
	"showroom/internal/persistence"
)

type Logger struct {
	store *persistence.Store
	now   func() string
	mu    sync.Mutex
	seq   int
}

func NewLogger(store *persistence.Store, now func() string) *Logger {
	if now == nil {
		now = func() string { return "fixture-now" }
	}
	return &Logger{store: store, now: now}
}

func (l *Logger) Record(ctx context.Context, sessionID, action, detail string) (model.AuditEntry, error) {
	l.mu.Lock()
	l.seq++
	entry := model.AuditEntry{ID: fmt.Sprintf("audit-%03d", l.seq), SessionID: sessionID, Action: action, Detail: detail, CreatedAt: l.now()}
	l.mu.Unlock()
	if entry.SessionID == "" || entry.Action == "" {
		return model.AuditEntry{}, fmt.Errorf("audit session and action are required")
	}
	if err := l.store.SaveAudit(ctx, entry); err != nil {
		return model.AuditEntry{}, err
	}
	return entry, nil
}

func (l *Logger) Sequence() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.seq
}
