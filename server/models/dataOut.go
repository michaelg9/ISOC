package models

// DataOut is the struct for the output data
type DataOut struct {
	Device []Device `xml:"device" json:"devices"`
}

// UserResponse is the response struct for /data/{email}
type UserResponse struct {
	User    User     `xml:"user" json:"user"`
	Devices []Device `xml:"devices" json:"devices"`
}

// LoginResponse is the response struct for /login
type LoginResponse struct {
	AccessToken string `json:"accessToken,omitempty"`
	ID          int    `json:"id,omitempty"` // The user id of the user that is logged in
}
