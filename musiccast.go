package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/LouisChrist/streamdeck-musiccast/sdplugin"
)

//Settings data for plugin
type Settings struct {
	IP string `json:"IP"`
}

//loadSettings from data or return erro if failed
func loadSettings(payload []byte) (Settings, error) {
	var settings Settings
	err := json.Unmarshal(payload, &settings)
	if err != nil {
		return Settings{}, err
	}
	return settings, nil
}

// contextUpdateWorker keeps the devices power state up to date
func (m *musicCastHandler) contextUpdateWorker(context context.Context, sender sdplugin.Sender, sdContext string) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	log.Println("Starting worker")
	for true {
		select {
		case <-ticker.C:
			log.Println("Worker update")
			m.stateUpdate(sender, sdContext)
		case <-context.Done():
			log.Println("Stopping worker")
			return
		}
	}
}

// stateUpdate fetches status from musiccast device and updates streamdeck icon
func (m *musicCastHandler) stateUpdate(sender sdplugin.Sender, context string) {
	m.contextMapMutex.Lock()
	defer m.contextMapMutex.Unlock()
	if settings, ok := m.contextMap[context]; ok {
		on, err := m.isMusicCastPowerOn(settings.IP)
		if err != nil {
			log.Printf("Could not check device power status: %v\n", err)
			return
		}

		log.Printf("Is on? %v\n", on)

		targetState := 1 //off
		if on {
			targetState = 0 //on
		}
		log.Printf("Settings state to %v\n", targetState)
		err = sender.SetState(context, targetState)
		if err != nil {
			log.Printf("Failed to set device state: %v\n", err)
			return
		}
	}
}

func (m *musicCastHandler) setMusicCastPower(ip string, on bool) error {
	desiredState := "standby"
	if on {
		desiredState = "on"
	}

	resp, err := m.client.Get(fmt.Sprintf("http://%v/YamahaExtendedControl/v1/main/setPower?power=%v", ip, desiredState))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Got wrong status code %v", resp.Status)
	}
	return nil
}

func (m *musicCastHandler) musicCastPowerToggle(ip string) error {
	resp, err := m.client.Get(fmt.Sprintf("http://%v/YamahaExtendedControl/v1/main/setPower?power=toggle", ip))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Got wrong status code %v", resp.Status)
	}
	return nil
}

// musicCastStatus struct for getStatus response from MusicCast device.
// Only needed field are present,
type musicCastStatus struct {
	Code  int    `json:"response_code"`
	Power string `json:"power"` // standby or on
}

func (m *musicCastHandler) isMusicCastPowerOn(ip string) (bool, error) {
	resp, err := m.client.Get(fmt.Sprintf("http://%v/YamahaExtendedControl/v1/main/getStatus", ip))
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("Got wrong status code %v", resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	resp.Body.Close()

	var musicCastStatus musicCastStatus
	err = json.Unmarshal(data, &musicCastStatus)
	if err != nil {
		return false, err
	}

	if musicCastStatus.Code != 0 {
		return false, fmt.Errorf("Error code: %v", musicCastStatus.Code)
	}

	if musicCastStatus.Power == "standby" {
		return false, nil
	}

	return true, nil
}

// propertyInspectorMessageType is used to differentiate between get and startup messages
type propertyInspectorMessageType struct {
	Type string `json:"type"`
}
