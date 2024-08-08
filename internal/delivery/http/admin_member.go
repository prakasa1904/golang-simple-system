package http

/**
Reference: https://github.com/themesberg/flowbite-admin-dashboard
*/

import (
	"net/http"

	"github.com/devetek/go-core/render"
	"github.com/devetek/golang-webapp-boilerplate/internal/helper"
	"github.com/devetek/golang-webapp-boilerplate/internal/services/member"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AdminMemberController struct {
	log       *logrus.Logger
	view      *render.Engine
	myUsecase *member.UseCase
	layout    string
}

func NewAdminMemberController(
	db *gorm.DB,
	log *logrus.Logger,
	view *render.Engine,
	validate *validator.Validate,
	layout string,
) *AdminMemberController {
	// init repositories
	myRepository := member.NewRepository(log)

	// init usecases
	myUsecase := member.NewUseCase(db, log, validate, myRepository)

	return &AdminMemberController{
		log:       log,
		myUsecase: myUsecase,
		view:      view,
		layout:    layout,
	}
}

func (c *AdminMemberController) Home(w http.ResponseWriter, r *http.Request) {

	filter := helper.ConvertQueryToFilter(r.URL, member.AllowedFilterQuery)
	limit := helper.ConvertQueryToLimit(r.URL)
	order := helper.ConvertQueryToOrder(r.URL)

	members, err := c.myUsecase.Find(r.Context(), filter, limit, order)
	if err != nil {
		c.log.Warnf("Find members error : %+v", err)
	}

	c.view.Set("members", members)

	// render page with template html (ejs)
	err = c.view.HTML(w).RenderWithLayout("views/pages/admin/member/admin-member.html", c.layout)
	if err != nil {
		c.log.Warnf("Render error : %+v", err)
	}

}
