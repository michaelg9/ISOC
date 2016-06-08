package models

// Battery is the struct for the battery
// percentage element
type Battery struct {
	Time  string `xml:"time,attr" json:"time"`
	Value int    `xml:",chardata" json:"value"`
}

// DeviceStored is the struct of the stored
// device data
type DeviceStored struct {
	ID           int    `json:"id"`
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	OS           string `json:"os"`
}
