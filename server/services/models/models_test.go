package models_test

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/michaelg9/ISOC/server/services/models"
)

var schemaUser = `
CREATE TABLE User (
  uid int(11) NOT NULL AUTO_INCREMENT,
  email varchar(20) NOT NULL,
  passwordHash char(64) NOT NULL,
  apiKey varchar(32) DEFAULT NULL,
  PRIMARY KEY (uid),
  UNIQUE KEY email (email),
  UNIQUE KEY apiKey (apiKey)
);`

var schemaDevice = `
CREATE TABLE Device (
  id int(11) NOT NULL AUTO_INCREMENT,
  manufacturer varchar(50) DEFAULT NULL,
  modelName varchar(50) DEFAULT NULL,
  osVersion varchar(50) DEFAULT NULL,
  user int(11) NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (user) REFERENCES User (uid) ON DELETE CASCADE
);`

var schemaBattery = `
CREATE TABLE BatteryStatus (
  id int(11) NOT NULL AUTO_INCREMENT,
  batteryPercentage int(11) NOT NULL,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  device int(11) NOT NULL,
  PRIMARY KEY (id),
  KEY device (device),
  CONSTRAINT BatteryStatus_ibfk_2 FOREIGN KEY (device) REFERENCES Device (id) ON DELETE CASCADE
);`

var insertionUser = `
INSERT INTO User VALUES (1,'user@usermail.com','$2a$10$539nT.CNbxpyyqrL9mro3OQEKuAjhTD3UjEa8JYPbZMZEM/HizvxK','37e72ff927f511e688adb827ebf7e157');
`

var insertionDevice = `
INSERT INTO Device VALUES (1,'Motorola','Moto X (2nd Generation)','Android 5.0',1);
`

var insertionBattery = `
INSERT INTO BatteryStatus VALUES (1,70,'2016-05-31 11:48:48',1);
`

var destroyUser = `
DROP TABLE User;
`

var destroyDevice = `
DROP TABLE Device;
`

var destroyBattery = `
DROP TABLE BatteryStatus;
`

func setup() (*models.DB, error) {
	db, err := models.NewDB("treigerm:Hip$terSWAG@/test_db")
	if err != nil {
		return &models.DB{}, err
	}
	db.MustExec(schemaUser)
	db.MustExec(schemaDevice)
	db.MustExec(schemaBattery)
	db.MustExec(insertionUser)
	db.MustExec(insertionDevice)
	db.MustExec(insertionBattery)
	return db, nil
}

func cleanUp(db *models.DB) {
	db.MustExec(destroyBattery)
	db.MustExec(destroyDevice)
	db.MustExec(destroyUser)
}

func TestGetUser(t *testing.T) {
	db, err := setup()
	if err != nil {
		t.Errorf("\n...error on db setup = %v", err.Error())
	}
	defer cleanUp(db)

	user := models.User{Email: "user@usermail.com"}
	expected := models.User{
		ID:           1,
		Email:        "user@usermail.com",
		PasswordHash: "$2a$10$539nT.CNbxpyyqrL9mro3OQEKuAjhTD3UjEa8JYPbZMZEM/HizvxK",
		APIKey:       "37e72ff927f511e688adb827ebf7e157",
	}
	result, err := db.GetUser(user)
	if err != nil {
		t.Errorf("\n...error = %v", err.Error())
	} else if expected != result {
		t.Errorf("\n...exptected = %v\n...obtained = %v", expected, result)
	}
}

func TestGetDevicesFromUser(t *testing.T) {
	db, err := setup()
	if err != nil {
		t.Errorf("\n...error on db setup = %v", err.Error())
	}
	defer cleanUp(db)

	user := models.User{Email: "user@usermail.com"}
	device := models.DeviceStored{
		ID:           1,
		Manufacturer: "Motorola",
		Model:        "Moto X (2nd Generation)",
		OS:           "Android 5.0",
	}
	expected := []models.DeviceStored{device}

	result, err := db.GetDevicesFromUser(user)
	if err != nil {
		t.Errorf("\n...error = %v", err.Error())
	} else if !reflect.DeepEqual(result, expected) {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, result)
	}
}

func TestGetBattery(t *testing.T) {
	// TODO: Implement
}

func TestCreateUser(t *testing.T) {
	// TODO: Implement
}

func TestCreateDeviceForUser(t *testing.T) {
	db, err := setup()
	if err != nil {
		t.Errorf("\n...error on db setup = %v", err.Error())
	}
	defer cleanUp(db)

	user := models.User{
		ID:           1,
		Email:        "user@usermail.com",
		PasswordHash: "$2a$10$539nT.CNbxpyyqrL9mro3OQEKuAjhTD3UjEa8JYPbZMZEM/HizvxK",
		APIKey:       "37e72ff927f511e688adb827ebf7e157",
	}
	newDevice := models.DeviceStored{
		Manufacturer: "One Plus",
		Model:        "Three",
		OS:           "Android 6.0",
	}

	err = db.CreateDeviceForUser(user, newDevice)
	if err != nil {
		t.Errorf("\n...error = %v", err.Error())
	}

	oldDevice := models.DeviceStored{
		ID:           1,
		Manufacturer: "Motorola",
		Model:        "Moto X (2nd Generation)",
		OS:           "Android 5.0",
	}
	expected := []models.DeviceStored{oldDevice, newDevice}
	result, _ := db.GetDevicesFromUser(user)
	if reflect.DeepEqual(result, expected) {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, result)
	}
}

func TestCreateBattery(t *testing.T) {
	// TODO: Implement
}

func TestDeleteUser(t *testing.T) {
	db, err := setup()
	if err != nil {
		t.Errorf("\n...error on db setup = %v", err.Error())
	}
	defer cleanUp(db)

	user := models.User{
		ID:           1,
		Email:        "user@usermail.com",
		PasswordHash: "$2a$10$539nT.CNbxpyyqrL9mro3OQEKuAjhTD3UjEa8JYPbZMZEM/HizvxK",
		APIKey:       "37e72ff927f511e688adb827ebf7e157"}
	err = db.DeleteUser(user)
	if err != nil {
		t.Errorf("\n...error = %v", err.Error())
	}

	expectedErr := sql.ErrNoRows
	_, err = db.GetUser(user)
	if err != expectedErr {
		t.Errorf("\n...error = %v", err.Error())
	}
}

func TestDeleteDevice(t *testing.T) {
	db, err := setup()
	if err != nil {
		t.Errorf("\n...error on db setup = %v", err.Error())
	}
	defer cleanUp(db)

	device := models.DeviceStored{
		ID:           1,
		Manufacturer: "Motorola",
		Model:        "Moto X (2nd Generation)",
		OS:           "Android 5.0",
	}
	err = db.DeleteDevice(device)
	if err != nil {
		t.Errorf("\n...error = %v", err.Error())
	}

	result, err := db.GetDevicesFromUser(models.User{Email: "user@usermail.com"})
	if err == nil && len(result) != 0 {
		t.Errorf("\n...expected error but got = %v", result)
	} else if err != nil {
		t.Errorf("\n...error = %v", err.Error())
	}
}
