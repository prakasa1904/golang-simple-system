package main

import (
	"net/http"

	"github.com/devetek/go-core/render"
	"github.com/devetek/golang-webapp-boilerplate/internal"
	"github.com/devetek/golang-webapp-boilerplate/internal/config"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.NewConfig()
	log := config.NewLogger(cfg)
	validate := config.NewValidator()
	db := config.New(config.DatabaseOption{
		Driver:   cfg.GetString("database.driver"),
		DBName:   cfg.GetString("database.name"),
		Username: cfg.GetString("database.username"),
		Password: cfg.GetString("database.password"),
	})
	router := chi.NewRouter()
	view := render.NewEngine(
		internal.Template,
		render.WithValues(
			map[string]any{
				"keywords": "golang, HTMX, Web Platform, SPA, Tailwind",
			}),
		render.WithDefaultLayout(cfg.GetString("view.default")),
	)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		Router:   router,
		View:     view,
		Log:      log,
		Config:   cfg,
		Validate: validate,
	})

	// run app and make fatal omn failed
	log.Fatal(http.ListenAndServe(":"+cfg.GetString("application.port"), router))
}
