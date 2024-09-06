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
	"go.mau.fi/whatsmeow/appstate"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

type waEvent struct {
	event   string
	number  string
	message string
}

type waInfo struct {
	isLoggedIn bool
	qrCode     string
}

type Repository struct {
	config   *viper.Viper
	appLog   *logrus.Logger
	instance *whatsmeow.Client
	event    chan waEvent
	info     *waInfo
}

var logger = createLog()
var storage = createStore()

// TODO: decide if need it or not
var info = &waInfo{
	isLoggedIn: false,
}

// Library https://github.com/tulir/whatsmeow
func NewRepository(config *viper.Viper, log *logrus.Logger) *Repository {

	// init whatsApp instance
	device := getFirstDevice(storage)
	waInstance := createInstance(info, device, logger)
	// init whatsApp instance

	return &Repository{
		config:   config,
		instance: waInstance,
		event:    make(chan waEvent),
		appLog:   log,
		info:     info,
	}
}

func (w *Repository) Connecting() {
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
					// future data used in next release
					w.info.qrCode = evt.Code

				} else {
					log.Printf("QR channel result: %s", evt.Event)
				}
			}
		}(w)
	}
}

func (w *Repository) Run() {
	w.Connecting()

	err := w.instance.Connect()
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return
	}

	// listen incoming event
	for {
		cmd := <-w.event
		go handleCmd(w, cmd)
	}
}

func (w *Repository) Reconnect() {
	err := w.instance.Connect()
	if err != nil {
		log.Printf("Failed to re-connect: %v", err)
		return
	}
}

func (w *Repository) IsConnected() bool {
	return w.instance.IsConnected()
}

func (w *Repository) IsLoggedIn() bool {
	return w.instance.IsLoggedIn()
}

func (w *Repository) GetQRCode() string {
	return w.info.qrCode
}

// Send Message To Number in whatsApp, usage SendMessage("628113468822", "Halo")
func (w *Repository) SendMessage(receiver string, message string) {
	w.event <- waEvent{event: "send", number: receiver, message: message}
}

func (w *Repository) Logout() error {
	w.event <- waEvent{event: "logout"}

	return nil
}

// set Log
func createLog() waLog.Logger {
	return waLog.Stdout("Database", "INFO", true)
}

// create store
func createStore() *sqlstore.Container {
	container, err := sqlstore.New(
		"sqlite",
		"file:wa_example.db?_pragma=foreign_keys(1)",
		createLog(),
	)
	if err != nil {
		log.Printf("Failed to create store: %v", err)
		return nil
	}

	return container
}

// get all device
func getAllDevices(store *sqlstore.Container) []*store.Device {
	devices, err := store.GetAllDevices()
	if err != nil {
		log.Printf("Failed to create store: %v", err)
	}

	return devices
}

// get new device
func getNewDevice(store *sqlstore.Container) *store.Device {
	device := store.NewDevice()

	return device
}

// get device by ID
func getDeviceByJID(store *sqlstore.Container, jid types.JID) *store.Device {
	device, err := store.GetDevice(jid)
	if err != nil {
		log.Printf("Failed to create store: %v", err)
	}

	return device
}

func getFirstDevice(store *sqlstore.Container) *store.Device {
	device, err := store.GetFirstDevice()
	if err != nil {
		log.Printf("Failed to get device: %v", err)
	}

	return device
}

// create wa instance
func createInstance(info *waInfo, device *store.Device, waLog waLog.Logger) *whatsmeow.Client {
	var pairRejectChan = make(chan bool, 1)

	whatsApp := whatsmeow.NewClient(device, waLog)
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
		// loggedIn
		info.isLoggedIn = true

		return true
	}

	// consider this section, do we need to listen whatsApp event ?
	whatsApp.AddEventHandler(waHandler(whatsApp))
	// consider this section, do we need to listen whatsApp event ?

	return whatsApp
}

// whatsApp interface handler
// WA event handler, listen all websocket event from whatsApp
func waHandler(wa *whatsmeow.Client) whatsmeow.EventHandler {
	return func(rawEvt interface{}) {
		switch evt := rawEvt.(type) {
		case *events.AppStateSyncComplete:
			if len(wa.Store.PushName) > 0 && evt.Name == appstate.WAPatchCriticalBlock {
				err := wa.SendPresence(types.PresenceAvailable)
				if err != nil {
					log.Printf("Failed to send available presence: %v", err)
				} else {
					log.Printf("Marked self as available on app state sync complete")
				}
			}
		case *events.Connected, *events.PushNameSetting:
			if len(wa.Store.PushName) == 0 {
				return
			}
			// Send presence available when connecting and when the pushname is changed.
			// This makes sure that outgoing messages always have the right pushname.
			err := wa.SendPresence(types.PresenceAvailable)
			if err != nil {
				log.Printf("Failed to send available presence: %v", err)
			} else {
				log.Printf("Marked self as available on connected")
			}
		case *events.Receipt:
			if evt.Type == types.ReceiptTypeRead || evt.Type == types.ReceiptTypeReadSelf {
				log.Printf("%v was read by %s at %s", evt.MessageIDs, evt.SourceString(), evt.Timestamp)
			} else if evt.Type == types.ReceiptTypeDelivered {
				log.Printf("%s was delivered to %s at %s", evt.MessageIDs[0], evt.SourceString(), evt.Timestamp)
			}
		case *events.StreamReplaced:
			log.Println(evt.PermanentDisconnectDescription())
		default:
		}
	}
}

// WA input handler, listen user input and interact with whatsApp
func handleCmd(repo *Repository, cmd waEvent) {
	switch cmd.event {
	case "send":
		msg := &waE2E.Message{Conversation: proto.String(cmd.message)}
		JID := types.NewJID(cmd.number, types.DefaultUserServer)

		resp, err := repo.instance.SendMessage(context.TODO(), JID, msg)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		} else {
			log.Printf("Message sent (server timestamp: %s)", resp.Timestamp)
		}
	case "logout":
		err := repo.instance.Logout()
		if err != nil {
			log.Printf("Error logout: %v", err)
		}
	}
}
