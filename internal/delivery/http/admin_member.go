package http

/**
Reference: https://github.com/themesberg/flowbite-admin-dashboard
*/

import (
	"net/http"

	"github.com/devetek/go-core/render"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AdminMemberController struct {
	log    *logrus.Logger
	view   *render.Engine
	layout string
}

func NewAdminMemberController(
	db *gorm.DB,
	log *logrus.Logger,
	view *render.Engine,
	validate *validator.Validate,
	layout string,
) *AdminMemberController {
	return &AdminMemberController{
		log:    log,
		view:   view,
		layout: layout,
	}
}

func (c *AdminMemberController) setHeaderMeta() {
	c.view.Set("title", "Member - Admin")
	c.view.Set("description", "Member administrator")
}

func (c *AdminMemberController) Home(w http.ResponseWriter, r *http.Request) {
	c.setHeaderMeta()

	// render page with template html (ejs)
	err := c.view.HTML(w).RenderWithLayout("views/pages/admin/member/admin-member.html", c.layout)
	if err != nil {
		c.log.Warnf("Render error : %+v", err)
	}

}
