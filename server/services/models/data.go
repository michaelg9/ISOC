package models

// Data is the XML struct for the data transfer
// from client to server
type Data struct {
	Meta    Meta      `xml:"metadata"`
	Battery []Battery `xml:"battery"`
}

// Battery is the struct for the battery
// percentage element
type Battery struct {
	Time  string `xml:"time,attr"`
	Unit  string `xml:"unit,attr"`
	Value int    `xml:",chardata"`
}

// Meta is the struct for the Metadata
type Meta struct {
	Device int `xml:"device"` // Device id to identify the sending device
}
