package models

// Device contains all stored information about one device
type Device struct {
	DeviceInfo DeviceStored `xml:"device-info" json:"deviceInfo"`
	Data       DeviceData   `xml:"device-data" json:"data"`
}

// DeviceStored is the struct of the stored
// device data
type DeviceStored struct {
	ID           int    `json:"id" db:"id"`
	IMEI         string `json:"imei" db:"imei"`
	Manufacturer string `json:"manufacturer" db:"manufacturer"`
	Model        string `json:"model" db:"modelName"`
	OS           string `json:"os" db:"osVersion"`
}

// GetDevicesFromUser gets all the devices and the tracked data from the given user.
func (db *DB) GetDevicesFromUser(user User) (devices []Device, err error) {
	// Query the database for the data which belongs to the API key
	devicesInfo, err := db.GetDeviceInfos(user)
	if err != nil {
		return
	}

	// Convert DeviceOut to Device
	devices = make([]Device, len(devicesInfo))
	for i, d := range devicesInfo {
		devices[i].DeviceInfo = d
	}

	// For each device get all its data and append it to the device
	for i, d := range devices {
		// Get pointers to the arrays which store the tracked data
		// and fill them with the data from the DB
		for _, data := range devices[i].Data.GetContents() {
			err = db.GetData(d.DeviceInfo, data)
			if err != nil {
				return
			}
		}
	}

	return
}

// GetDeviceInfos gets all the registered devices without the stored data from the given user.
func (db *DB) GetDeviceInfos(user User) (devices []DeviceStored, err error) {
	getDevicesQuery := `SELECT dev.id, dev.imei, dev.manufacturer, dev.modelName, dev.osVersion
                	   FROM Device dev, User u
                	   WHERE (u.email = :email OR u.apiKey = :apiKey) AND u.uid = dev.user;`
	stmt, err := db.PrepareNamed(getDevicesQuery)
	if err != nil {
		return
	}

	err = stmt.Select(&devices, user)
	return
}

// CreateDeviceForUser creates the device which was passed as a struct
func (db *DB) CreateDeviceForUser(user User, device DeviceStored) error {
	// TODO: Use named query with interface
	createDeviceQuery := "INSERT INTO Device (imei, manufacturer, modelName, osVersion, user) VALUES (?, ?, ?, ?, ?);"
	_, err := db.Exec(createDeviceQuery, device.IMEI, device.Manufacturer, device.Model, device.OS, user.ID)
	return err
}

// UpdateDevice updates the given field of the device with the given id
// NOTE: What should one be able to update?
func (db *DB) UpdateDevice(device DeviceStored) error {
	queries := map[string]string{
		"Manufacturer": `UPDATE Device SET manufacturer = :manufacturer WHERE id = :id;`,
		"Model":        `UPDATE Device SET modelName = :modelName WHERE id = :id;`,
		"OS":           `UPDATE Device SET osVersion = :osVersion WHERE id = :id;`,
	}

	return db.update(queries, device)
}

// DeleteDevice deletes the device with that specified in the struct.
// Also deletes all the collected data from device.
// IDEA: pass optional value to see if data should be deleted as well
func (db *DB) DeleteDevice(device DeviceStored) error {
	deleteDeviceQuery := `DELETE FROM Device WHERE id = :id`
	_, err := db.NamedExec(deleteDeviceQuery, device)
	return err
}
