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

// SetDeviceInfo sets the fields which give information about the device
// NOTE: Is this the right place for the function?
func (device *Device) SetDeviceInfo(deviceInfo DeviceStored) {
	device.ID = deviceInfo.ID
	device.Manufacturer = deviceInfo.Manufacturer
	device.Model = deviceInfo.Model
	device.OS = deviceInfo.OS
}
