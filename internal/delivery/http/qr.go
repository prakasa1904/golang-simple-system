package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/devetek/go-core/render"
	"github.com/prakasa1904/panji-express/internal/services/qr"
	"github.com/sirupsen/logrus"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/spf13/viper"
)

type QRController struct {
	log          *logrus.Logger
	myRepository *qr.Repository
	view         *render.Engine
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

func (c *QRController) Generate(qrCodeVal simpleQRCode) ([]byte, error) {
	qrCode, err := qrcode.Encode(qrCodeVal.Content, qrCodeVal.RecoveryLevel, qrCodeVal.Size)
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

	qrCodeVal := simpleQRCode{Content: content, RecoveryLevel: 1, Size: qrCodeSize}
	codeData, err = c.Generate(qrCodeVal)
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
