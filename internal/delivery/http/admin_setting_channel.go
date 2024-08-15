package http

// initial module to interact with whatsApp and Telegram channel

import (
	"log"
	"net/http"

	"github.com/devetek/go-core/render"
	"github.com/go-playground/validator/v10"
	"github.com/prakasa1904/panji-express/internal/helper"
	"github.com/prakasa1904/panji-express/internal/services/whatsapp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type AdminSettingChannelController struct {
	config       *logrus.Logger
	log          *logrus.Logger
	view         *render.Engine
	myRepository *whatsapp.Repository
	layout       string
}

func NewAdminSettingChannelController(
	config *viper.Viper,
	db *gorm.DB,
	log *logrus.Logger,
	view *render.Engine,
	validate *validator.Validate,
	layout string,
) *AdminSettingChannelController {
	// init repositories
	myRepository := whatsapp.NewRepository(config, log)

	// run background from setting channel
	go myRepository.Run()

	return &AdminSettingChannelController{
		log:          log,
		myRepository: myRepository,
		view:         view,
		layout:       layout,
	}
}

func (c *AdminSettingChannelController) Home(w http.ResponseWriter, r *http.Request) {
	c.view.Set("pageTitle", "Channel Setting")
	c.view.Set("qrcode", "")

	// get whatsApp login info
	IsConnected := c.myRepository.IsConnected()
	isLoggedIn := c.myRepository.IsLoggedIn()

	log.Println("isLoggedIn IsConnected")
	log.Println(isLoggedIn, IsConnected)
	log.Println("isLoggedIn IsConnected")

	if isLoggedIn {
		c.myRepository.SendMessage("6287772440255", "Halo courier Nicky Roly, ada pesanan baru dari PT Terpusat untuk Indonesia")
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
