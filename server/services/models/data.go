package models

import "time"

// Data is the XML struct for the data transfer
// from client to server
type Data struct {
	Meta    string    `xml:"metadata"`
	Battery []Battery `xml:"battery"`
}

// Battery is the struct for the battery
// percentage element
type Battery struct {
	Time  time.Time `xml:"time,attr"`
	Unit  string    `xml:"unit,attr"`
	Value string    `xml:",chardata"`
}
