package models_test

import (
	"database/sql"
	"reflect"
	"testing"

	"golang.org/x/crypto/bcrypt"

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

var users = []models.User{
	models.User{
		ID:           1,
		Email:        "user@usermail.com",
		PasswordHash: "$2a$10$539nT.CNbxpyyqrL9mro3OQEKuAjhTD3UjEa8JYPbZMZEM/HizvxK",
		APIKey:       "37e72ff927f511e688adb827ebf7e157",
	},
	models.User{
		ID:     2,
		Email:  "user@mail.com",
		APIKey: "",
	},
}

var deviceInfos = []models.DeviceStored{
	models.DeviceStored{
		ID:           1,
		Manufacturer: "Motorola",
		Model:        "Moto X (2nd Generation)",
		OS:           "Android 5.0",
	},
	models.DeviceStored{
		ID:           2,
		Manufacturer: "One Plus",
		Model:        "Three",
		OS:           "Android 6.0",
	},
}

var devices = []models.Device{
	models.Device{
		DeviceInfo: deviceInfos[0],
		Data: models.DeviceData{
			Battery: batteryData[:1],
		},
	},
}

var batteryData = []models.Battery{
	models.Battery{
		Value: 70,
		Time:  "2016-05-31 11:48:48",
	},
	models.Battery{
		Value: 71,
		Time:  "2016-05-31 11:50:31",
	},
}

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

func checkDBErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf("\n...error on db setup = %v", err.Error())
	}
}

func cleanUp(db *models.DB) {
	db.MustExec(destroyBattery)
	db.MustExec(destroyDevice)
	db.MustExec(destroyUser)
}

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf("\n...error = %v", err.Error())
	}
}

func checkEqual(t *testing.T, expected, obtained interface{}) {
	if !reflect.DeepEqual(expected, obtained) {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}
}

func TestGetUser(t *testing.T) {
	db, err := setup()
	checkDBErr(t, err)
	defer cleanUp(db)

	user := models.User{Email: "user@usermail.com"}
	expected := users[0]
	result, err := db.GetUser(user)
	checkErr(t, err)
	checkEqual(t, expected, result)
}

func TestGetDevicesFromUser(t *testing.T) {
	db, err := setup()
	checkDBErr(t, err)
	defer cleanUp(db)

	user := models.User{Email: "user@usermail.com"}
	expected := devices

	result, err := db.GetDevicesFromUser(user)
	checkErr(t, err)
	checkEqual(t, expected, result)
}

func TestGetDeviceInfos(t *testing.T) {
	db, err := setup()
	checkDBErr(t, err)
	defer cleanUp(db)

	user := models.User{Email: "user@usermail.com"}
	device := deviceInfos[0]
	expected := []models.DeviceStored{device}

	result, err := db.GetDeviceInfos(user)
	checkErr(t, err)
	checkEqual(t, expected, result)
}

func TestGetBattery(t *testing.T) {
	db, err := setup()
	checkDBErr(t, err)
	defer cleanUp(db)

	device := deviceInfos[0]
	var batteries []models.Battery
	expected := batteryData[:1]
	err = db.GetBattery(device, &batteries)
	checkErr(t, err)
	checkEqual(t, expected, batteries)
}

func TestCreateUser(t *testing.T) {
	db, err := setup()
	checkDBErr(t, err)
	defer cleanUp(db)

	newUser := users[1]
	pwd, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	newUser.PasswordHash = string(pwd)
	err = db.CreateUser(newUser)
	checkErr(t, err)

	result, err := db.GetUser(newUser)
	checkErr(t, err)
	checkEqual(t, newUser, result)
}

func TestCreateDeviceForUser(t *testing.T) {
	db, err := setup()
	checkDBErr(t, err)
	defer cleanUp(db)

	user := users[0]
	newDevice := deviceInfos[1]

	err = db.CreateDeviceForUser(user, newDevice)
	checkErr(t, err)

	oldDevice := deviceInfos[0]
	expected := []models.DeviceStored{oldDevice, newDevice}
	result, err := db.GetDeviceInfos(user)
	checkErr(t, err)
	checkEqual(t, expected, result)
}

func TestCreateBattery(t *testing.T) {
	db, err := setup()
	checkDBErr(t, err)
	defer cleanUp(db)

	device := deviceInfos[0]
	batteriesToInsert := batteryData[1:]
	err = db.CreateBattery(device, &batteriesToInsert)
	checkErr(t, err)

	var batteries []models.Battery
	err = db.GetBattery(device, &batteries)
	expected := batteryData
	checkErr(t, err)
	checkEqual(t, expected, batteries)
}

func TestUpdateUser(t *testing.T) {
	db, err := setup()
	checkDBErr(t, err)
	defer cleanUp(db)

	user := users[0]
	user.Email = "user2@mail.com"
	pwd, err := bcrypt.GenerateFromPassword([]byte("12345"), bcrypt.DefaultCost)
	checkErr(t, err)
	user.PasswordHash = string(pwd)
	user.APIKey = "12345678abcdefg"
	err = db.UpdateUser(user, "Email")
	checkErr(t, err)
	err = db.UpdateUser(user, "PasswordHash")
	checkErr(t, err)
	err = db.UpdateUser(user, "APIKey")
	checkErr(t, err)

	result, err := db.GetUser(user)
	checkErr(t, err)
	checkEqual(t, user, result)
}

func TestUpdateDevice(t *testing.T) {
	db, err := setup()
	checkDBErr(t, err)
	defer cleanUp(db)

	device := deviceInfos[0]
	device.Manufacturer = "Apple"
	device.Model = "iPhone 6"
	device.OS = "iOS 10"
	err = db.UpdateDevice(device, "Manufacturer")
	checkErr(t, err)
	err = db.UpdateDevice(device, "Model")
	checkErr(t, err)
	err = db.UpdateDevice(device, "OS")
	checkErr(t, err)

	expected := []models.DeviceStored{device}
	result, err := db.GetDeviceInfos(users[0])
	checkErr(t, err)
	checkEqual(t, expected, result)
}

func TestDeleteUser(t *testing.T) {
	db, err := setup()
	checkDBErr(t, err)
	defer cleanUp(db)

	user := users[0]
	err = db.DeleteUser(user)
	checkErr(t, err)

	expectedErr := sql.ErrNoRows
	_, err = db.GetUser(user)
	checkEqual(t, expectedErr, err)
}

func TestDeleteDevice(t *testing.T) {
	db, err := setup()
	checkDBErr(t, err)
	defer cleanUp(db)

	device := deviceInfos[0]
	err = db.DeleteDevice(device)
	checkErr(t, err)

	result, err := db.GetDeviceInfos(models.User{Email: "user@usermail.com"})
	if err == nil && len(result) != 0 {
		t.Errorf("\n...expected error but got = %v", result)
	} else if err != nil {
		t.Errorf("\n...error = %v", err.Error())
	}
}
