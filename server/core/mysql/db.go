package mysql

import (
	"database/sql"
	"fmt"
	"os"

	// mysql database driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/michaelg9/ISOC/server/services/models"
)

const (
	user = "treigerm"
	pwd  = "123"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:3306)/mobile_data?parseTime=true", user, pwd, os.Getenv("DB_HOST")))
	if err != nil {
		// TODO: Error handling
		panic(err)
	}
}

// GetHashedPassword gets the hash of the password from a given user
func GetHashedPassword(username string) (hash string, err error) {
	stmt, err := db.Prepare(passwordQuery)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&hash)
	if err != nil {
		return "", err
	}

	return
}

// InsertBatteryData inserts the given data for the battery status
// TODO: Change timestamp type to time.Time
func InsertBatteryData(deviceID, batteryStatus int, timestamp string) (err error) {
	stmt, err := db.Prepare(insertIntoData)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(deviceID, timestamp)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	stmt, err = db.Prepare(insertIntoBattery)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id, batteryStatus)
	if err != nil {
		return err
	}

	return nil
}

// GetBatteryData gets all battery data from a given API key
func GetBatteryData(apiKey string) (*[]models.Battery, error) {
	stmt, err := db.Prepare(getAllBattery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(apiKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var batteries []models.Battery
	for rows.Next() {
		var b models.Battery
		err = rows.Scan(&b.Time, &b.Value)
		if err != nil {
			return nil, err
		}
		batteries = append(batteries, b)
	}

	return &batteries, nil
}

// GetDeviceData gets all the devices from a user with
// given API key
func GetDeviceData(apiKey string) (*[]models.Device, error) {
	stmt, err := db.Prepare(getAllDevices)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(apiKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []models.Device
	for rows.Next() {
		var d models.Device
		err = rows.Scan(&d.ID, &d.Manufacturer, &d.Model, &d.OS)
		if err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}

	return &devices, nil
}
