package http

// initial module to interact with whatsApp and Telegram channel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/devetek/go-core/render"
	"github.com/go-playground/validator/v10"
	"github.com/prakasa1904/panji-express/internal/helper"
	"github.com/prakasa1904/panji-express/internal/model"
	"github.com/prakasa1904/panji-express/internal/services/whatsapp"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AdminSettingChannelController struct {
	log          *logrus.Logger
	view         *render.Engine
	myRepository *whatsapp.Repository
	layout       string
}

func NewAdminSettingChannelController(
	wa *whatsapp.Repository,
	db *gorm.DB,
	log *logrus.Logger,
	view *render.Engine,
	validate *validator.Validate,
	layout string,
) *AdminSettingChannelController {
	return &AdminSettingChannelController{
		log:          log,
		myRepository: wa,
		view:         view,
		layout:       layout,
	}
}

func (c *AdminSettingChannelController) Home(w http.ResponseWriter, r *http.Request) {
	c.view.Set("pageTitle", "Channel Setting")
	c.view.Set("IsLoggedIn", false)
	c.view.Set("qrcode", "")

	// get whatsApp login info
	IsConnected := c.myRepository.IsConnected()
	isLoggedIn := c.myRepository.IsLoggedIn()

	if isLoggedIn && IsConnected {
		c.view.Set("IsLoggedIn", true)
	} else {
		qrCode := c.myRepository.GetQRCode()
		c.view.Set("qrcode", qrCode)
	}

	// render page with template html (ejs)
	if helper.IsHTMXRequest(r.Header) {
		err := c.view.HTML(w).RenderClean("views/pages/admin/setting-channel/admin-setting-channel.html")
		if err != nil {
			c.log.Warnf("Render error : %+v", err)
		}
	} else {
		err := c.view.HTML(w).RenderWithLayout("views/pages/admin/setting-channel/admin-setting-channel.html", c.layout)
		if err != nil {
			c.log.Warnf("Render error : %+v", err)
		}
	}
}

func (c *AdminSettingChannelController) SendMessage(w http.ResponseWriter, r *http.Request) {
	var frontendResp = model.ResponseMutation{
		Status:  "Sukses",
		Message: "Pesan akan dikirim secara async",
	}
	var payload = new(whatsapp.SendMessagePayload)
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		frontendResp.Message = fmt.Sprintf("Json decoder error : %+v", err)
	}

	c.myRepository.SendMessage(payload.Receiver, payload.Message)

	c.view.Set("toasterTitle", frontendResp.Status)
	c.view.Set("toasterDescription", frontendResp.Message)

	// render page with template html (ejs)
	err = c.view.HTML(w).RenderClean("views/components/toaster/toaster.html")
	if err != nil {
		c.log.Warnf("Render error : %+v", err)
	}
}

func (c *AdminSettingChannelController) Logout(w http.ResponseWriter, r *http.Request) {
	c.view.Set("pageTitle", "Channel Setting")
	c.view.Set("IsLoggedIn", true)
	c.view.Set("qrcode", "")

	// get whatsApp login info
	IsConnected := c.myRepository.IsConnected()
	isLoggedIn := c.myRepository.IsLoggedIn()

	if isLoggedIn || IsConnected {
		c.myRepository.Logout()
		c.view.Set("IsLoggedIn", false)

		time.Sleep(10000)

		qrCode := c.myRepository.GetQRCode()
		c.view.Set("qrcode", qrCode)
	}

	// render page with template html (ejs)
	if helper.IsHTMXRequest(r.Header) {
		err := c.view.HTML(w).RenderClean("views/pages/admin/setting-channel/admin-setting-channel.html")
		if err != nil {
			c.log.Warnf("Render error : %+v", err)
		}
	} else {
		err := c.view.HTML(w).RenderWithLayout("views/pages/admin/setting-channel/admin-setting-channel.html", c.layout)
		if err != nil {
			c.log.Warnf("Render error : %+v", err)
		}
	}
}
