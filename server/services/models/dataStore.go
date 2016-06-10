package models

// User is the struct of the stored user data
type User struct {
	ID           int
	Username     string
	PasswordHash string
	APIKey       string
}

// DeviceStored is the struct of the stored
// device data
// TODO: Make function which converts this to device struct
type DeviceStored struct {
	ID           int    `json:"id"`
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	OS           string `json:"os"`
}

// Battery is the struct for the battery
// percentage element
// TODO: Change timestamp type to time.Time
type Battery struct {
	Time  string `xml:"time,attr" json:"time"`
	Value int    `xml:",chardata" json:"value"`
}
