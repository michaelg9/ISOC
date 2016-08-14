package models

// Upload is the XML struct for the data transfer
// from client to server
type Upload struct {
	Meta        Meta        `xml:"metadata"`
	TrackedData TrackedData `xml:"device-data"`
}

// Meta is the struct for the Metadata
type Meta struct {
	Device int `xml:"device"` // Device id to identify the sending device
}
