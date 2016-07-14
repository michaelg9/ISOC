package models

import "reflect"

// DataOut is the struct for the output data
type DataOut struct {
	Device []Device `xml:"device" json:"devices"`
}

// DeviceData contains all the tracked data of the device
type DeviceData struct {
	Battery []Battery `xml:"battery" json:"battery"`
}

// GetContents returns a slice of pointers to all the data of the device in the struct
func (deviceData *DeviceData) GetContents() []interface{} {
	v := reflect.Indirect(reflect.ValueOf(deviceData))
	contents := make([]interface{}, v.NumField())

	for i := range contents {
		contents[i] = v.Field(i).Addr().Interface()
	}

	return contents
}
