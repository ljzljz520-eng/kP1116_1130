package config

import (
	"fmt"
	"net"
	"path/filepath"
	"strings"
)

func (c Config) Validate() error {
	if strings.TrimSpace(c.Address) == "" {
		return fmt.Errorf("address is required")
	}
	if _, _, err := net.SplitHostPort(c.Address); err != nil {
		return fmt.Errorf("address must contain host and port: %w", err)
	}
	if filepath.Clean(c.DatabasePath) == "." || strings.TrimSpace(c.DatabasePath) == "" {
		return fmt.Errorf("database path is required")
	}
	if c.MaxParticles < 20 || c.MaxParticles > 5000 {
		return fmt.Errorf("max particles must be between 20 and 5000")
	}
	if len(strings.TrimSpace(c.GalleryName)) < 2 {
		return fmt.Errorf("gallery name is too short")
	}
	return nil
}

func (c Config) IsLoopback() bool {
	host, _, err := net.SplitHostPort(c.Address)
	if err != nil {
		return false
	}
	return host == "localhost" || host == "127.0.0.1" || host == "::1" || host == ""
}
