package http

// initial module to interact with whatsApp and Telegram channel

import (
	"net/http"

	"github.com/devetek/go-core/render"
	"github.com/go-playground/validator/v10"
	"github.com/prakasa1904/panji-express/internal/helper"
	"github.com/prakasa1904/panji-express/internal/services/group"
	"github.com/prakasa1904/panji-express/internal/services/member"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AdminSettingController struct {
	log          *logrus.Logger
	view         *render.Engine
	groupUsecase *group.UseCase
	myUsecase    *member.UseCase
	layout       string
}

func NewAdminSettingController(
	db *gorm.DB,
	log *logrus.Logger,
	view *render.Engine,
	validate *validator.Validate,
	layout string,
) *AdminSettingController {
	// init repositories
	groupRepository := group.NewRepository(log)
	myRepository := member.NewRepository(log)

	// init usecases
	groupUsecase := group.NewUseCase(db, log, validate, groupRepository)
	myUsecase := member.NewUseCase(db, log, validate, myRepository)

	return &AdminSettingController{
		log:          log,
		groupUsecase: groupUsecase,
		myUsecase:    myUsecase,
		view:         view,
		layout:       layout,
	}
}

func (c *AdminSettingController) Home(w http.ResponseWriter, r *http.Request) {
	filter := helper.ConvertQueryToFilter(r.URL, member.AllowedFilterQuery)
	limit := helper.ConvertQueryToLimit(r.URL)
	order := helper.ConvertQueryToOrder(r.URL, "")

	searchQuery := r.URL.Query().Get("channel")

	members, err := c.myUsecase.Find(r.Context(), filter, limit, order)
	if err != nil {
		c.log.Warnf("Find settings error : %+v", err)
	}

	c.view.Set("pageTitle", "Setting")
	c.view.Set("search", searchQuery)
	// require to validate because members is just pointer
	if members != nil {
		c.view.Set("members", members)
	}

	// render page with template html (ejs)
	if helper.IsHTMXRequest(r.Header) {
		err = c.view.HTML(w).RenderClean("views/pages/admin/setting/admin-setting.html")
		if err != nil {
			c.log.Warnf("Render error : %+v", err)
		}
	} else {
		err = c.view.HTML(w).RenderWithLayout("views/pages/admin/setting/admin-setting.html", c.layout)
		if err != nil {
			c.log.Warnf("Render error : %+v", err)
		}
	}
}
