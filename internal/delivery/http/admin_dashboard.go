package http

import (
	"net/http"

	"github.com/devetek/go-core/render"
	"github.com/go-playground/validator/v10"
	"github.com/prakasa1904/panji-express/internal/helper"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AdminDashboardController struct {
	log *logrus.Logger
	// myUsecase *member.UseCase
	view   *render.Engine
	layout string
}

func NewAdminDashboardController(
	db *gorm.DB,
	log *logrus.Logger,
	view *render.Engine,
	validate *validator.Validate,
	layout string,
) *AdminDashboardController {
	// init module repositories
	// myRepository := member.NewRepository(log)

	// init module usecase
	// myUsecase := member.NewUseCase(db, log, validate, myRepository)

	return &AdminDashboardController{
		log: log,
		// myUsecase: myUsecase,
		view:   view,
		layout: layout,
	}
}

func (c *AdminDashboardController) setHeaderMeta() {
	c.view.Set("title", "Dashboard - Admin")
	c.view.Set("description", "Dashboard administrator")
}

func (c *AdminDashboardController) Home(w http.ResponseWriter, r *http.Request) {
	c.setHeaderMeta()

	// render page with template html (ejs)
	if helper.IsHTMXRequest(r.Header) {
		err := c.view.HTML(w).RenderClean("views/pages/admin/dashboard/admin-dashboard.html")
		if err != nil {
			c.log.Warnf("Render error : %+v", err)
		}
	} else {
		err := c.view.HTML(w).RenderWithLayout("views/pages/admin/dashboard/admin-dashboard.html", c.layout)
		if err != nil {
			c.log.Warnf("Render error : %+v", err)
		}
	}
}
