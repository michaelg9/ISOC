package models

// DataIn is the XML struct for the data transfer
// from client to server
type DataIn struct {
	Meta       Meta       `xml:"metadata"`
	DeviceData DeviceData `xml:"device-data"`
}

// Meta is the struct for the Metadata
// TODO: Use device struct
type Meta struct {
	Device int `xml:"device"` // Device id to identify the sending device
}
