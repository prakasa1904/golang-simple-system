package http

import (
	"net/http"

	"github.com/devetek/go-core/render"
	"github.com/devetek/golang-webapp-boilerplate/internal/services/member"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ServiceController struct {
	log       *logrus.Logger
	myUsecase *member.UseCase
	view      *render.Engine
}

func NewServiceController(
	db *gorm.DB,
	log *logrus.Logger,
	view *render.Engine,
	validate *validator.Validate,
) *ServiceController {
	// init module repositories
	myRepository := member.NewRepository(log)

	// init module usecase
	myUsecase := member.NewUseCase(db, log, validate, myRepository)

	return &ServiceController{
		log:       log,
		myUsecase: myUsecase,
		view:      view,
	}
}

func (c *ServiceController) setHeaderMeta() {
	c.view.Set("title", "Service")
	c.view.Set("description", "My Service")
}

func (c *ServiceController) Home(w http.ResponseWriter, r *http.Request) {
	c.setHeaderMeta()

	// render page with template html (ejs)
	err := c.view.HTML(w).Render("views/pages/service/index.html")
	if err != nil {
		c.log.Warnf("Render error : %+v", err)
	}

}

func (c *ServiceController) Component(w http.ResponseWriter, r *http.Request) {
	c.setHeaderMeta()

	err := c.view.HTML(w).RenderClean("views/pages/service/component.html")
	if err != nil {
		c.log.Warnf("RenderClean error : %+v", err)
	}
}
