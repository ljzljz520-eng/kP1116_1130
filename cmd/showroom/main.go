package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"showroom/internal/audit"
	"showroom/internal/catalog"
	"showroom/internal/config"
	"showroom/internal/display"
	"showroom/internal/httpapi"
	"showroom/internal/model"
	"showroom/internal/particles"
	"showroom/internal/persistence"
	"showroom/internal/welcome"
	"showroom/internal/workflow"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		return err
	}
	store, err := persistence.Open(cfg.DatabasePath)
	if err != nil {
		return err
	}
	defer store.Close()
	now := func() string { return time.Now().UTC().Format(time.RFC3339) }
	book := welcome.NewPhraseBook()
	fixtureCatalog := catalog.Fixtures()
	for _, mode := range []model.SceneMode{model.ModeWelcome, model.ModeTour, model.ModeClosing, model.ModeIdle} {
		for _, phrase := range fixtureCatalog.Search("", mode) {
			book.Add(phrase)
		}
	}
	welcomeService := welcome.NewService(store, book, now)
	if err := welcomeService.SeedDefaults(context.Background()); err != nil {
		return err
	}
	emitter := particles.NewEmitter(cfg.MaxParticles)
	controller := display.NewControllerWithPanel(store, emitter, now, cfg.ControlHidden)
	logger := audit.NewLogger(store, now)
	flow := workflow.New(store, welcomeService, controller, logger)
	handler := httpapi.NewHandler(flow, cfg.GalleryName)
	server := &http.Server{Addr: cfg.Address, Handler: httpapi.Routes(handler), ReadHeaderTimeout: 3 * time.Second}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stop
		shutdownContext, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownContext)
	}()
	err = server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
