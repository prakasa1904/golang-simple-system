package http

import (
	"net/http"

	"github.com/devetek/go-core/render"
	"github.com/devetek/golang-webapp-boilerplate/internal/services/member"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AboutController struct {
	log       *logrus.Logger
	myUsecase *member.UseCase
	view      *render.Engine
}

func NewAboutController(
	db *gorm.DB,
	log *logrus.Logger,
	view *render.Engine,
	validate *validator.Validate,
) *AboutController {
	// init module repositories
	myRepository := member.NewRepository(log)

	// init module usecase
	myUsecase := member.NewUseCase(db, log, validate, myRepository)

	return &AboutController{
		log:       log,
		myUsecase: myUsecase,
		view:      view,
	}
}

func (c *AboutController) setHeaderMeta() {
	c.view.Set("title", "About")
	c.view.Set("description", "About Me")
}

func (c *AboutController) Home(w http.ResponseWriter, r *http.Request) {
	c.setHeaderMeta()

	// render page with template html (ejs)
	err := c.view.HTML(w).Render("views/pages/about/index.html")
	if err != nil {
		c.log.Warnf("Render error : %+v", err)
	}

}

func (c *AboutController) Component(w http.ResponseWriter, r *http.Request) {
	c.setHeaderMeta()

	err := c.view.HTML(w).RenderClean("views/pages/about/component.html")
	if err != nil {
		c.log.Warnf("RenderClean error : %+v", err)
	}
}
