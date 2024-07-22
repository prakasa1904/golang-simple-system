package whatsapp

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type Repository struct {
	config *viper.Viper
	log    *logrus.Logger
	wa     *whatsmeow.Client
	qr     chan string
}

// Library https://github.com/tulir/whatsmeow
func NewRepository(config *viper.Viper, log *logrus.Logger) *Repository {
	dbLog := waLog.Stdout("Database", config.GetString("whatsapp.debug_level"), true)

	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New(
		config.GetString("whatsapp.dbdriver"),
		"file:"+config.GetString("whatsapp.dbname")+"?_pragma=foreign_keys(1)",
		dbLog,
	)

	if err != nil {
		log.Fatalf("Connect to whatsapp db error : %+v", err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		log.Fatalf("Get device whatsapp error : %+v", err)
	}

	clientLog := waLog.Stdout("Client", config.GetString("whatsapp.debug_level"), true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	qr := make(chan string, 1)

	return &Repository{
		config: config,
		log:    log,
		wa:     client,
		qr:     qr,
	}
}

func (r Repository) ListenQRCode() {
	// running non blocking
	go func() {
		if !r.wa.IsConnected() {
			if r.wa.Store.ID == nil {
				// No ID stored, new login
				qrChan, _ := r.wa.GetQRChannel(context.Background())

				err := r.wa.Connect()
				if err != nil {
					r.log.Fatalf("Get QR code whatsapp error : %+v", err)
				}

				for evt := range qrChan {
					if evt.Event == "code" {
						// Render the QR code here
						// e.g. qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
						// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
						r.log.Warnf("New code incoming:", evt.Code)

						r.qr <- evt.Code
					} else {
						r.log.Warnf("Whatsapp login event : %s", evt.Event)
					}
				}
			} else {
				// Already logged in, just connect
				err := r.wa.Connect()
				if err != nil {
					r.log.Warnf("Login status whatsapp error : %+v", err)
				}
			}
		}

		// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		r.wa.Disconnect()
	}()
}

func (r Repository) GetQRCode() string {
	qrCodeResult := <-r.qr

	return qrCodeResult
}
