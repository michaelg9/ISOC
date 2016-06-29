package models

// User is the struct of the stored user data
type User struct {
	ID           int
	Email        string
	PasswordHash string
	APIKey       string
}

// DeviceStored is the struct of the stored
// device data
type DeviceStored struct {
	ID           int    `json:"id"`
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	OS           string `json:"os"`
}

// Battery is the struct for the battery
// percentage element
type Battery struct {
	Time  string `xml:"time,attr" json:"time"`
	Value int    `xml:",chardata" json:"value"`
}
