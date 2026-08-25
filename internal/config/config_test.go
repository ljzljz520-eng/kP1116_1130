package config

import "testing"

func TestEnvironmentConfig(t *testing.T) {
	values := map[string]string{"SHOWROOM_ADDRESS": "127.0.0.1:9090", "SHOWROOM_DB": "fixture.db", "SHOWROOM_HIDE_CONTROLS": "true", "SHOWROOM_PARTICLES": "80", "SHOWROOM_GALLERY": "Atrium"}
	cfg := FromEnvironment(func(key string) string { return values[key] })
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
	if !cfg.ControlHidden || cfg.MaxParticles != 80 || cfg.Address != "127.0.0.1:9090" {
		t.Fatalf("unexpected config: %#v", cfg)
	}
	if !cfg.IsLoopback() {
		t.Fatal("expected loopback")
	}
}

func TestInvalidConfig(t *testing.T) {
	cfg := Default()
	cfg.Address = "not-an-address"
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected invalid address")
	}
}
