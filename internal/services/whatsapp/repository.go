package whatsapp

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
	"time"

	_ "github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

type waEvent struct {
	event   string
	number  string
	message string
}

type Repository struct {
	config   *viper.Viper
	appLog   *logrus.Logger
	instance *whatsmeow.Client
	event    chan waEvent
	qrCode   string
}

// Library https://github.com/tulir/whatsmeow
func NewRepository(config *viper.Viper, log *logrus.Logger) *Repository {
	var pairRejectChan = make(chan bool, 1)
	dbLog := waLog.Stdout("Database", "INFO", true)

	container, err := sqlstore.New(
		"sqlite",
		"file:wa_example.db?_pragma=foreign_keys(1)",
		dbLog,
	)
	if err != nil {
		log.Printf("Failed to create store: %v", err)
	}

	device, err := container.GetFirstDevice()
	if err != nil {
		log.Printf("Failed to get device: %v", err)
	}

	whatsApp := whatsmeow.NewClient(device, dbLog)
	var isWaitingForPair atomic.Bool

	whatsApp.PrePairCallback = func(jid types.JID, platform, businessName string) bool {
		isWaitingForPair.Store(true)
		defer isWaitingForPair.Store(false)
		log.Printf("Pairing %s (platform: %q, business name: %q). Type r within 3 seconds to reject pair", jid, platform, businessName)
		select {
		case reject := <-pairRejectChan:
			if reject {
				log.Println("Rejecting pair")
				return false
			}
		case <-time.After(3 * time.Second):
		}
		log.Println("Accepting pair")
		return true
	}

	// consider this section, do we need to listen whatsApp event ?
	whatsApp.AddEventHandler(waHandler)
	// consider this section, do we need to listen whatsApp event ?

	return &Repository{
		config:   config,
		instance: whatsApp,
		event:    make(chan waEvent),
		appLog:   log,
	}
}

func (w *Repository) Run() {
	// get login or not
	ch, err := w.instance.GetQRChannel(context.Background())
	if err != nil {
		// This error means that we're already logged in, so ignore it.
		if !errors.Is(err, whatsmeow.ErrQRStoreContainsID) {
			log.Printf("Failed to get QR channel: %v", err)
		}
	} else {
		go func(waInstance *Repository) {
			for evt := range ch {
				if evt.Event == "code" {
					w.qrCode = evt.Code

				} else {
					log.Printf("QR channel result: %s", evt.Event)
				}
			}
		}(w)
	}

	err = w.instance.Connect()
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return
	}

	// listen incoming event
	for {
		select {
		case cmd := <-w.event:
			go handleCmd(w.instance, cmd)
		}
	}
}

func (w *Repository) IsConnected() bool {
	return w.instance.IsConnected()
}

func (w *Repository) IsLoggedIn() bool {
	return w.instance.IsLoggedIn()
}

func (w *Repository) GetQRCode() string {
	return w.qrCode
}

// Send Message To Number in whatsApp, usage SendMessage("628113468822", "Halo")
func (w *Repository) SendMessage(receiver string, message string) {
	w.event <- waEvent{event: "send", number: receiver, message: message}

}

// whatsApp interface handler
// WA event handler, listen all websocket event from whatsApp
func waHandler(rawEvt interface{}) {
	log.Println("rawEvt Triggered!")
}

// WA input handler, listen user input and interact with whatsApp
func handleCmd(waInstance *whatsmeow.Client, cmd waEvent) {
	switch cmd.event {
	case "send":
		if waInstance.IsLoggedIn() && waInstance.IsConnected() {
			log.Println("cmd.message: ", cmd.message)
			log.Println("cmd.number: ", cmd.number)

			msg := &waE2E.Message{Conversation: proto.String(cmd.message)}
			JID := types.NewJID(cmd.number, types.DefaultUserServer)

			resp, err := waInstance.SendMessage(context.TODO(), JID, msg)
			if err != nil {
				log.Printf("Error sending message: %v", err)
			} else {
				log.Printf("Message sent (server timestamp: %s)", resp.Timestamp)
			}
		}
	case "logout":
		err := waInstance.Logout()
		if err != nil {
			log.Printf("Error logout: %v", err)
		}
	}
}
