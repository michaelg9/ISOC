package models

// DataIn is the XML struct for the data transfer
// from client to server
type DataIn struct {
	Meta    Meta      `xml:"metadata"`
	Battery []Battery `xml:"battery"`
}

// Battery is the struct for the battery
// percentage element
// TODO: Refactor into new file
type Battery struct {
	Time  string `xml:"time,attr" json:"time"`
	Value int    `xml:",chardata" json:"value"`
}

// Meta is the struct for the Metadata
type Meta struct {
	Device int `xml:"device"` // Device id to identify the sending device
}
