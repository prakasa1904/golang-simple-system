package http

/**
Reference: https://github.com/themesberg/flowbite-admin-dashboard
*/

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/devetek/go-core/render"
	"github.com/devetek/golang-webapp-boilerplate/internal/helper"
	"github.com/devetek/golang-webapp-boilerplate/internal/model"
	"github.com/devetek/golang-webapp-boilerplate/internal/services/group"
	"github.com/devetek/golang-webapp-boilerplate/internal/services/member"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AdminMemberController struct {
	log          *logrus.Logger
	view         *render.Engine
	groupUsecase *group.UseCase
	myUsecase    *member.UseCase
	layout       string
}

func NewAdminMemberController(
	db *gorm.DB,
	log *logrus.Logger,
	view *render.Engine,
	validate *validator.Validate,
	layout string,
) *AdminMemberController {
	// init repositories
	groupRepository := group.NewRepository(log)
	myRepository := member.NewRepository(log)

	// init usecases
	groupUsecase := group.NewUseCase(db, log, validate, groupRepository)
	myUsecase := member.NewUseCase(db, log, validate, myRepository)

	return &AdminMemberController{
		log:          log,
		groupUsecase: groupUsecase,
		myUsecase:    myUsecase,
		view:         view,
		layout:       layout,
	}
}

func (c *AdminMemberController) Home(w http.ResponseWriter, r *http.Request) {
	filter := helper.ConvertQueryToFilter(r.URL, member.AllowedFilterQuery)
	limit := helper.ConvertQueryToLimit(r.URL)
	order := helper.ConvertQueryToOrder(r.URL, "`member`.`id`")

	searchQuery := r.URL.Query().Get("fullname")

	members, err := c.myUsecase.Find(r.Context(), filter, limit, order)
	if err != nil {
		c.log.Warnf("Find members error : %+v", err)
	}

	c.view.Set("pageTitle", "Semua Member")
	c.view.Set("search", searchQuery)
	// require to validate because members is just pointer
	if members != nil {
		c.view.Set("members", members)
	}

	// render page with template html (ejs)
	if helper.IsHTMXRequest(r.Header) {
		err = c.view.HTML(w).RenderClean("views/pages/admin/member/admin-member.html")
		if err != nil {
			c.log.Warnf("Render error : %+v", err)
		}
	} else {
		err = c.view.HTML(w).RenderWithLayout("views/pages/admin/member/admin-member.html", c.layout)
		if err != nil {
			c.log.Warnf("Render error : %+v", err)
		}
	}
}

func (c *AdminMemberController) ComponentList(w http.ResponseWriter, r *http.Request) {
	filter := helper.ConvertQueryToFilter(r.URL, member.AllowedFilterQuery)
	limit := helper.ConvertQueryToLimit(r.URL)
	order := helper.ConvertQueryToOrder(r.URL, "`member`.`id`")

	members, err := c.myUsecase.Find(r.Context(), filter, limit, order)
	if err != nil {
		c.log.Warnf("Find members error : %+v", err)
	}

	// require to validate because members is just pointer
	if members != nil {
		c.view.Set("members", members)
	}

	// render page with template html (ejs)
	err = c.view.HTML(w).RenderClean("views/pages/admin/member/member-item.html")
	if err != nil {
		c.log.Warnf("Render error : %+v", err)
	}
}

func (c *AdminMemberController) ComponentForm(w http.ResponseWriter, r *http.Request) {
	// get action
	action := chi.URLParam(r, "action")
	id := chi.URLParam(r, "id")

	// set all view data, if not set data will use cache and causing invalid data
	c.view.Set("formTitle", strings.ToTitle(action))
	c.view.Set("action", action)
	c.view.Set("groups", nil)
	c.view.Set("member", nil)

	filter := helper.ConvertQueryToFilter(r.URL, group.AllowedFilterQuery)
	limit := helper.ConvertQueryToLimit(r.URL)
	order := helper.ConvertQueryToOrder(r.URL, "")

	// get list of group show select field
	groups, err := c.groupUsecase.Find(r.Context(), filter, limit, order)
	if err != nil {
		c.log.Warnf("Find groups error : %+v", err)
	}

	if groups != nil {
		c.view.Set("groups", groups)
	}

	// edit form
	if id != "" {
		member, err := c.myUsecase.GetById(r.Context(), id)
		if err != nil {
			c.log.Warnf("Find members error : %+v", err)
		}

		// require to validate because members is just pointer
		if member != nil {
			c.view.Set("member", member)
		}
	}

	// render page with template html (ejs)
	err = c.view.HTML(w).RenderClean("views/pages/admin/member/form-member-content.html")
	if err != nil {
		c.log.Warnf("Render error : %+v", err)
	}
}

