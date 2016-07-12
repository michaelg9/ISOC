package models

import "errors"

// Device contains all stored information about one device
type Device struct {
	DeviceInfo DeviceStored `xml:"device-info" json:"deviceInfo"`
	Data       DeviceData   `xml:"device-data" json:"data"`
}

// DeviceStored is the struct of the stored
// device data
type DeviceStored struct {
	ID           int    `json:"id" db:"id"`
	Manufacturer string `json:"manufacturer" db:"manufacturer"`
	Model        string `json:"model" db:"modelName"`
	OS           string `json:"os" db:"osVersion"`
}

// GetDevicesFromUser gets all the registered devices from the given user
func (db *DB) GetDevicesFromUser(user User) (devices []DeviceStored, err error) {
	getDevicesQuery := `SELECT dev.id, dev.manufacturer, dev.modelName, dev.osVersion
                	   FROM Device dev, User u
                	   WHERE u.email = :email AND u.uid = dev.user;`
	stmt, err := db.PrepareNamed(getDevicesQuery)
	if err != nil {
		return
	}

	err = stmt.Select(&devices, user)
	return
}

// CreateDeviceForUser creates the device which was passed as a struct
func (db *DB) CreateDeviceForUser(user User, device DeviceStored) error {
	createDeviceQuery := "INSERT INTO Device (manufacturer, modelName, osVersion, user) VALUES (?, ?, ?, ?);"
	_, err := db.Exec(createDeviceQuery, device.Manufacturer, device.Model, device.OS, user.ID)
	return err
}

// UpdateDevice updates the given field of the device with the given id
func (db *DB) UpdateDevice(id int, field string, value interface{}) error {
	// TODO: Implement
	return errors.New("Not implemented")
}

// DeleteDevice deletes the device with that specified in the struct
// IDEA: pass optional value to see if data should be deleted as well
func (db *DB) DeleteDevice(device DeviceStored) error {
	deleteDeviceQuery := `DELETE FROM Device WHERE id = :id`
	_, err := db.NamedExec(deleteDeviceQuery, device)
	return err
}
