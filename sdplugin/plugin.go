// Package sdplugin implements a simple wrapper around the streamdeck API.
//
// Create an instance of Plugin with New(handler Handler) and provide your own
// Handler implementation.
package sdplugin

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	port          int
	pluginUUID    string
	registerEvent string
	info          string
)

func init() {
	// flags from streamdeck
	flag.IntVar(&port, "port", 8080, "Port for webserver")
	flag.StringVar(&pluginUUID, "pluginUUID", "", "UUID of plugin")
	flag.StringVar(&registerEvent, "registerEvent", "", "Name of register event")
	flag.StringVar(&info, "info", "", "JSON info object from StreamDeck app")
	flag.Parse()
}

// Handler must implement all possible events. Each event gets called on its own goroutine.
// The Sender can bee used to send data back to the streamdeck app.
// Sender is implemented by Plugin itself.
type Handler interface {
	HandleKeyDownEvent(sender Sender, event KeyEventMessage) error
	HandleKeyUpEvent(sender Sender, event KeyEventMessage) error
	HandleWillAppearEvent(sender Sender, event AppearanceEventMessage) error
	HandleWillDisappearEvent(sender Sender, event AppearanceEventMessage) error
	HandleSendToPluginEvent(sender Sender, event SendToPluginEventMessage) error
	HandleTitleParametersDidChangeEvent(sender Sender, event TitleParametersDidChangeEventMessage) error
	HandleDeviceDidConnectEvent(sender Sender, event DeviceDidConnectEventMessage) error
	HandleDeviceDidDisconnectEvent(sender Sender, event DeviceDidDisconnectEventMessage) error
	HandleApplicationDidLaunchEvent(sender Sender, event ApplicationEventMessage) error
	HandleApplicationDidTerminateEvent(sender Sender, event ApplicationEventMessage) error
}

// Plugin communicates with StreamDeck websocket.
// Command line args from StreamDeck are automaticaly parsed.
// All send methods are threadsafe.
// Use New(...) to create an instance.
type Plugin struct {
	connSendMutex *sync.Mutex
	conn          *websocket.Conn
	handler       Handler
}

// New plugin instance. The instance is already registered with the streamdeck app.
func New(handler Handler) (*Plugin, error) {
	// connect to streamdeck app
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://localhost:%v", port), nil)
	if err != nil {
		return nil, err
	}

	// register plugin with app
	registerEventMessage := RegisterEventMessage{
		Event: registerEvent,
		UUID:  pluginUUID,
	}
	conn.WriteJSON(&registerEventMessage)

	return &Plugin{
		connSendMutex: &sync.Mutex{},
		conn:          conn,
		handler:       handler,
	}, nil
}

// Run plugin and receive messages in a loop.
func (p *Plugin) Run() error {
	baseEventMessage := BaseEventMessage{}
	for true {
		// raw message
		messageType, data, err := p.conn.ReadMessage()
		if err != nil {
			return err
		}

		// ignore all but text messages
		if messageType != websocket.TextMessage {
			continue
		}

		// parsed to get event name
		err = json.Unmarshal(data, &baseEventMessage)
		if err != nil {
			return err
		}

		// parse and handle different event types
		if baseEventMessage.Event == "keyDown" || baseEventMessage.Event == "keyUp" {
			var keyEventMessage KeyEventMessage
			err = json.Unmarshal(data, &keyEventMessage)
			if err != nil {
				return err
			}

			if baseEventMessage.Event == "keyDown" {
				go func() {
					err = p.handler.HandleKeyDownEvent(p, keyEventMessage)
					if err != nil {
						log.Println(err)
					}
				}()
			} else {
				go func() {
					err = p.handler.HandleKeyUpEvent(p, keyEventMessage)
					if err != nil {
						log.Println(err)
					}
				}()
			}
		} else if baseEventMessage.Event == "willAppear" || baseEventMessage.Event == "willDisappear" {
			var appearanceEventMessage AppearanceEventMessage
			err = json.Unmarshal(data, &appearanceEventMessage)
			if err != nil {
				return err
			}

			if baseEventMessage.Event == "willAppear" {
				go func() {
					err = p.handler.HandleWillAppearEvent(p, appearanceEventMessage)
					if err != nil {
						log.Println(err)
					}
				}()
			} else {
				go func() {
					err = p.handler.HandleWillDisappearEvent(p, appearanceEventMessage)
					if err != nil {
						log.Println(err)
					}
				}()
			}
		} else if baseEventMessage.Event == "sendToPlugin" {
			var sendToPluginEventMessage SendToPluginEventMessage
			err = json.Unmarshal(data, &sendToPluginEventMessage)
			if err != nil {
				return err
			}

			go func() {
				err = p.handler.HandleSendToPluginEvent(p, sendToPluginEventMessage)
				if err != nil {
					log.Println(err)
				}
			}()
		} else if baseEventMessage.Event == "titleParametersDidChange" {
			var titleParametersDidChangeEventMessage TitleParametersDidChangeEventMessage
			err = json.Unmarshal(data, &titleParametersDidChangeEventMessage)
			if err != nil {
				return err
			}

			go func() {
				err = p.handler.HandleTitleParametersDidChangeEvent(p, titleParametersDidChangeEventMessage)
				if err != nil {
					log.Println(err)
				}
			}()
		} else if baseEventMessage.Event == "deviceDidConnect" {
			var deviceDidConnectEventMessage DeviceDidConnectEventMessage
			err = json.Unmarshal(data, &deviceDidConnectEventMessage)
			if err != nil {
				return err
			}

			go func() {
				err = p.handler.HandleDeviceDidConnectEvent(p, deviceDidConnectEventMessage)
				if err != nil {
					log.Println(err)
				}
			}()
		} else if baseEventMessage.Event == "deviceDidDisconnect" {
			var deviceDidDisconnectEventMessage DeviceDidDisconnectEventMessage
			err = json.Unmarshal(data, &deviceDidDisconnectEventMessage)
			if err != nil {
				return err
			}

			go func() {
				err = p.handler.HandleDeviceDidDisconnectEvent(p, deviceDidDisconnectEventMessage)
				if err != nil {
					log.Println(err)
				}
			}()
		} else if baseEventMessage.Event == "applicationDidLaunch" || baseEventMessage.Event == "applicationDidTerminate" {
			var applicationEventMessage ApplicationEventMessage
			err = json.Unmarshal(data, &applicationEventMessage)
			if err != nil {
				return err
			}

			if baseEventMessage.Event == "applicationDidLaunch" {
				go func() {
					err = p.handler.HandleApplicationDidLaunchEvent(p, applicationEventMessage)
					if err != nil {
						log.Println(err)
					}
				}()
			} else {

				go func() {
					err = p.handler.HandleApplicationDidTerminateEvent(p, applicationEventMessage)
					if err != nil {
						log.Println(err)
					}
				}()
			}
		} else {
			log.Printf("Not handled: %v", baseEventMessage.Event)
		}
	}
	return nil
}

// Close underlying connection
func (p *Plugin) Close() error {
	return p.conn.Close()
}
