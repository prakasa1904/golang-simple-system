package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/devetek/go-core/render"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/prakasa1904/panji-express/internal/helper"
	"github.com/prakasa1904/panji-express/internal/model"
	"github.com/prakasa1904/panji-express/internal/services/group"
	"github.com/prakasa1904/panji-express/internal/services/member"
	"github.com/prakasa1904/panji-express/internal/services/order"
	"github.com/prakasa1904/panji-express/internal/services/whatsapp"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AdminOrderController struct {
	log           *logrus.Logger
	view          *render.Engine
	groupUsecase  *group.UseCase
	memberUsecase *member.UseCase
	myUsecase     *order.UseCase
	waRepository  *whatsapp.Repository
	layout        string
}

func NewAdminOrderController(
	wa *whatsapp.Repository,
	db *gorm.DB,
	log *logrus.Logger,
	view *render.Engine,
	validate *validator.Validate,
	layout string,
) *AdminOrderController {
	// init repositories
	groupRepository := group.NewRepository(log)
	memberUsecaseRepository := member.NewRepository(log)
	myRepository := order.NewRepository(log)

	// run background from setting channel
	// go waRepository.Run()

	// init usecases
	groupUsecase := group.NewUseCase(db, log, validate, groupRepository)
	memberUsecase := member.NewUseCase(db, log, validate, memberUsecaseRepository)
	myUsecase := order.NewUseCase(db, log, validate, myRepository)

	return &AdminOrderController{
		log:           log,
		groupUsecase:  groupUsecase,
		memberUsecase: memberUsecase,
		myUsecase:     myUsecase,
		waRepository:  wa,
		view:          view,
		layout:        layout,
	}
}

func (c *AdminOrderController) Home(w http.ResponseWriter, r *http.Request) {
	filter := helper.ConvertQueryToFilter(r.URL, order.AllowedFilterQuery)
	limit := helper.ConvertQueryToLimit(r.URL)
	order := helper.ConvertQueryToOrder(r.URL, "`order`.`id`")

	searchQuery := r.URL.Query().Get("fullname")

	orders, err := c.myUsecase.Find(r.Context(), filter, limit, order)
	if err != nil {
		c.log.Warnf("Find orders error : %+v", err)
	}

	c.view.Set("pageTitle", "Semua order")
	c.view.Set("search", searchQuery)
	// require to validate because orders is just pointer
	if orders != nil {
		c.view.Set("orders", orders)
	}

	// render page with template html (ejs)
	if helper.IsHTMXRequest(r.Header) {
		err = c.view.HTML(w).RenderClean("views/pages/admin/order/admin-order.html")
		if err != nil {
			c.log.Warnf("Render error : %+v", err)
		}
	} else {
		err = c.view.HTML(w).RenderWithLayout("views/pages/admin/order/admin-order.html", c.layout)
		if err != nil {
			c.log.Warnf("Render error : %+v", err)
		}
	}
}

func (c *AdminOrderController) ComponentList(w http.ResponseWriter, r *http.Request) {
	filter := helper.ConvertQueryToFilter(r.URL, order.AllowedFilterQuery)
	limit := helper.ConvertQueryToLimit(r.URL)
	order := helper.ConvertQueryToOrder(r.URL, "`order`.`id`")

	orders, err := c.myUsecase.Find(r.Context(), filter, limit, order)
	if err != nil {
		c.log.Warnf("Find orders error : %+v", err)
	}

	// require to validate because orders is just pointer
	if orders != nil {
		c.view.Set("orders", orders)
	}

	// render page with template html (ejs)
	err = c.view.HTML(w).RenderClean("views/pages/admin/order/order-item.html")
	if err != nil {
		c.log.Warnf("Render error : %+v", err)
	}
}

func (c *AdminOrderController) ComponentForm(w http.ResponseWriter, r *http.Request) {
	// get action
	action := chi.URLParam(r, "action")
	id := chi.URLParam(r, "id")

	// set all view data, if not set data will use cache and causing invalid data
	c.view.Set("formTitle", strings.ToTitle(action))
	c.view.Set("action", action)
	c.view.Set("order", nil)

	// edit form
	if id != "" {
		order, err := c.myUsecase.GetById(r.Context(), id)
		if err != nil {
			c.log.Warnf("Find orders error : %+v", err)
		}

		// require to validate because orders is just pointer
		if order != nil {
			c.view.Set("order", order)
		}
	}

	// render page with template html (ejs)
	err := c.view.HTML(w).RenderClean("views/pages/admin/order/form-order-content.html")
	if err != nil {
		c.log.Warnf("Render error : %+v", err)
	}
}

