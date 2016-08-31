package models

// TODO: Test time input

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var (
	testDBHost     = os.Getenv("TEST_MYSQL_HOST")
	testDBUser     = os.Getenv("TEST_MYSQL_USER")
	testDBPassword = os.Getenv("TEST_MYSQL_PWD")
)

/* Database schema */

var schemaUser = `
CREATE TABLE User (
  uid int(11) NOT NULL AUTO_INCREMENT,
  email varchar(20) NOT NULL,
  passwordHash char(64) NOT NULL,
  apiKey varchar(32) DEFAULT NULL,
  admin tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (uid),
  UNIQUE KEY email (email),
  UNIQUE KEY apiKey (apiKey)
);`

var schemaDevice = `
CREATE TABLE Device (
  id int(11) NOT NULL AUTO_INCREMENT,
  imei varchar(17) DEFAULT NULL,
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

var schemaCall = `
CREATE TABLE ` + "`" + `Call` + "`" + ` (
    id int(11) NOT NULL AUTO_INCREMENT,
    callType varchar(8) NOT NULL,
    ` + "`" + `start` + "`" + ` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ` + "`" + `end` + "`" + ` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    contact varchar(50),
    device int(11) NOT NULL,
    PRIMARY KEY (id),
    KEY device (device),
    CONSTRAINT Call_ibfk_1 FOREIGN KEY (device) REFERENCES Device (id) ON DELETE CASCADE
);`

var schemaApp = `
CREATE TABLE App (
    id int(11) NOT NULL AUTO_INCREMENT,
    name varchar(100) NOT NULL,
    uid int(11) NOT NULL,
    version varchar(10) NOT NULL,
    installed timestamp NOT NULL,
    label varchar(50) NOT NULL,
    device int(11) NOT NULL,
    PRIMARY KEY (id),
    KEY device (device),
    CONSTRAINT App_ibfk_1 FOREIGN KEY (device) REFERENCES Device (id) ON DELETE CASCADE
);`

var schemaRunservice = `
CREATE TABLE Runservice (
    id int(11) NOT NULL AUTO_INCREMENT,
    appName varchar(100) NOT NULL,
    rx int(11) NOT NULL,
    tx int(11) NOT NULL,
    ` + "`" + `start` + "`" + ` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ` + "`" + `end` + "`" + ` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    device int(11) NOT NULL,
    PRIMARY KEY (id),
    KEY device (device),
    CONSTRAINT Runservice_ibfk_1 FOREIGN KEY (device) REFERENCES Device (id) ON DELETE CASCADE
);`

/* Database data */

var insertionUser = `
INSERT INTO User VALUES (1,'user@usermail.com','$2a$10$539nT.CNbxpyyqrL9mro3OQEKuAjhTD3UjEa8JYPbZMZEM/HizvxK','37e72ff927f511e688adb827ebf7e157', TRUE);
`

var insertionDevice = `
INSERT INTO Device VALUES (1,'123456789012345','Motorola','Moto X (2nd Generation)','Android 5.0',1);
`

var insertionBattery = `
INSERT INTO BatteryStatus VALUES (1,70,'2016-05-31 11:48:48',1);
`

var insertionCall = `
INSERT INTO ` + "`" + `Call` + "`" + ` VALUES (1,'Outgoing','2016-05-31 11:50:00','2016-05-31 11:51:00','43A',1);
`

var insertionApp = `
INSERT INTO App VALUES (1,'com.isoc.Monitor',123,'2.3','2016-05-31 11:50:00','Monitor',1);
`

var insertionRunservice = `
INSERT INTO Runservice VALUES (1,'com.isoc.Monitor',12,10,'2016-05-31 11:50:00','2016-05-31 13:50:00',1);
`

/* Database deletion */

var destroyUser = `
DROP TABLE User;
`

var destroyDevice = `
DROP TABLE Device;
`

var destroyBattery = `
DROP TABLE BatteryStatus;
`

var destroyCall = `
DROP TABLE ` + "`" + `Call` + "`" + `;
`

var destroyApp = `
DROP TABLE App;
`

var destroyRunservice = `
DROP TABLE Runservice;
`

var users = []User{
	User{
		ID:           1,
		Email:        "user@usermail.com",
		PasswordHash: "$2a$10$539nT.CNbxpyyqrL9mro3OQEKuAjhTD3UjEa8JYPbZMZEM/HizvxK",
		APIKey:       "37e72ff927f511e688adb827ebf7e157",
		Admin:        true,
	},
	User{
		ID:     2,
		Email:  "user@mail.com",
		APIKey: "",
		Admin:  true,
	},
}

var deviceInfos = []AboutDevice{
	AboutDevice{
		ID:           1,
		IMEI:         "123456789012345",
		Manufacturer: "Motorola",
		Model:        "Moto X (2nd Generation)",
		OS:           "Android 5.0",
	},
	AboutDevice{
		ID:           2,
		IMEI:         "12345678901234567",
		Manufacturer: "One Plus",
		Model:        "Three",
		OS:           "Android 6.0",
	},
}

var devices = []Device{
	Device{
		AboutDevice: deviceInfos[0],
		Data: Features{
			Battery: batteryData[:1],
		},
	},
}

/* Slices for the several data types. The first entry is always the data that's already in the
   database and the second entry is used to test the insert.
*/

var batteryData = []Battery{
	Battery{
		Value: 70,
		Time:  "2016-05-31 11:48:48",
	},
	Battery{
		Value: 71,
		Time:  "2016-05-31 11:50:31",
	},
}

var callData = []Call{
	Call{
		Type:    "Outgoing",
		Start:   "2016-05-31 11:50:00",
		End:     "2016-05-31 11:51:00",
		Contact: "43A",
	},
	Call{
		Type:    "Ingoing",
		Start:   "2016-06-30 11:50:00",
		End:     "2016-06-30 11:51:00",
		Contact: "43A",
	},
}

var appData = []App{
	App{
		Name:      "com.isoc.Monitor",
		UID:       123,
		Version:   "2.3",
		Installed: "2016-05-31 11:50:00",
		Label:     "Monitor",
	},
	App{
		Name:      "com.isoc.MobileApp",
		UID:       124,
		Version:   "2.0",
		Installed: "2016-05-31 11:51:00",
		Label:     "MobileApp",
	},
}

var runserviceData = []Runservice{
	Runservice{
		AppName: "com.isoc.Monitor",
		RX:      12,
		TX:      10,
		Start:   "2016-05-31 11:50:00",
		End:     "2016-05-31 13:50:00",
	},
	Runservice{
		AppName: "com.isoc.Monitor",
		RX:      15,
		TX:      0,
		Start:   "2016-05-31 14:50:00",
		End:     "2016-05-31 15:50:00",
	},
}

func setup() *DB {
	// Panics if there is a connection error
	dsn := fmt.Sprintf("%v:%v@tcp(%v:3306)/test_db?", testDBUser, testDBPassword, testDBHost)
	db := NewDB(dsn)

	// Start transaction
	tx := db.MustBegin()

	// Create tables
	tx.MustExec(schemaUser)
	tx.MustExec(schemaDevice)
	tx.MustExec(schemaBattery)
	tx.MustExec(schemaCall)
	tx.MustExec(schemaApp)
	tx.MustExec(schemaRunservice)

	// Populate tables
	tx.MustExec(insertionUser)
	tx.MustExec(insertionDevice)
	tx.MustExec(insertionBattery)
	tx.MustExec(insertionCall)
	tx.MustExec(insertionApp)
	tx.MustExec(insertionRunservice)

	// Commit transaction
	tx.Commit()

	return db
}

func cleanUp(db *DB) {
	tx := db.MustBegin()
	tx.MustExec(destroyBattery)
	tx.MustExec(destroyCall)
	tx.MustExec(destroyApp)
	tx.MustExec(destroyRunservice)
	tx.MustExec(destroyDevice)
	tx.MustExec(destroyUser)
	tx.Commit()
}

func TestGetUser(t *testing.T) {
	var tests = []struct {
		id       int
		email    string
		key      string
		expected User
	}{
		{0, users[0].Email, "", users[0]},
		{1, "", users[0].APIKey, users[0]},
		{1, users[0].Email, users[0].APIKey, users[0]},
		{0, users[0].Email, "1234", users[0]},
	}
	db := setup()
	defer cleanUp(db)

	for _, test := range tests {
		user := User{ID: test.id, Email: test.email, APIKey: test.key}
		result, err := db.GetUser(user)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, result)
	}
}

func TestGetDeviceFromUser(t *testing.T) {
	db := setup()
	defer cleanUp(db)

	var tests = []struct {
		device      Device
		user        User
		expected    Device
		expectedErr error
	}{
		{devices[0], users[0], devices[0], nil},
		{devices[0], User{ID: 42}, Device{}, sql.ErrNoRows},
	}

	for _, test := range tests {
		result, err := db.GetDeviceFromUser(test.user, test.device)
		if test.expectedErr != nil {
			assert.EqualError(t, err, test.expectedErr.Error())
		} else if assert.NoError(t, err) {
			assert.Equal(t, test.expected, result)
		}
	}
}

func TestGetDevicesFromUser(t *testing.T) {
	db := setup()
	defer cleanUp(db)

	tests := []struct {
		user        User
		expected    []Device
		expectedErr error
	}{
		{User{ID: 1}, devices, nil},
		{User{ID: 42}, []Device{}, nil},
	}

	for _, test := range tests {
		result, err := db.GetDevicesFromUser(test.user)
		if test.expectedErr != nil {
			assert.EqualError(t, err, test.expectedErr.Error())
		} else if assert.NoError(t, err) {
			assert.Equal(t, test.expected, result)
		}
	}
}

func TestGetDeviceInfos(t *testing.T) {
	db := setup()
	defer cleanUp(db)

	tests := []struct {
		user        User
		expected    []AboutDevice
		expectedErr error
	}{
		{User{ID: 1}, deviceInfos[:1], nil},
		{User{ID: 42}, []AboutDevice(nil), nil},
	}

	for _, test := range tests {
		result, err := db.getDeviceInfos(test.user)
		if test.expectedErr != nil {
			t.Log(result)
			assert.EqualError(t, err, test.expectedErr.Error())
		} else if assert.NoError(t, err) {
			assert.Equal(t, test.expected, result)
		}
	}
}

func TestGetFeatureOfDevice(t *testing.T) {
	db := setup()
	defer cleanUp(db)

	var tests = []struct {
		ptrToResult interface{}
		expected    interface{}
	}{
		{&[]Battery{}, batteryData[:1]},
		{&[]Call{}, callData[:1]},
		{&[]App{}, appData[:1]},
		{&[]Runservice{}, runserviceData[:1]},
	}

	device := deviceInfos[0]
	for _, test := range tests {
		err := db.GetFeatureOfDevice(device, test.ptrToResult)
		assert.NoError(t, err)
		result := reflect.Indirect(reflect.ValueOf(test.ptrToResult)).Interface()
		assert.Equal(t, test.expected, result)
	}
}

func TestGetAllUsers(t *testing.T) {
	db := setup()
	defer cleanUp(db)

	result, err := db.GetAllUsers()
	assert.NoError(t, err)
	assert.Equal(t, users[:1], result)
}

func TestGetAllDevices(t *testing.T) {
	db := setup()
	defer cleanUp(db)

	result, err := db.GetAllDevices()
	assert.NoError(t, err)
	assert.Equal(t, deviceInfos[:1], result)
}

func TestGetAllFeatureData(t *testing.T) {
	db := setup()
	defer cleanUp(db)

	var tests = []struct {
		ptrToResult interface{}
		expected    interface{}
	}{
		{&[]Battery{}, batteryData[:1]},
		{&[]Call{}, callData[:1]},
		{&[]App{}, appData[:1]},
		{&[]Runservice{}, runserviceData[:1]},
	}

	for _, test := range tests {
		err := db.GetAllFeatureData(test.ptrToResult)
		assert.NoError(t, err)
		result := reflect.Indirect(reflect.ValueOf(test.ptrToResult)).Interface()
		assert.Equal(t, test.expected, result)
	}
}

func TestCreateUser(t *testing.T) {
	db := setup()
	defer cleanUp(db)

	tests := []struct {
		newUser     User
		password    string
		expectedErr error
	}{
		{users[1], "123456", nil},
	}

	for _, test := range tests {
		pwd, _ := bcrypt.GenerateFromPassword([]byte(test.password), bcrypt.DefaultCost)
		test.newUser.PasswordHash = string(pwd)
		_, err := db.CreateUser(test.newUser)
		if test.expectedErr != nil {
			assert.EqualError(t, test.expectedErr, err.Error())
		} else if assert.NoError(t, err) {
			result, err := db.GetUser(test.newUser)
			assert.NoError(t, err)
			assert.Equal(t, test.newUser, result)
		}
	}
}

func TestCreateDeviceForUser(t *testing.T) {
	db := setup()
	defer cleanUp(db)

	tests := []struct {
		user            User
		newDevice       AboutDevice
		existingDevices []AboutDevice
		expectedErr     error
	}{
		{users[0], deviceInfos[1], deviceInfos[:1], nil},
		{users[0], AboutDevice{ID: 3}, deviceInfos, nil},
	}

	for _, test := range tests {
		_, err := db.CreateDeviceForUser(test.user, test.newDevice)
		if test.expectedErr != nil {
			assert.EqualError(t, test.expectedErr, err.Error())
		} else if assert.NoError(t, err) {
			expected := append(test.existingDevices, test.newDevice)
			result, err := db.getDeviceInfos(test.user)
			assert.NoError(t, err)
			assert.Equal(t, expected, result)
		}
	}
}

func TestCreateFeatureForDevice(t *testing.T) {
	db := setup()
	defer cleanUp(db)

	toInsertBattery := batteryData[1:]
	toInsertCall := callData[1:]
	toInsertApp := appData[1:]
	toInsertRunservice := runserviceData[1:]

	errWrongDeviceID := errors.New("Error 1452: Cannot add or update a child row")

	var tests = []struct {
		ptrToResult interface{}
		toInsert    interface{}
		device      AboutDevice
		expectedErr error
		expected    interface{}
	}{
		{&[]Battery{}, &toInsertBattery, deviceInfos[0], nil, batteryData},
		{&[]Call{}, &toInsertCall, deviceInfos[0], nil, callData},
		{&[]App{}, &toInsertApp, deviceInfos[0], nil, appData},
		{&[]Runservice{}, &toInsertRunservice, deviceInfos[0], nil, runserviceData},
		{&[]Battery{}, &toInsertBattery, deviceInfos[1], errWrongDeviceID, toInsertBattery},
	}

	for _, test := range tests {
		err := db.CreateFeatureForDevice(test.device, test.toInsert)
		if test.expectedErr != nil {
			assert.Contains(t, err.Error(), test.expectedErr.Error())
		} else if assert.NoError(t, err) {
			err = db.GetFeatureOfDevice(test.device, test.ptrToResult)
			assert.NoError(t, err)
			result := reflect.Indirect(reflect.ValueOf(test.ptrToResult)).Interface()
			assert.Equal(t, test.expected, result)
		}
	}
}

func TestUpdateUser(t *testing.T) {
	db := setup()
	defer cleanUp(db)

	tests := []struct {
		newEmail    string
		newPassword string
		newAPIKey   string
		//    isAdmin bool
	}{
		{"user2@mail.com", "12345", "1"},
		{"user@mail.com", "", ""},
		{"", "123456", ""},
		{"", "", "1"},
		{"", "", ""},
	}

	for _, test := range tests {
		oldUser, _ := db.GetUser(User{ID: 1})

		newUser := User{ID: 1, Admin: true}
		newUser.Email = test.newEmail

		if test.newPassword != "" {
			pwd, _ := bcrypt.GenerateFromPassword([]byte(test.newPassword), bcrypt.DefaultCost)
			newUser.PasswordHash = string(pwd)
		} else {
			newUser.PasswordHash = ""
		}

		newUser.APIKey = test.newAPIKey
		err := db.UpdateUser(newUser)
		assert.NoError(t, err)

		result, err := db.GetUser(User{ID: 1})
		assert.NoError(t, err)

		if test.newEmail == "" {
			newUser.Email = oldUser.Email
		}
		if test.newPassword == "" {
			newUser.PasswordHash = oldUser.PasswordHash
		}
		if test.newAPIKey == "" {
			newUser.APIKey = oldUser.APIKey
		} else {
			newUser.APIKey = result.APIKey
		}

		assert.Equal(t, newUser, result)
	}
}

func TestUpdateDevice(t *testing.T) {
	db := setup()
	defer cleanUp(db)

	tests := []struct {
		newOS string
	}{
		{"iOS 10"},
		{""},
	}

	for _, test := range tests {
		oldDevices, _ := db.getDeviceInfos(users[0])
		newDevice := deviceInfos[0]
		newDevice.OS = test.newOS
		err := db.UpdateDevice(newDevice)
		assert.NoError(t, err)

		if test.newOS == "" {
			newDevice.OS = oldDevices[0].OS
		}
		expected := []AboutDevice{newDevice}
		result, _ := db.getDeviceInfos(users[0])
		assert.Equal(t, expected, result)
	}
}

func TestDeleteUser(t *testing.T) {
	db := setup()
	defer cleanUp(db)

	tests := []struct {
		userID int
	}{
		{1},
		{42},
	}

	for _, test := range tests {
		user := User{ID: test.userID}
		err := db.DeleteUser(user)
		assert.NoError(t, err)

		expectedErr := sql.ErrNoRows
		_, err = db.GetUser(user)
		assert.EqualError(t, expectedErr, err.Error())
	}
}

func TestDeleteDevice(t *testing.T) {
	db := setup()
	defer cleanUp(db)

	tests := []struct {
		deviceID int
	}{
		{1},
		{42},
	}

	for _, test := range tests {
		device := AboutDevice{ID: test.deviceID}
		err := db.DeleteDevice(device)
		assert.NoError(t, err)

		expectedErr := sql.ErrNoRows
		_, err = db.GetDeviceFromUser(users[0], Device{AboutDevice: device})
		assert.EqualError(t, expectedErr, err.Error())
	}
}
