package models

import "reflect"

// DataOut is the struct for the output data
type DataOut struct {
	Device []Device `xml:"devices" json:"devices"`
}

// Device contains all stored information about one device
// TODO: Refactor the first fields into DeviceInfo DeviceStored
type Device struct {
	ID           int        `xml:"id,attr" json:"id"`
	Manufacturer string     `xml:"manufacturer" json:"manufacturer"`
	Model        string     `xml:"model" json:"model"`
	OS           string     `xml:"os" json:"os"`
	Data         DeviceData `xml:"data" json:"data"`
}

// SetDeviceInfo sets the fields which give information about the device
// NOTE: Is this the right place for the function?
func (device *Device) SetDeviceInfo(deviceInfo DeviceStored) {
	device.ID = deviceInfo.ID
	device.Manufacturer = deviceInfo.Manufacturer
	device.Model = deviceInfo.Model
	device.OS = deviceInfo.OS
}

// DeviceData contains all the tracked data of the device
type DeviceData struct {
	Battery []Battery `xml:"battery" json:"battery"`
}

// GetContents returns pointers to all the data of the device in the struct
func (deviceData *DeviceData) GetContents() []interface{} {
	v := reflect.Indirect(reflect.ValueOf(deviceData))
	contents := make([]interface{}, v.NumField())

	for i := range contents {
		contents[i] = v.Field(i).Addr().Interface()
	}

	return contents
}
