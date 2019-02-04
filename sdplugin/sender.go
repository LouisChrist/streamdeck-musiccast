package sdplugin

import (
	"encoding/json"
)

// Sender can send message to the StreamDeck app
type Sender interface {
	SetState(context string, state int) error
	ShowAlert(context string) error
	ShowOk(context string) error
	SetSettings(context string, payload interface{}) error
	SendToPropertyInspector(context string, action string, payload interface{}) error
	SetTitle(context string, title string, target string) error
	SetImage(context string, image string, target string) error
	SwitchToProfile(context string, device string, profile string) error
	OpenURL(url string) error
}

func (p *Plugin) sendMessage(v interface{}) error {
	p.connSendMutex.Lock()
	defer p.connSendMutex.Unlock()
	return p.conn.WriteJSON(v)
}

// SetState of action
func (p *Plugin) SetState(context string, state int) error {
	return p.sendMessage(&SetStateEventMessage{
		Event:   "setState",
		Context: context,
		Payload: SetStatePayload{
			State: state,
		},
	})
}

// ShowAlert on action button
func (p *Plugin) ShowAlert(context string) error {
	return p.sendMessage(&ShowNotifyEventMessage{
		Event:   "showAlert",
		Context: context,
	})
}

// ShowOk on action button
func (p *Plugin) ShowOk(context string) error {
	return p.sendMessage(&ShowNotifyEventMessage{
		Event:   "showOk",
		Context: context,
	})
}

// SetSettings associated to each action
func (p *Plugin) SetSettings(context string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return p.sendMessage(&SetSettingsEventMessage{
		Event:   "setSettings",
		Context: context,
		Payload: data,
	})
}

// SendToPropertyInspector send data to the PropertyInspector
func (p *Plugin) SendToPropertyInspector(context string, action string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return p.sendMessage(&SendToPropertyInspectorEventMessage{
		Event:   "sendToPropertyInspector",
		Context: context,
		Action:  action,
		Payload: data,
	})
}

// SetTitle to new value.
//Target defines if the title should be shown in software, on hardware or in both places.
func (p *Plugin) SetTitle(context string, title string, target string) error {
	return p.sendMessage(&SetTitleEventMessage{
		Event:   "setTitle",
		Context: context,
		Payload: SetTitlePayload{
			Title:  title,
			Target: target,
		},
	})
}

//SetImage for action. The image must be a base64 decoded string.
//Target defines if the image should be shown in software, on hardware or in both places.
func (p *Plugin) SetImage(context string, image string, target string) error {
	return p.sendMessage(&SetImageEventMessage{
		Event:   "setImage",
		Context: context,
		Payload: SetImagePayload{
			Image:  image,
			Target: target,
		},
	})
}

//SwitchToProfile with the given profile name.
func (p *Plugin) SwitchToProfile(context string, device string, profile string) error {
	return p.sendMessage(&SwitchToProfileEventMessage{
		Event:   "switchToProfile",
		Context: context,
		Payload: SwitchToProfilePayload{
			Profile: profile,
		},
	})
}

//OpenURL in default browser.
func (p *Plugin) OpenURL(url string) error {
	return p.sendMessage(&OpenURLEventMessage{
		Event: "openUrl",
		Payload: OpenURLPayload{
			URL: url,
		},
	})
}
