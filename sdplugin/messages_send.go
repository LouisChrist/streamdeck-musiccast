package sdplugin

import "encoding/json"

// SetTitle and SetImage can decide to show the content only in software, only on hardware
// or in both places.
const (
	TargetSoftware = "software" // TargetSoftware will only show in software
	TargetHardware = "hardware" // TargetHardware will only show on streamdeck
	TargetBoth     = "both"     // TargetBoth will show in software and on streamdeck
)

// RegisterEventMessage must be sent as first message. This happens automaticaly in the New(...) function.
type RegisterEventMessage struct {
	Event string `json:"event"`
	UUID  string `json:"uuid"`
}

// SetStateEventMessage contains all data needed to SetState
type SetStateEventMessage struct {
	Event   string          `json:"event"`
	Context string          `json:"context"`
	Payload SetStatePayload `json:"payload"`
}

// SetStatePayload data in SetStateEventMessage
type SetStatePayload struct {
	State int `json:"state"`
}

// SendToPropertyInspectorEventMessage contains data send to PropertyInspector
type SendToPropertyInspectorEventMessage struct {
	Event   string          `json:"event"`
	Action  string          `json:"action"`
	Context string          `json:"context"`
	Payload json.RawMessage `json:"payload"`
}

// SetSettingsEventMessage contains all data needed for SetSettingsEvent
type SetSettingsEventMessage struct {
	Event   string          `json:"event"`
	Context string          `json:"context"`
	Payload json.RawMessage `json:"payload"`
}

// ShowNotifyEventMessage contains all data needed for showAlert and showOk events
type ShowNotifyEventMessage struct {
	Event   string `json:"event"`
	Context string `json:"context"`
}

//SetTitleEventMessage contains all data needed for setTitle event
type SetTitleEventMessage struct {
	Event   string          `json:"event"`
	Context string          `json:"context"`
	Payload SetTitlePayload `json:"payload"`
}

//SetTitlePayload data for SetTitleEventMessage
type SetTitlePayload struct {
	Title string `json:"title"`
	// "software", "hardware" or "both". See Target constants
	Target string `json:"target"`
}

//SetImageEventMessage contains all data needed for setImage event
type SetImageEventMessage struct {
	Event   string          `json:"event"`
	Context string          `json:"context"`
	Payload SetImagePayload `json:"payload"`
}

//SetImagePayload data for SetImageEventMessage
type SetImagePayload struct {
	// Base64 encoded image
	Image string `json:"image"`
	// "software", "hardware" or "both". See Target constants
	Target string `json:"target"`
}

//SwitchToProfileEventMessage contains all data needed for switchToProfile event
type SwitchToProfileEventMessage struct {
	Event   string                 `json:"event"`
	Context string                 `json:"context"`
	Device  string                 `json:"device"`
	Payload SwitchToProfilePayload `json:"payload"`
}

// SwitchToProfilePayload data for SwitchToProfileEventMessage
type SwitchToProfilePayload struct {
	// Name of profile
	Profile string `json:"profile"`
}

//OpenURLEventMessage contains all data needed for openURL event
type OpenURLEventMessage struct {
	Event   string         `json:"event"`
	Payload OpenURLPayload `json:"payload"`
}

// OpenURLPayload data for OpenURLEventMessage
type OpenURLPayload struct {
	URL string `json:"url"`
}
