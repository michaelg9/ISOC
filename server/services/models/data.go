package models

// Data is the XML struct for the data transfer
// from client to server
type Data struct {
	Meta    Meta      `xml:"metadata" json:"metadata"`
	Battery []Battery `xml:"battery" json:"battery"`
}

// Battery is the struct for the battery
// percentage element
type Battery struct {
	Time  string `xml:"time,attr" json:"time"`
	Value int    `xml:",chardata" json:"value"`
}

// Meta is the struct for the Metadata
type Meta struct {
	Device int `xml:"device" json:"device"` // Device id to identify the sending device
}
