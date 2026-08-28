package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"showroom/internal/audit"
	"showroom/internal/display"
	"showroom/internal/particles"
	"showroom/internal/persistence"
	"showroom/internal/welcome"
	"showroom/internal/workflow"
)

func TestRoutesHealthAndGesture(t *testing.T) {
	store, err := persistence.Open(filepath.Join(t.TempDir(), "http.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	now := func() string { return "http-time" }
	service := welcome.NewService(store, welcome.NewPhraseBook(), now)
	if err := service.SeedDefaults(context.Background()); err != nil {
		t.Fatal(err)
	}
	flow := workflow.New(store, service, display.NewController(store, particles.NewEmitter(5), now), audit.NewLogger(store, now))
	handler := NewHandler(flow, "Gallery")
	routes := Routes(handler)
	health := httptest.NewRecorder()
	routes.ServeHTTP(health, httptest.NewRequest(http.MethodGet, "/health", nil))
	if health.Code != http.StatusOK || !strings.Contains(health.Body.String(), "ready") {
		t.Fatalf("unexpected health response: %d %s", health.Code, health.Body.String())
	}
	body := strings.NewReader(`{"session_id":"s1","signal":"fist:88","frame":1}`)
	gestureResponse := httptest.NewRecorder()
	routes.ServeHTTP(gestureResponse, httptest.NewRequest(http.MethodPost, "/gesture", body))
	if gestureResponse.Code != http.StatusOK || !strings.Contains(gestureResponse.Body.String(), "heart") {
		t.Fatalf("unexpected gesture response: %d %s", gestureResponse.Code, gestureResponse.Body.String())
	}
	panelResponse := httptest.NewRecorder()
	routes.ServeHTTP(panelResponse, httptest.NewRequest(http.MethodPost, "/panel", strings.NewReader(`{"action":"hide"}`)))
	if panelResponse.Code != http.StatusOK || !strings.Contains(panelResponse.Body.String(), `"visible":false`) {
		t.Fatalf("unexpected panel response: %d %s", panelResponse.Code, panelResponse.Body.String())
	}
}
