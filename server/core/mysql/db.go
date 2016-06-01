package mysql

import (
	"database/sql"
	"fmt"

	// mysql database driver
	_ "github.com/go-sql-driver/mysql"
)

const (
	user = "treigerm"
	pwd  = "Hip$terSWAG"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%v:%v@/mobile_data?parseTime=true", user, pwd))
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
func GetBatteryData(apiKey string) {
	// TODO: Implement
}
