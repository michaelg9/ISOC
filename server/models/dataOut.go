package models

// DataOut is the struct for the output data
type DataOut struct {
	Device []Device `xml:"device" json:"devices"`
}

// SessionData is the data which is available to the web page during one session
type SessionData struct {
	DataOut
	User User `json:"user"`
}

// UserResponse is the response struct for /data/{email}
type UserResponse struct {
	User    User     `xml:"user" json:"user"`
	Devices []Device `xml:"devices" json:"devices"`
}
