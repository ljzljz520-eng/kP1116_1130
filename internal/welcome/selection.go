package welcome

import (
	"context"
	"fmt"

	"showroom/internal/model"
)

type Selector struct {
	service *Service
}

func NewSelector(service *Service) *Selector {
	return &Selector{service: service}
}

func (s *Selector) Choose(ctx context.Context, mode model.SceneMode, id string) (model.Phrase, error) {
	if s == nil || s.service == nil {
		return model.Phrase{}, fmt.Errorf("welcome selector is unavailable")
	}
	phrase, err := s.service.book.Find(mode, id)
	if err != nil {
		return model.Phrase{}, err
	}
	if err := s.service.ConfirmPhrase(ctx, phrase); err != nil {
		return model.Phrase{}, err
	}
	return phrase, nil
}

func (s *Selector) ChooseText(ctx context.Context, mode model.SceneMode, id, text string) (model.Phrase, error) {
	phrase := model.Phrase{ID: id, Text: text, Mode: mode, Priority: 20, Enabled: true, CreatedAt: s.service.now()}
	if err := s.service.ConfirmPhrase(ctx, phrase); err != nil {
		return model.Phrase{}, err
	}
	return phrase, nil
}
