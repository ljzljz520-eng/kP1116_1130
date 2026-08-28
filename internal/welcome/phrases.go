package welcome

import (
	"fmt"
	"strings"

	"showroom/internal/model"
)

type PhraseBook struct {
	items map[model.SceneMode][]model.Phrase
}

func NewPhraseBook() *PhraseBook {
	book := &PhraseBook{items: make(map[model.SceneMode][]model.Phrase)}
	book.Add(model.Phrase{ID: "welcome-default", Text: "Welcome to the Gallery", Mode: model.ModeWelcome, Priority: 10, Enabled: true, CreatedAt: "fixture-1"})
	book.Add(model.Phrase{ID: "tour-default", Text: "Guided Tour Starts Here", Mode: model.ModeTour, Priority: 10, Enabled: true, CreatedAt: "fixture-1"})
	book.Add(model.Phrase{ID: "closing-default", Text: "Thank You for Visiting", Mode: model.ModeClosing, Priority: 10, Enabled: true, CreatedAt: "fixture-1"})
	book.Add(model.Phrase{ID: "idle-default", Text: "Move your hand to begin", Mode: model.ModeIdle, Priority: 10, Enabled: true, CreatedAt: "fixture-1"})
	return book
}

func (b *PhraseBook) Add(phrase model.Phrase) {
	if b.items == nil {
		b.items = make(map[model.SceneMode][]model.Phrase)
	}
	b.items[phrase.Mode] = append(b.items[phrase.Mode], phrase)
}

func (b *PhraseBook) ForMode(mode model.SceneMode) []model.Phrase {
	items := append([]model.Phrase(nil), b.items[mode]...)
	return items
}

func (b *PhraseBook) Find(mode model.SceneMode, id string) (model.Phrase, error) {
	for _, phrase := range b.items[mode] {
		if phrase.ID == id && phrase.Enabled {
			return phrase, nil
		}
	}
	return model.Phrase{}, fmt.Errorf("phrase %q is not enabled for %s", id, mode)
}

func NormalizePhrase(text string) string {
	return strings.Join(strings.Fields(text), " ")
}
