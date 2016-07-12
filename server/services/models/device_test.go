package models_test

import (
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
  KEY user (user),
  CONSTRAINT Device_ibfk_1 FOREIGN KEY (user) REFERENCES User (uid)
);`

var insertionUser = `
INSERT INTO User VALUES (1,'user@usermail.com','$2a$10$539nT.CNbxpyyqrL9mro3OQEKuAjhTD3UjEa8JYPbZMZEM/HizvxK','37e72ff927f511e688adb827ebf7e157');
`

var insertionDevice = `
INSERT INTO Device VALUES (1,'Motorola','Moto X (2nd Generation)','Android 5.0',1);
`

var destroyUser = `
DROP TABLE User;
`

var destroyDevice = `
DROP TABLE Device;
`

func setupDeviceDB() (*models.DB, error) {
	db, err := models.NewDB("treigerm:Hip$terSWAG@/test_db")
	if err != nil {
		return &models.DB{}, err
	}
	db.MustExec(schemaUser)
	db.MustExec(schemaDevice)
	db.MustExec(insertionUser)
	db.MustExec(insertionDevice)
	return db, nil
}

func cleanUpDeviceDB(db *models.DB) {
	db.MustExec(destroyDevice)
	db.MustExec(destroyUser)
}

func TestGetDevicesFromUser(t *testing.T) {
	db, err := setupDeviceDB()
	if err != nil {
		t.Errorf("\n...error on db setup = %v", err.Error())
	}
	defer cleanUpDeviceDB(db)

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

func TestCreateDeviceForUser(t *testing.T) {
	db, err := setupDeviceDB()
	if err != nil {
		t.Errorf("\n...error on db setup = %v", err.Error())
	}
	defer cleanUpDeviceDB(db)

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

func TestDeleteDevice(t *testing.T) {
	db, err := setupDeviceDB()
	if err != nil {
		t.Errorf("\n...error on db setup = %v", err.Error())
	}
	defer cleanUpDeviceDB(db)

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
