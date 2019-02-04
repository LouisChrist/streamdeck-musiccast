package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/LouisChrist/streamdeck-musiccast/sdplugin"
)

type musicCastHandler struct {
	contextMapMutex *sync.Mutex
	contextMap      map[string]Settings

	cancelMapMutex *sync.Mutex
	cancelMap      map[string]context.CancelFunc

	client *http.Client
}

//newMusicCastHandler initializes a new musicCastHandler
func newMusicCastHandler() *musicCastHandler {
	return &musicCastHandler{
		contextMapMutex: &sync.Mutex{},
		contextMap:      make(map[string]Settings),
		cancelMapMutex:  &sync.Mutex{},
		cancelMap:       make(map[string]context.CancelFunc),
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (m *musicCastHandler) HandleKeyDownEvent(sender sdplugin.Sender, event sdplugin.KeyEventMessage) error {
	settings, err := loadSettings(event.Payload.Settings)
	if err != nil {
		return err
	}

	err = m.musicCastPowerToggle(settings.IP)
	if err != nil {
		sender.ShowAlert(event.Context)
		return err
	}

	on, err := m.isMusicCastPowerOn(settings.IP)
	if err != nil {
		return err
	}
	log.Printf("Is MusicCast device on? %v", on)

	targetState := 1 // off
	if on {
		targetState = 0 // on
	}
	// delay setState for one second to avoid wrong values, because streamdeck switches button by iteself shortly after press
	go func() {
		time.Sleep(time.Second)
		err = sender.SetState(event.Context, targetState)
		if err != nil {
			return
		}
		log.Printf("State set to: %v", targetState)
	}()

	return nil
}

func (m *musicCastHandler) HandleKeyUpEvent(sender sdplugin.Sender, event sdplugin.KeyEventMessage) error {
	return nil
}

func (m *musicCastHandler) HandleWillAppearEvent(sender sdplugin.Sender, event sdplugin.AppearanceEventMessage) error {
	// load settings once at appearance for UI update
	settings, err := loadSettings(event.Payload.Settings)
	if err != nil {
		return err
	}

	log.Printf("[willAppear]Settings loaded: %#v\n", settings)

	// update settings map
	m.contextMapMutex.Lock()
	m.contextMap[event.Context] = settings
	m.contextMapMutex.Unlock()

	// start background update worker
	context, cancelFunc := context.WithCancel(context.Background())
	go m.contextUpdateWorker(context, sender, event.Context)

	// update cancel func map
	m.cancelMapMutex.Lock()
	m.cancelMap[event.Context] = cancelFunc
	m.cancelMapMutex.Unlock()

	return nil
}

func (m *musicCastHandler) HandleWillDisappearEvent(sender sdplugin.Sender, event sdplugin.AppearanceEventMessage) error {
	// delete settings
	m.contextMapMutex.Lock()
	delete(m.contextMap, event.Context)
	m.contextMapMutex.Unlock()

	m.cancelMapMutex.Lock()
	cancelFunc := m.cancelMap[event.Context]
	m.cancelMapMutex.Unlock()

	cancelFunc()

	return nil
}

func (m *musicCastHandler) HandleSendToPluginEvent(sender sdplugin.Sender, event sdplugin.SendToPluginEventMessage) error {
	var propertyInspectorMessageType propertyInspectorMessageType
	err := json.Unmarshal(event.Payload, &propertyInspectorMessageType)
	if err != nil {
		return err
	}

	// send settings to UI at startup
	if propertyInspectorMessageType.Type == "startup" {
		// get settings from map
		m.contextMapMutex.Lock()
		defer m.contextMapMutex.Unlock()
		if settings, ok := m.contextMap[event.Context]; ok {
			// settings request by property view
			err = sender.SendToPropertyInspector(event.Context, "de.louischrist.musiccast.power", &settings)
			if err != nil {
				return err
			}
		}
	} else {
		// settings update from UI
		settings, err := loadSettings(event.Payload)
		if err != nil {
			return err
		}

		err = sender.SetSettings(event.Context, settings)
		if err != nil {
			return err
		}

		// update settings map
		m.contextMapMutex.Lock()
		m.contextMap[event.Context] = settings
		m.contextMapMutex.Unlock()
	}

	return nil
}

func (m *musicCastHandler) HandleTitleParametersDidChangeEvent(sender sdplugin.Sender, event sdplugin.TitleParametersDidChangeEventMessage) error {
	return nil
}

func (m *musicCastHandler) HandleDeviceDidConnectEvent(sender sdplugin.Sender, event sdplugin.DeviceDidConnectEventMessage) error {
	return nil
}

func (m *musicCastHandler) HandleDeviceDidDisconnectEvent(sender sdplugin.Sender, event sdplugin.DeviceDidDisconnectEventMessage) error {
	return nil
}

func (m *musicCastHandler) HandleApplicationDidLaunchEvent(sender sdplugin.Sender, event sdplugin.ApplicationEventMessage) error {
	return nil
}

func (m *musicCastHandler) HandleApplicationDidTerminateEvent(sender sdplugin.Sender, event sdplugin.ApplicationEventMessage) error {
	return nil
}
