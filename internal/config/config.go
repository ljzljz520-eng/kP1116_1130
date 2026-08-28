package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Address       string
	DatabasePath  string
	ControlHidden bool
	MaxParticles  int
	GalleryName   string
}

func Default() Config {
	return Config{Address: ":8080", DatabasePath: "showroom.db", MaxParticles: 240, GalleryName: "Lumen Gallery"}
}

func FromEnvironment(getenv func(string) string) Config {
	cfg := Default()
	if value := strings.TrimSpace(getenv("SHOWROOM_ADDRESS")); value != "" {
		cfg.Address = value
	}
	if value := strings.TrimSpace(getenv("SHOWROOM_DB")); value != "" {
		cfg.DatabasePath = value
	}
	if value := strings.TrimSpace(getenv("SHOWROOM_GALLERY")); value != "" {
		cfg.GalleryName = value
	}
	if value := strings.TrimSpace(getenv("SHOWROOM_HIDE_CONTROLS")); value != "" {
		cfg.ControlHidden = value == "1" || strings.EqualFold(value, "true")
	}
	if value := strings.TrimSpace(getenv("SHOWROOM_PARTICLES")); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			cfg.MaxParticles = parsed
		}
	}
	return cfg
}

func Load() Config {
	return FromEnvironment(os.Getenv)
}
