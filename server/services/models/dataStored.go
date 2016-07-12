package models

// Battery is the struct for the battery
// percentage element
type Battery struct {
	Time  string `xml:"time,attr" json:"time"`
	Value int    `xml:",chardata" json:"value"`
}

// GetBattery gets all the battery data from the given device and scans it into the given slice
func (db *DB) GetBattery(device DeviceStored, batteries *[]Battery) error {
	getBatteryQuery := `SELECT timestamp, batteryPercentage FROM BatteryStatus WHERE device = :id;`
	stmt, err := db.PrepareNamed(getBatteryQuery)
	if err != nil {
		return err
	}

	return stmt.Select(batteries, device)
}

// CreateBattery inserts the given data into the database
func (db *DB) CreateBattery(device DeviceStored, batteries *[]Battery) error {
	createBatteryQuery := "INSERT INTO BatteryStatus (timestamp, batteryPercentage, device) VALUES (?, ?, ?);"
	for _, battery := range *batteries {
		_, err := db.Exec(createBatteryQuery, battery.Time, battery.Value, device.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
