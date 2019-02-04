package sdplugin

import "encoding/json"

// BaseEventMessage contains only the event field.
// Used to decied how to serialize the message.
type BaseEventMessage struct {
	Event string `json:"event"`
}

// Coordinates on stream deck
type Coordinates struct {
	Column int `json:"column"`
	Row    int `json:"row"`
}

// KeyEventMessage contains all data from keyDown and keyUp events
type KeyEventMessage struct {
	Event   string     `json:"event"`
	Action  string     `json:"action"`
	Context string     `json:"context"`
	Device  string     `json:"device"`
	Payload KeyPayload `json:"payload"`
}

// KeyPayload data in KeyEventMessage
type KeyPayload struct {
	Settings         json.RawMessage `json:"settings"`
	Coordinates      Coordinates     `json:"coordinates"`
	State            int             `json:"state"`
	UserDesiredState int             `json:"userDesiredState"`
	IsInMultiAction  bool            `json:"isInMultiAction"`
}

// AppearanceEventMessage contains all data from willAppear and willDisappear events
type AppearanceEventMessage struct {
	Event   string            `json:"event"`
	Action  string            `json:"action"`
	Context string            `json:"context"`
	Device  string            `json:"device"`
	Payload AppearancePayload `json:"payload"`
}

// AppearancePayload data in AppearanceEventMessage
type AppearancePayload struct {
	Settings        json.RawMessage `json:"settings"`
	Coordinates     Coordinates     `json:"coordinates"`
	State           int             `json:"state"`
	IsInMultiAction bool            `json:"isInMultiAction"`
}

// SendToPluginEventMessage contains data send from PropertyInspector
type SendToPluginEventMessage struct {
	Event   string          `json:"event"`
	Action  string          `json:"action"`
	Context string          `json:"context"`
	Payload json.RawMessage `json:"payload"`
}

//TitleParametersDidChangeEventMessage contains all data from titleParametersDidChange event
type TitleParametersDidChangeEventMessage struct {
	Event   string                          `json:"event"`
	Action  string                          `json:"action"`
	Context string                          `json:"context"`
	Device  string                          `json:"device"`
	Payload TitleParametersDidChangePayload `json:"payload"`
}

// TitleParametersDidChangePayload contains all payload data for titleParametersDidChange event
type TitleParametersDidChangePayload struct {
	Coordinates     Coordinates     `json:"coordinates"`
	Settings        json.RawMessage `json:"settings"`
	State           int             `json:"state"`
	Title           string          `json:"title"`
	TitleParameters TitleParameters `json:"titleParameters"`
}

// TitleParameters contains all title formating information
type TitleParameters struct {
	FontFamily     string `json:"fontFamily"`
	FontSize       int    `json:"fontSize"`
	FontStyle      string `json:"fontStyle"`
	FontUnderline  bool   `json:"fontUnderline"`
	ShowTitle      bool   `json:"showTitle"`
	TitleAlignment string `json:"titleAlignment"`
	TitleColor     string `json:"titleColor"`
}

//DeviceDidConnectEventMessage contains all data from the deviceDidConnect event
type DeviceDidConnectEventMessage struct {
	Event      string     `json:"event"`
	Device     string     `json:"device"`
	DeviceInfo DeviceInfo `json:"deviceInfo"`
}

//DeviceInfo describes the connected device
type DeviceInfo struct {
	Type int  `json:"type"`
	Size Size `json:"size"`
}

//Size of the connected device
type Size struct {
	Columns int `json:"columns"`
	Rows    int `json:"rows"`
}

//DeviceDidDisconnectEventMessage contains all data from the deviceDidDisconnect event
type DeviceDidDisconnectEventMessage struct {
	Event  string `json:"event"`
	Device string `json:"device"`
}

//ApplicationEventMessage contains all data from the applicationDidLaunch and applicationDidTerminate events
type ApplicationEventMessage struct {
	Event   string             `json:"event"`
	Payload ApplicationPayload `json:"payload"`
}

//ApplicationPayload contains the launched or terminated application
type ApplicationPayload struct {
	Application string `json:"application"`
}
