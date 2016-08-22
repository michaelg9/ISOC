package models

import (
	"reflect"

	"github.com/fatih/structs"
)

// Device contains all stored information about one device
type Device struct {
	AboutDevice AboutDevice `xml:"about-device" json:"aboutDevice"`
	Data        TrackedData `xml:"device-data" json:"data"`
}

// AboutDevice is the struct of the stored
// device data
type AboutDevice struct {
	ID           int    `json:"id" db:"id"`
	IMEI         string `json:"imei" db:"imei"`
	Manufacturer string `json:"manufacturer" db:"manufacturer"`
	Model        string `json:"model" db:"modelName"`
	OS           string `json:"os" db:"osVersion"`
}

// TrackedData contains all the tracked data of the device
type TrackedData struct {
	Battery []Battery `xml:"battery,omitempty" json:"battery,omitempty"`
}

// GetContents returns a slice of pointers to all the data of the device in the struct.
// We need to return pointers in order to preserve the underlying type.
func (deviceData *TrackedData) GetContents() []interface{} {
	v := reflect.Indirect(reflect.ValueOf(deviceData))
	contents := make([]interface{}, v.NumField())

	for i := range contents {
		contents[i] = v.Field(i).Addr().Interface()
	}

	return contents
}

// GetDeviceFromUser gets a device and its tracked data by the Device ID and the user ID.
func (db *DB) GetDeviceFromUser(user User, device Device) (fullDevice Device, err error) {
	getDeviceQuery := `SELECT id, imei, manufacturer, modelName, osVersion
                	   FROM Device
                	   WHERE id = :ID AND user = :userID;`
	stmt, err := db.PrepareNamed(getDeviceQuery)
	if err != nil {
		return
	}

	// Transform aboutDevice struct to map to add user ID
	args := structs.Map(device.AboutDevice)
	args["userID"] = user.ID

	err = stmt.Get(&fullDevice.AboutDevice, args)
	if err != nil {
		return
	}

	for _, data := range fullDevice.Data.GetContents() {
		err = db.GetData(device.AboutDevice, data)
		if err != nil {
			return
		}
	}

	return
}

// GetDevicesFromUser gets all the devices and the tracked data from the given user.
func (db *DB) GetDevicesFromUser(user User) (devices []Device, err error) {
	// Query the database for the data which belongs to the user
	devicesInfo, err := db.getDeviceInfos(user)
	if err != nil {
		return
	}

	// Convert DeviceOut to Device
	devices = make([]Device, len(devicesInfo))
	for i, d := range devicesInfo {
		devices[i].AboutDevice = d
	}

	// For each device get all its data and append it to the device
	for i, d := range devices {
		// Get pointers to the arrays which store the tracked data
		// and fill them with the data from the DB
		for _, data := range devices[i].Data.GetContents() {
			err = db.GetData(d.AboutDevice, data)
			if err != nil {
				return
			}
		}
	}

	return
}

// getDeviceInfos gets all the registered devices without the stored data from the given user.
func (db *DB) getDeviceInfos(user User) (devices []AboutDevice, err error) {
	getDevicesQuery := `SELECT id, imei, manufacturer, modelName, osVersion
                	   FROM Device, User
                	   WHERE (uid = :uid OR apiKey = :apiKey) AND User.uid = Device.user;`
	stmt, err := db.PrepareNamed(getDevicesQuery)
	if err != nil {
		return
	}

	err = stmt.Select(&devices, user)
	return
}

// CreateDeviceForUser creates the device which was passed as a struct
func (db *DB) CreateDeviceForUser(user User, aboutDevice AboutDevice) (insertedID int, err error) {
	createDeviceQuery := `INSERT INTO Device (imei, manufacturer, modelName, osVersion, user)
                          VALUES (:IMEI, :Manufacturer, :Model, :OS, :User);`

	// Transform the device struct into a map[string]interface{} so we can add the user ID
	args := structs.Map(aboutDevice)
	args["User"] = user.ID
	result, err := db.NamedExec(createDeviceQuery, args)
	if err != nil {
		return
	}

	id, err := result.LastInsertId()
	return int(id), err
}

// UpdateDevice updates the given field of the device with the given id
// NOTE: What should one be able to update?
func (db *DB) UpdateDevice(aboutDevice AboutDevice) error {
	queries := map[string]string{
		"Manufacturer": `UPDATE Device SET manufacturer = :manufacturer WHERE id = :id;`,
		"Model":        `UPDATE Device SET modelName = :modelName WHERE id = :id;`,
		"OS":           `UPDATE Device SET osVersion = :osVersion WHERE id = :id;`,
	}

	return db.update(queries, aboutDevice)
}

// DeleteDevice deletes the device with that specified in the struct.
// Also deletes all the collected data from device.
func (db *DB) DeleteDevice(aboutDevice AboutDevice) error {
	deleteDeviceQuery := `DELETE FROM Device WHERE id = :id;`
	_, err := db.NamedExec(deleteDeviceQuery, aboutDevice)
	return err
}
