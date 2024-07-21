package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/devetek/go-core/render"
	"github.com/devetek/golang-webapp-boilerplate/internal/services/qr"
	"github.com/sirupsen/logrus"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/spf13/viper"
)

type QRController struct {
	log          *logrus.Logger
	myRepository *qr.Repository
	view         *render.Engine
	qrCode       *simpleQRCode
}

type simpleQRCode struct {
	Content       string
	Size          int
	RecoveryLevel qrcode.RecoveryLevel
}

func NewQRController(
	config *viper.Viper,
	log *logrus.Logger,
) *QRController {
	// init module repositories
	myRepository := qr.NewRepository(config, log)

	return &QRController{
		log:          log,
		myRepository: myRepository,
	}
}

func (c *QRController) Generate() ([]byte, error) {
	qrCode, err := qrcode.Encode(c.qrCode.Content, c.qrCode.RecoveryLevel, c.qrCode.Size)
	if err != nil {
		return nil, fmt.Errorf("could not generate a QR code: %v", err)
	}
	return qrCode, nil
}

func (c *QRController) View(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	var size, content string = r.FormValue("size"), r.FormValue("content")
	var codeData []byte

	w.Header().Set("Content-Type", "application/json")

	if content == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(
			"Could not determine the desired QR code content.",
		)
		return
	}

	qrCodeSize, err := strconv.Atoi(size)
	if err != nil || size == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode("Could not determine the desired QR code size.")
		return
	}

	// qrCode := simpleQRCode{Content: content, Size: qrCodeSize}
	c.qrCode.Content = content
	c.qrCode.Size = qrCodeSize

	codeData, err = c.Generate()
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(
			fmt.Sprintf("Could not generate QR code. %v", err),
		)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(codeData)
}
