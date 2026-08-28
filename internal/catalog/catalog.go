package catalog

import (
	"fmt"
	"sort"
	"strings"

	"showroom/internal/model"
)

type Catalog struct {
	phrases map[string]model.Phrase
}

func New() *Catalog {
	return &Catalog{phrases: make(map[string]model.Phrase)}
}

func (c *Catalog) Add(phrase model.Phrase) error {
	if err := phrase.Validate(); err != nil {
		return err
	}
	if _, exists := c.phrases[phrase.ID]; exists {
		return fmt.Errorf("phrase %s already exists", phrase.ID)
	}
	c.phrases[phrase.ID] = phrase
	return nil
}

func (c *Catalog) Replace(phrase model.Phrase) error {
	if err := phrase.Validate(); err != nil {
		return err
	}
	c.phrases[phrase.ID] = phrase
	return nil
}

func (c *Catalog) Get(id string) (model.Phrase, bool) {
	phrase, ok := c.phrases[id]
	return phrase, ok
}

func (c *Catalog) Search(query string, mode model.SceneMode) []model.Phrase {
	query = strings.ToLower(strings.TrimSpace(query))
	items := make([]model.Phrase, 0)
	for _, phrase := range c.phrases {
		if !phrase.Enabled || (mode != "" && phrase.Mode != mode) {
			continue
		}
		if query == "" || strings.Contains(strings.ToLower(phrase.Text), query) || strings.Contains(strings.ToLower(phrase.ID), query) {
			items = append(items, phrase)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Priority == items[j].Priority {
			return items[i].ID < items[j].ID
		}
		return items[i].Priority > items[j].Priority
	})
	return items
}

func (c *Catalog) Modes() []model.SceneMode {
	seen := make(map[model.SceneMode]bool)
	for _, phrase := range c.phrases {
		seen[phrase.Mode] = true
	}
	result := make([]model.SceneMode, 0, len(seen))
	for mode := range seen {
		result = append(result, mode)
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}

func (c *Catalog) EnabledCount() int {
	count := 0
	for _, phrase := range c.phrases {
		if phrase.Enabled {
			count++
		}
	}
	return count
}
