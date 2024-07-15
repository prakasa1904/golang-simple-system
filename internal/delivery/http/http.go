package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// register your controller
type RouteConfig struct {
	Router            *chi.Mux
	HomeController    *HomeController
	FindController    *FindController
	AboutController   *AboutController
	ServiceController *ServiceController

	// register API by service
	MemberAPIController *MemberAPIController
}

func (c *RouteConfig) Setup() {
	c.SetupMiddleware()
	c.SetupStaticFileServing()
	c.SetupGuestRoute()
	c.SetupComponentRoute()
	c.SetupAPItRoute()
}

// setup global middleware
func (c *RouteConfig) SetupMiddleware() {
	c.Router.Use(middleware.Logger)
}

func (c *RouteConfig) SetupStaticFileServing() {
	var fs = http.FileServer(http.Dir("public"))

	c.Router.Handle("/static/*", http.StripPrefix("/static/", fs))
}

func (c *RouteConfig) SetupGuestRoute() {
	c.Router.Get("/", c.HomeController.Home)
	c.Router.Get("/find", c.FindController.Home)
	c.Router.Get("/service", c.ServiceController.Home)
	c.Router.Get("/about", c.AboutController.Home)
}

func (c *RouteConfig) SetupComponentRoute() {
	// server side UI component
	c.Router.Route("/component", func(r chi.Router) {
		r.Route("/home", func(r chi.Router) {
			r.Get("/", c.HomeController.Component)
		})
		r.Route("/find", func(r chi.Router) {
			r.Get("/", c.FindController.Component)
		})
		r.Route("/service", func(r chi.Router) {
			r.Get("/", c.ServiceController.Component)
		})
		r.Route("/about", func(r chi.Router) {
			r.Get("/", c.AboutController.Component)
		})
	})
}

func (c *RouteConfig) SetupAPItRoute() {
	// post to mutate data based on repository
	c.Router.Route("/api", func(r chi.Router) {
		r.Route("/member", func(r chi.Router) {
			r.Post("/add", c.MemberAPIController.Add)
			r.Post("/find", c.MemberAPIController.Find)
		})
	})
}
