package config

import (
	"github.com/devetek/go-core/render"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/prakasa1904/panji-express/internal/delivery/http"
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
	findController := http.NewFindController(config.DB, config.Log, config.View, config.Validate)
	aboutController := http.NewAboutController(config.DB, config.Log, config.View, config.Validate)
	serviceController := http.NewServiceController(config.DB, config.Log, config.View, config.Validate)
	whatsappController := http.NewWhatsappController(config.Config, config.Log, config.View)
	memberAPIController := http.NewMemberAPIController(config.DB, config.Log, config.View, config.Validate)
	qrController := http.NewQRController(config.Config, config.Log)

	// administrator controller
	adminDashboardController := http.NewAdminDashboardController(
		config.DB,
		config.Log,
		config.View,
		config.Validate,
		config.Config.GetString("view.administrator"),
	)
	adminGroupController := http.NewAdminGroupController(
		config.DB,
		config.Log,
		config.View,
		config.Validate,
		config.Config.GetString("view.administrator"),
	)
	adminMemberController := http.NewAdminMemberController(
		config.DB,
		config.Log,
		config.View,
		config.Validate,
		config.Config.GetString("view.administrator"),
	)

	route := &http.RouteConfig{
		Router:                   config.Router,
		FindController:           findController,
		AboutController:          aboutController,
		HomeController:           homeController,
		ServiceController:        serviceController,
		WhatsappController:       whatsappController,
		MemberAPIController:      memberAPIController,
		QRController:             qrController,
		AdminDashboardController: adminDashboardController,
		AdminGroupController:     adminGroupController,
		AdminMemberController:    adminMemberController,
	}

	// init registered router
	// add additional config here if you need to use it on the middleware
	route.Setup(config.View)
}
