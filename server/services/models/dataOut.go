package models

// DataOut is the struct for the output data
type DataOut struct {
	Device []Device `json:"devices"`
}

// Device contains all stored information about one device
type Device struct {
	ID           int       `json:"id"`
	Manufacturer string    `json:"manufacturer"`
	Model        string    `json:"model"`
	OS           string    `json:"os"`
	Battery      []Battery `json:"battery"`
}
