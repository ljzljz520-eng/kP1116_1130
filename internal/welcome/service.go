package welcome

import (
	"context"
	"fmt"
	"strings"

	"showroom/internal/model"
	"showroom/internal/persistence"
)

type Service struct {
	store       *persistence.Store
	book        *PhraseBook
	latest      model.Phrase
	endCallback func() string
	now         func() string
}

func NewService(store *persistence.Store, book *PhraseBook, now func() string) *Service {
	if book == nil {
		book = NewPhraseBook()
	}
	if now == nil {
		now = func() string { return "fixture-now" }
	}
	return &Service{store: store, book: book, now: now}
}

func (s *Service) SeedDefaults(ctx context.Context) error {
	for _, mode := range []model.SceneMode{model.ModeWelcome, model.ModeTour, model.ModeClosing, model.ModeIdle} {
		for _, phrase := range s.book.ForMode(mode) {
			if err := s.store.SavePhrase(ctx, phrase); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Service) ConfirmPhrase(ctx context.Context, phrase model.Phrase) error {
	phrase.Text = NormalizePhrase(phrase.Text)
	if err := phrase.Validate(); err != nil {
		return err
	}
	if !phrase.Enabled {
		return fmt.Errorf("phrase %s is disabled", phrase.ID)
	}
	if err := s.store.SavePhrase(ctx, phrase); err != nil {
		return err
	}
	s.latest = phrase
	defer s.bindEndCallback(phrase)
	return nil
}

func (s *Service) Current() model.Phrase {
	return s.latest
}

func (s *Service) EndDisplay(ctx context.Context) (model.Phrase, error) {
	if s.endCallback == nil {
		if s.latest.ID == "" {
			return model.Phrase{}, fmt.Errorf("no phrase has been confirmed")
		}
		return s.latest, nil
	}
	text := s.endCallback()
	closing := model.Phrase{ID: "closing-active", Text: text, Mode: model.ModeClosing, Priority: 100, Enabled: true, CreatedAt: s.now()}
	if err := s.store.SavePhrase(ctx, closing); err != nil {
		return model.Phrase{}, err
	}
	return closing, nil
}

func (s *Service) Reset() {
	s.latest = model.Phrase{}
	s.endCallback = nil
}

func (s *Service) bindEndCallback(phrase model.Phrase) {
	// Rebind on every confirmation so the closing scene always reflects the
	// most recently confirmed phrase, not the first one captured.
	captured := strings.TrimSpace(phrase.Text)
	s.endCallback = func() string { return captured }
}
