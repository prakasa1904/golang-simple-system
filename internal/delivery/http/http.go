package http

import (
	"net/http"

	"github.com/devetek/go-core/render"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prakasa1904/panji-express/internal/delivery/http/middlewares"
)

// register your controller
type RouteConfig struct {
	Router            *chi.Mux
	HomeController    *HomeController
	FindController    *FindController
	AboutController   *AboutController
	ServiceController *ServiceController
	QRController      *QRController

	// admin controller
	AdminDashboardController      *AdminDashboardController
	AdminGroupController          *AdminGroupController
	AdminMemberController         *AdminMemberController
	AdminOrderController          *AdminOrderController
	AdminSettingController        *AdminSettingController
	AdminSettingChannelController *AdminSettingChannelController

	// register API by service
	MemberAPIController *MemberAPIController
}

func (c *RouteConfig) Setup(view *render.Engine) {
	c.SetupMiddleware()
	c.SetupStaticFileServing()
	c.SetupGuestRoute()
	c.SetupComponentRoute()
	c.SetupAPItRoute()
	c.SetupAdminRoute(view)
}

// setup global middleware
func (c *RouteConfig) SetupMiddleware() {
	c.Router.Use(middleware.Logger)
}

// func (c *RouteConfig) SetupAPItRoute() {
// 	// post to mutate data based on repository
// 	c.Router.Route("/api", func(r chi.Router) {
// 		r.Route("/group", func(r chi.Router) {
// 			r.Post("/create", c.GroupAPIController.Create)
// 		})
// 		r.Route("/member", func(r chi.Router) {
// 			r.Post("/add", c.MemberAPIController.Add)
// 			r.Post("/find", c.MemberAPIController.Find)
// 		})
// 		r.Route("/qr", func(r chi.Router) {
// 			r.Post("/", c.QRController.View)
// 		})
// 	})
// }

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

func (c *RouteConfig) SetupAdminRoute(view *render.Engine) {
	// TODO: add admin middleware later

	c.Router.Route("/admin", func(r chi.Router) {
		r.Use(middlewares.AdminSidebar(view))
		r.Route("/", func(r chi.Router) {
			r.Get("/", c.AdminDashboardController.Home)
		})
		r.Route("/group", func(r chi.Router) {
			r.Get("/", c.AdminGroupController.Home)
			// partial component UI
			r.Route("/component", func(r chi.Router) {
				r.Get("/list", c.AdminGroupController.ComponentList)
				r.Get("/form/{action}", c.AdminGroupController.ComponentForm)
				r.Get("/form/{action}/{id}", c.AdminGroupController.ComponentForm)
				r.Get("/delete/{id}", c.AdminGroupController.ComponentDelete)
			})
			// mutation data and return status UI notification depend
			r.Route("/mutation", func(r chi.Router) {
				r.Post("/create", c.AdminGroupController.MutationCreate)
				r.Post("/update", c.AdminGroupController.MutationUpdate)
				r.Delete("/delete/{id}", c.AdminGroupController.MutationDelete)
			})
		})
		r.Route("/member", func(r chi.Router) {
			r.Get("/", c.AdminMemberController.Home)
			// partial component UI
			r.Route("/component", func(r chi.Router) {
				r.Get("/list", c.AdminMemberController.ComponentList)
				r.Get("/form/{action}", c.AdminMemberController.ComponentForm)
				r.Get("/form/{action}/{id}", c.AdminMemberController.ComponentForm)
				r.Get("/delete/{id}", c.AdminMemberController.ComponentDelete)
			})
			// mutation data and return status UI notification depend
			r.Route("/mutation", func(r chi.Router) {
				r.Post("/create", c.AdminMemberController.MutationCreate)
				r.Post("/update", c.AdminMemberController.MutationUpdate)
				r.Delete("/delete/{id}", c.AdminMemberController.MutationDelete)
			})
		})
		r.Route("/order", func(r chi.Router) {
			r.Get("/", c.AdminOrderController.Home)
			// partial component UI
			r.Route("/component", func(r chi.Router) {
				r.Get("/list", c.AdminOrderController.ComponentList)
				r.Get("/form/{action}", c.AdminOrderController.ComponentForm)
				r.Get("/form/{action}/{id}", c.AdminOrderController.ComponentForm)
				r.Get("/delete/{id}", c.AdminOrderController.ComponentDelete)
			})
			// mutation data and return status UI notification depend
			r.Route("/mutation", func(r chi.Router) {
				r.Post("/create", c.AdminOrderController.MutationCreate)
				r.Post("/update", c.AdminOrderController.MutationUpdate)
				r.Delete("/delete/{id}", c.AdminOrderController.MutationDelete)
			})
		})
		// init config for interact with core configuration
		// admin profile
		// whatsApp config
		// etc...
		r.Route("/setting", func(r chi.Router) {
			r.Get("/", c.AdminSettingController.Home)
			r.Route("/channel", func(r chi.Router) {
				r.Get("/", c.AdminSettingChannelController.Home)
				r.Post("/send-message", c.AdminSettingChannelController.SendMessage)
				r.Post("/logout", c.AdminSettingChannelController.Logout)
			})
		})
	})
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

/*
*
API use for frontend
*/
func (c *RouteConfig) SetupAPItRoute() {
	// post to mutate data based on repository
	c.Router.Route("/api", func(r chi.Router) {
		r.Route("/member", func(r chi.Router) {
			r.Post("/add", c.MemberAPIController.Add)
			r.Post("/find", c.MemberAPIController.Find)
		})
		r.Route("/qr", func(r chi.Router) {
			r.Post("/", c.QRController.View)
		})
	})
}
