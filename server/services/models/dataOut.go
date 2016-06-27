package models

import "reflect"

// DataOut is the struct for the output data
type DataOut struct {
	Device []Device `xml:"device" json:"devices"`
}

// Device contains all stored information about one device
type Device struct {
	DeviceInfo DeviceStored `xml:"device-info" json:"deviceInfo"`
	Data       DeviceData   `xml:"device-data" json:"data"`
}

// DeviceData contains all the tracked data of the device
type DeviceData struct {
	Battery []Battery `xml:"battery" json:"battery"`
}

// GetContents returns a slice of pointers to all the data of the device in the struct
// TODO: Check if better to return value and not pointers
func (deviceData *DeviceData) GetContents() []interface{} {
	v := reflect.Indirect(reflect.ValueOf(deviceData))
	contents := make([]interface{}, v.NumField())

	for i := range contents {
		contents[i] = v.Field(i).Addr().Interface()
	}

	return contents
}
