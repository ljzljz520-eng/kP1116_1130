package catalog

import "showroom/internal/model"

func Fixtures() *Catalog {
	catalog := New()
	items := []model.Phrase{
		{ID: "welcome-default", Text: "Welcome to the Gallery", Mode: model.ModeWelcome, Priority: 10, Enabled: true, CreatedAt: "fixture-1"},
		{ID: "tour-default", Text: "Guided Tour Starts Here", Mode: model.ModeTour, Priority: 10, Enabled: true, CreatedAt: "fixture-1"},
		{ID: "closing-default", Text: "Thank You for Visiting", Mode: model.ModeClosing, Priority: 10, Enabled: true, CreatedAt: "fixture-1"},
		{ID: "idle-default", Text: "Move your hand to begin", Mode: model.ModeIdle, Priority: 10, Enabled: true, CreatedAt: "fixture-1"},
		{ID: "welcome-family", Text: "Welcome, families and friends", Mode: model.ModeWelcome, Priority: 5, Enabled: true, CreatedAt: "fixture-1"},
		{ID: "tour-lights", Text: "Follow the light trail to the next room", Mode: model.ModeTour, Priority: 5, Enabled: true, CreatedAt: "fixture-1"},
	}
	for _, item := range items {
		_ = catalog.Add(item)
	}
	return catalog
}

func SelectDefault(catalog *Catalog, mode model.SceneMode) (model.Phrase, bool) {
	items := catalog.Search("", mode)
	if len(items) == 0 {
		return model.Phrase{}, false
	}
	return items[0], true
}
