package models

// TODO: Test time input

import (
	"database/sql"
	"reflect"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

/* Database schema */

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
INSERT INTO User VALUES (1,'user@usermail.com','$2a$10$539nT.CNbxpyyqrL9mro3OQEKuAjhTD3UjEa8JYPbZMZEM/HizvxK','37e72ff927f511e688adb827ebf7e157');
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
	},
	User{
		ID:     2,
		Email:  "user@mail.com",
		APIKey: "",
	},
}

var deviceInfos = []DeviceStored{
	DeviceStored{
		ID:           1,
		IMEI:         "123456789012345",
		Manufacturer: "Motorola",
		Model:        "Moto X (2nd Generation)",
		OS:           "Android 5.0",
	},
	DeviceStored{
		ID:           2,
		IMEI:         "12345678901234567",
		Manufacturer: "One Plus",
		Model:        "Three",
		OS:           "Android 6.0",
	},
}

var devices = []Device{
	Device{
		DeviceInfo: deviceInfos[0],
		Data: DeviceData{
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

func setup() (*DB, error) {
	// Panics if there is a connection error
	db := NewDB("treigerm:Hip$terSWAG@/test_db?")

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

	return db, nil
}

func checkDBErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf("\n...error on db setup = %v", err.Error())
	}
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

	user := User{Email: "user@usermail.com"}
	expected := users[0]
	result, err := db.GetUser(user)
	checkErr(t, err)
	checkEqual(t, expected, result)
}

func TestGetDevicesFromUser(t *testing.T) {
	db, err := setup()
	checkDBErr(t, err)
	defer cleanUp(db)

	user := User{Email: "user@usermail.com"}
	expected := devices

	result, err := db.GetDevicesFromUser(user)
	checkErr(t, err)
	checkEqual(t, expected, result)
}

func TestGetDeviceInfos(t *testing.T) {
	db, err := setup()
	checkDBErr(t, err)
	defer cleanUp(db)

	user := User{Email: "user@usermail.com"}
	expected := deviceInfos[:1]

	result, err := db.getDeviceInfos(user)
	checkErr(t, err)
	checkEqual(t, expected, result)
}

func TestGetData(t *testing.T) {
	db, err := setup()
	checkDBErr(t, err)
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
		err = db.GetData(device, test.ptrToResult)
		checkErr(t, err)
		result := reflect.Indirect(reflect.ValueOf(test.ptrToResult)).Interface()
		checkEqual(t, test.expected, result)
	}
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
	expected := []DeviceStored{oldDevice, newDevice}
	result, err := db.getDeviceInfos(user)
	checkErr(t, err)
	checkEqual(t, expected, result)
}

func TestCreateData(t *testing.T) {
	db, err := setup()
	checkDBErr(t, err)
	defer cleanUp(db)

	toInsertBattery := batteryData[1:]
	toInsertCall := callData[1:]
	toInsertApp := appData[1:]
	toInsertRunservice := runserviceData[1:]

	var tests = []struct {
		ptrToResult interface{}
		toInsert    interface{}
		expected    interface{}
	}{
		{&[]Battery{}, &toInsertBattery, batteryData},
		{&[]Call{}, &toInsertCall, callData},
		{&[]App{}, &toInsertApp, appData},
		{&[]Runservice{}, &toInsertRunservice, runserviceData},
	}

	device := deviceInfos[0]
	for _, test := range tests {
		err = db.CreateData(device, test.toInsert)
		checkErr(t, err)

		err = db.GetData(device, test.ptrToResult)
		checkErr(t, err)
		result := reflect.Indirect(reflect.ValueOf(test.ptrToResult)).Interface()
		checkEqual(t, test.expected, result)
	}
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
	user.APIKey = "1"
	err = db.UpdateUser(user)
	checkErr(t, err)

	result, err := db.GetUser(user)
	user.APIKey = result.APIKey
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
	err = db.UpdateDevice(device)
	checkErr(t, err)

	expected := []DeviceStored{device}
	result, err := db.getDeviceInfos(users[0])
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

	result, err := db.getDeviceInfos(User{Email: "user@usermail.com"})
	if err == nil && len(result) != 0 {
		t.Errorf("\n...expected error but got = %v", result)
	} else if err != nil {
		t.Errorf("\n...error = %v", err.Error())
	}
}
