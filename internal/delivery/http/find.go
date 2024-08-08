package http

import (
	"net/http"

	"github.com/devetek/go-core/render"
	"github.com/devetek/golang-webapp-boilerplate/internal/helper"
	"github.com/devetek/golang-webapp-boilerplate/internal/services/member"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type FindController struct {
	log       *logrus.Logger
	myUsecase *member.UseCase
	view      *render.Engine
}

func NewFindController(
	db *gorm.DB,
	log *logrus.Logger,
	view *render.Engine,
	validate *validator.Validate,
) *FindController {
	// init module repositories
	myRepository := member.NewRepository(log)

	// init module usecase
	myUsecase := member.NewUseCase(db, log, validate, myRepository)

	return &FindController{
		log:       log,
		myUsecase: myUsecase,
		view:      view,
	}
}

func (c *FindController) setHeaderMeta() {
	c.view.Set("title", "Panji Express")
	c.view.Set("description", "Jasa pengiriman paket Jakarta - Kalimantan same day")
}

func (c *FindController) Home(w http.ResponseWriter, r *http.Request) {
	c.setHeaderMeta()

	filter := helper.ConvertQueryToFilter(r.URL, member.AllowedFilterQuery)
	limit := helper.ConvertQueryToLimit(r.URL)
	order := helper.ConvertQueryToOrder(r.URL, "")

	users, err := c.myUsecase.Find(r.Context(), filter, limit, order)
	if err != nil {
		c.log.Warnf("Find users error : %+v", err)
	}

	c.view.Set("users", users)

	// render page with template html (ejs)
	err = c.view.HTML(w).Render("views/pages/find/index.html")
	if err != nil {
		c.log.Warnf("Render error : %+v", err)
	}

}

func (c *FindController) Component(w http.ResponseWriter, r *http.Request) {
	c.setHeaderMeta()

	filter := helper.ConvertQueryToFilter(r.URL, member.AllowedFilterQuery)
	limit := helper.ConvertQueryToLimit(r.URL)
	order := helper.ConvertQueryToOrder(r.URL, "")

	users, err := c.myUsecase.Find(r.Context(), filter, limit, order)
	if err != nil {
		c.log.Warnf("Find users error : %+v", err)
	}

	c.view.Set("users", users)

	err = c.view.HTML(w).RenderClean("views/pages/find/component.html")
	if err != nil {
		c.log.Warnf("RenderClean error : %+v", err)
	}
}