func (c *AdminMemberController) ComponentDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// set all view data, if not set data will use cache and causing invalid data
	c.view.Set("member", nil)

	// edit form
	if id != "" {
		member, err := c.myUsecase.GetById(r.Context(), id)
		if err != nil {
			c.log.Warnf("Find members error : %+v", err)
		}

		// require to validate because members is just pointer
		if member != nil {
			c.view.Set("member", member)
		}
	}

	// render page with template html (ejs)
	err := c.view.HTML(w).RenderClean("views/pages/admin/member/delete-member-content.html")
	if err != nil {
		c.log.Warnf("Render error : %+v", err)
	}
}

func (c *AdminMemberController) MutationCreate(w http.ResponseWriter, r *http.Request) {
	var frontendResp = model.ResponseMutation{
		Status: "Sukses",
	}
	var payload = new(member.RequestPayload)
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		c.log.Warnf("Json decoder error : %+v", err)
	}

	// generate temporary Username
	payload.Username = helper.GenerateUsernameFromEmail(payload.Email)

	newmember, err := c.myUsecase.Create(r.Context(), payload)
	if err != nil {
		frontendResp.Status = "Gagal"
		frontendResp.Message = fmt.Sprintf("Gagal membuat member %s, %v", payload.Fullname, err)

		c.log.Warnf("Create member error : %+v", err)
	}

	if newmember != nil {
		// success message
		frontendResp.Message = "Berhasil membuat member " + newmember.Fullname
	}

	c.view.Set("toasterTitle", frontendResp.Status)
	c.view.Set("toasterDescription", frontendResp.Message)

	// render page with template html (ejs)
	err = c.view.HTML(w).RenderClean("views/components/toaster/toaster.html")
	if err != nil {
		c.log.Warnf("Render error : %+v", err)
	}
}

func (c *AdminMemberController) MutationUpdate(w http.ResponseWriter, r *http.Request) {
	var frontendResp = model.ResponseMutation{
		Status: "Sukses",
	}
	var payload = new(member.RequestPayload)
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		c.log.Warnf("Json decoder error : %+v", err)
	}

	updatedmember, err := c.myUsecase.Update(r.Context(), payload)
	if err != nil {
		frontendResp.Status = "Gagal"
		frontendResp.Message = fmt.Sprintf("Gagal memperbaharui member, %v", err)

		c.log.Warnf("Update member error : %+v", err)
	}

	if updatedmember != nil {
		// success message
		frontendResp.Message = "Berhasil memperbaharui member " + updatedmember.Fullname
	}

	c.view.Set("toasterTitle", frontendResp.Status)
	c.view.Set("toasterDescription", frontendResp.Message)

	// render page with template html (ejs)
	err = c.view.HTML(w).RenderClean("views/components/toaster/toaster.html")
	if err != nil {
		c.log.Warnf("Render error : %+v", err)
	}
}

func (c *AdminMemberController) MutationDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var frontendResp = model.ResponseMutation{
		Status: "Sukses",
	}

	var payload = new(member.DeletePayload)
	if id != "" {
		data, err := c.myUsecase.GetById(r.Context(), id)
		if err != nil {
			c.log.Warnf("Find members error : %+v", err)
		}

		// require to validate because members is just pointer
		if data != nil {
			payload.ID = id
		}
	}

	deletedmember, err := c.myUsecase.Delete(r.Context(), payload)
	if err != nil {
		frontendResp.Status = "Gagal"
		frontendResp.Message = fmt.Sprintf("Gagal menghapus member, %v", err)

		c.log.Warnf("Delete member error : %+v", err)
	}

	if deletedmember != nil {
		// success message
		frontendResp.Message = "Berhasil menghapus member " + deletedmember.Fullname
	}

	c.view.Set("toasterTitle", frontendResp.Status)
	c.view.Set("toasterDescription", frontendResp.Message)

	// render page with template html (ejs)
	err = c.view.HTML(w).RenderClean("views/components/toaster/toaster.html")
	if err != nil {
		c.log.Warnf("Render error : %+v", err)
	}
}
