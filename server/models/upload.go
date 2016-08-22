package models

// Upload is the XML struct for the data transfer
// from client to server
type Upload struct {
	Meta        AboutDevice `xml:"metadata"`
	TrackedData TrackedData `xml:"device-data"`
}