func (c *AdminOrderController) ComponentDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// set all view data, if not set data will use cache and causing invalid data
	c.view.Set("order", nil)

	// edit form
	if id != "" {
		order, err := c.myUsecase.GetById(r.Context(), id)
		if err != nil {
			c.log.Warnf("Find orders error : %+v", err)
		}

		// require to validate because orders is just pointer
		if order != nil {
			c.view.Set("order", order)
		}
	}

	// render page with template html (ejs)
	err := c.view.HTML(w).RenderClean("views/pages/admin/order/delete-order-content.html")
	if err != nil {
		c.log.Warnf("Render error : %+v", err)
	}
}

func (c *AdminOrderController) MutationCreate(w http.ResponseWriter, r *http.Request) {
	var frontendResp = model.ResponseMutation{
		Status: "Sukses",
	}
	var payload = new(order.CreatePayload)
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		c.log.Warnf("Json decoder error : %+v", err)
	}

	// file will enable when order require to upload file
	// err := r.ParseForm()
	// fileResume, fileHeader, err := r.FormFile("file");
	// if err != nil {
	// 	c.log.Warnf("Reading file error : %+v", err)
	// }

	neworder, err := c.myUsecase.Create(r.Context(), payload)
	if err != nil {
		frontendResp.Status = "Gagal"
		frontendResp.Message = fmt.Sprintf("Gagal membuat order, %v", err)

		c.log.Warnf("Create order error : %+v", err)
	}

	if neworder != nil {
		// success message
		frontendResp.Message = "Berhasil membuat order"

		// need to improve to get dynamic courier data
		courier, err := c.memberUsecase.GetByGroupName(r.Context(), "Courier")
		if err != nil {
			c.log.Warnf("Failed to send notification to courier : %+v", err)
		}

		message := fmt.Sprintf("Hi %s, ada pesanan pengiriman dokumen baru, dengan detail '%s'", courier.Fullname, neworder.Description)

		// send message to courier
		c.waRepository.SendMessage(courier.Phone, message)
	}

	c.view.Set("toasterTitle", frontendResp.Status)
	c.view.Set("toasterDescription", frontendResp.Message)

	// render page with template html (ejs)
	err = c.view.HTML(w).RenderClean("views/components/toaster/toaster.html")
	if err != nil {
		c.log.Warnf("Render error : %+v", err)
	}
}

func (c *AdminOrderController) MutationUpdate(w http.ResponseWriter, r *http.Request) {
	var frontendResp = model.ResponseMutation{
		Status: "Sukses",
	}
	var payload = new(order.UpdatePayload)
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		c.log.Warnf("Json decoder error : %+v", err)
	}

	updatedorder, err := c.myUsecase.Update(r.Context(), payload)
	if err != nil {
		frontendResp.Status = "Gagal"
		frontendResp.Message = fmt.Sprintf("Gagal memperbaharui order, %v", err)

		c.log.Warnf("Update order error : %+v", err)
	}

	if updatedorder != nil {
		// success message
		frontendResp.Message = "Berhasil memperbaharui order"
	}

	c.view.Set("toasterTitle", frontendResp.Status)
	c.view.Set("toasterDescription", frontendResp.Message)

	// render page with template html (ejs)
	err = c.view.HTML(w).RenderClean("views/components/toaster/toaster.html")
	if err != nil {
		c.log.Warnf("Render error : %+v", err)
	}
}

func (c *AdminOrderController) MutationDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var frontendResp = model.ResponseMutation{
		Status: "Sukses",
	}

	var payload = new(order.DeletePayload)
	if id != "" {
		data, err := c.myUsecase.GetById(r.Context(), id)
		if err != nil {
			c.log.Warnf("Find orders error : %+v", err)
		}

		// require to validate because orders is just pointer
		if data != nil {
			payload.ID = id
		}
	}

	deletedorder, err := c.myUsecase.Delete(r.Context(), payload)
	if err != nil {
		frontendResp.Status = "Gagal"
		frontendResp.Message = fmt.Sprintf("Gagal menghapus order, %v", err)

		c.log.Warnf("Delete order error : %+v", err)
	}

	if deletedorder != nil {
		// success message
		frontendResp.Message = "Berhasil menghapus order"
	}

	c.view.Set("toasterTitle", frontendResp.Status)
	c.view.Set("toasterDescription", frontendResp.Message)

	// render page with template html (ejs)
	err = c.view.HTML(w).RenderClean("views/components/toaster/toaster.html")
	if err != nil {
		c.log.Warnf("Render error : %+v", err)
	}
}
