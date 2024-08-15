package http

import (
	"net/http"

	"github.com/devetek/go-core/render"
	"github.com/prakasa1904/panji-express/internal/services/whatsapp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type WhatsappController struct {
	log          *logrus.Logger
	myRepository *whatsapp.Repository
	view         *render.Engine
}

func NewWhatsappController(
	config *viper.Viper,
	log *logrus.Logger,
	view *render.Engine,
) *WhatsappController {
	// init module repositories
	myRepository := whatsapp.NewRepository(config, log)

	return &WhatsappController{
		log:          log,
		myRepository: myRepository,
		view:         view,
	}
}

func (c *WhatsappController) setHeaderMeta() {
	c.view.Set("title", "WhatsApp - Panji Express")
	c.view.Set("description", "Jasa pengiriman paket Jakarta - Kalimantan same day")
}

func (c *WhatsappController) Home(w http.ResponseWriter, r *http.Request) {
	c.setHeaderMeta()

	// listen WA event in paralel go routine
	// c.myRepository.ListenQRCode()

	// qrCode := c.myRepository.GetQRCode()

	c.view.Set("qrcode", "")

	// render page with template html (ejs)
	err := c.view.HTML(w).Render("views/pages/whatsapp/index.html")
	if err != nil {
		c.log.Warnf("Render error : %+v", err)
	}

}
