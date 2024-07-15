package config

import (
	"github.com/devetek/go-core/render"
	"github.com/devetek/golang-webapp-boilerplate/internal/delivery/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	Config   *viper.Viper
	DB       *gorm.DB
	Router   *chi.Mux
	View     *render.Engine
	Log      *logrus.Logger
	Validate *validator.Validate
}

func Bootstrap(config *BootstrapConfig) {
	homeController := http.NewHomeController(config.DB, config.Log, config.View, config.Validate)
	aboutController := http.NewAboutController(config.DB, config.Log, config.View, config.Validate)
	serviceController := http.NewServiceController(config.DB, config.Log, config.View, config.Validate)
	memberAPIController := http.NewMemberAPIController(config.DB, config.Log, config.View, config.Validate)

	route := &http.RouteConfig{
		Router:              config.Router,
		AboutController:     aboutController,
		HomeController:      homeController,
		ServiceController:   serviceController,
		MemberAPIController: memberAPIController,
	}

	// init registered router
	route.Setup()
}
