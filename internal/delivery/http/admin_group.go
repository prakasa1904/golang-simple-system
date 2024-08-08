package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/devetek/go-core/render"
	"github.com/devetek/golang-webapp-boilerplate/internal/helper"
	"github.com/devetek/golang-webapp-boilerplate/internal/services/group"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AdminGroupController struct {
	log       *logrus.Logger
	view      *render.Engine
	myUsecase *group.UseCase
	layout    string
}

func NewAdminGroupController(
	db *gorm.DB,
	log *logrus.Logger,
	view *render.Engine,
	validate *validator.Validate,
	layout string,
) *AdminGroupController {
	// init repositories
	myRepository := group.NewRepository(log)

	// init usecases
	myUsecase := group.NewUseCase(db, log, validate, myRepository)

	return &AdminGroupController{
		log:       log,
		myUsecase: myUsecase,
		view:      view,
		layout:    layout,
	}
}

func (c *AdminGroupController) Home(w http.ResponseWriter, r *http.Request) {
	filter := helper.ConvertQueryToFilter(r.URL, group.AllowedFilterQuery)
	limit := helper.ConvertQueryToLimit(r.URL)
	order := helper.ConvertQueryToOrder(r.URL)

	searchQuery := r.URL.Query().Get("name")

	groups, err := c.myUsecase.Find(r.Context(), filter, limit, order)
	if err != nil {
		c.log.Warnf("Find groups error : %+v", err)
	}

	c.view.Set("search", searchQuery)
	c.view.Set("groups", groups)

	// render page with template html (ejs)
	err = c.view.HTML(w).RenderWithLayout("views/pages/admin/group/admin-group.html", c.layout)
	if err != nil {
		c.log.Warnf("Render error : %+v", err)
	}
}

func (c *AdminGroupController) ComponentList(w http.ResponseWriter, r *http.Request) {
	filter := helper.ConvertQueryToFilter(r.URL, group.AllowedFilterQuery)
	limit := helper.ConvertQueryToLimit(r.URL)
	order := helper.ConvertQueryToOrder(r.URL)

	groups, err := c.myUsecase.Find(r.Context(), filter, limit, order)
	if err != nil {
		c.log.Warnf("Find groups error : %+v", err)
	}

	// require to validate because groups is just pointer
	if groups != nil {
		c.view.Set("groups", groups)
	}

	// render page with template html (ejs)
	err = c.view.HTML(w).RenderClean("views/pages/admin/group/group-item.html")
	if err != nil {
		c.log.Warnf("Render error : %+v", err)
	}
}

func (c *AdminGroupController) ComponentForm(w http.ResponseWriter, r *http.Request) {
	// get action
	action := chi.URLParam(r, "action")
	id := chi.URLParam(r, "id")

	// set all view data, if not set data will use cache and causing invalid data
	c.view.Set("formTitle", strings.ToTitle(action))
	c.view.Set("action", action)
	c.view.Set("group", nil)

	// edit form
	if id != "" {
		group, err := c.myUsecase.GetById(r.Context(), id)
		if err != nil {
			c.log.Warnf("Find groups error : %+v", err)
		}

		// require to validate because groups is just pointer
		if group != nil {
			c.view.Set("group", group)
		}
	}

	// render page with template html (ejs)
	err := c.view.HTML(w).RenderClean("views/pages/admin/group/form-group-content.html")
	if err != nil {
		c.log.Warnf("Render error : %+v", err)
	}
}

func (c *AdminGroupController) MutationCreate(w http.ResponseWriter, r *http.Request) {
	var frontendResp = group.ResponseMutation{
		Status: "Sukses",
	}
	var payload = new(group.RequestPayload)
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		c.log.Warnf("Json decoder error : %+v", err)
	}

	log.Println("payloadpayloadpayload")
	log.Println(payload)
	log.Println("payloadpayloadpayload")

	newGroup, err := c.myUsecase.Create(r.Context(), payload)
	if err != nil {
		frontendResp.Status = "Gagal"
		frontendResp.Message = fmt.Sprintf("Gagal membuat group %s, %v", payload.Name, err)

		c.log.Warnf("Create group error : %+v", err)
	}

	if newGroup != nil {
		// success message
		frontendResp.Message = "Berhasil membuat group " + newGroup.Name
	}

	c.view.Set("toasterTitle", frontendResp.Status)
	c.view.Set("toasterDescription", frontendResp.Message)

	// render page with template html (ejs)
	err = c.view.HTML(w).RenderClean("views/components/toaster/toaster.html")
	if err != nil {
		c.log.Warnf("Render error : %+v", err)
	}
}

func (c *AdminGroupController) MutationUpdate(w http.ResponseWriter, r *http.Request) {
	var frontendResp = group.ResponseMutation{
		Status: "Sukses",
	}
	var payload = new(group.RequestPayload)
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		c.log.Warnf("Json decoder error : %+v", err)
	}

	newGroup, err := c.myUsecase.Update(r.Context(), payload)
	if err != nil {
		frontendResp.Status = "Gagal"
		frontendResp.Message = fmt.Sprintf("Gagal memperbaharui group %s, %v", payload.Name, err)

		c.log.Warnf("Create group error : %+v", err)
	}

	if newGroup != nil {
		// success message
		frontendResp.Message = "Berhasil memperbaharui group " + newGroup.Name
	}

	c.view.Set("toasterTitle", frontendResp.Status)
	c.view.Set("toasterDescription", frontendResp.Message)

	// render page with template html (ejs)
	err = c.view.HTML(w).RenderClean("views/components/toaster/toaster.html")
	if err != nil {
		c.log.Warnf("Render error : %+v", err)
	}
}
