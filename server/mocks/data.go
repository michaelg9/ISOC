package mocks

import "github.com/michaelg9/ISOC/server/models"

// JWT is a sample JWT.
var JWT = "123"

// Users is a slice of sample Users.
var Users = []models.User{
	models.User{
		ID:           1,
		Email:        "user@usermail.com",
		PasswordHash: "$2a$10$539nT.CNbxpyyqrL9mro3OQEKuAjhTD3UjEa8JYPbZMZEM/HizvxK", // Passord: 123456
		APIKey:       "37e72ff927f511e688adb827ebf7e157",
	},
	models.User{
		ID:     2,
		Email:  "user@mail.com",
		APIKey: "",
	},
}

// AboutDevices is a slice of sample AboutDevice structs.
var AboutDevices = []models.AboutDevice{
	models.AboutDevice{
		ID:           1,
		Manufacturer: "Motorola",
		Model:        "Moto X (2nd Generation)",
		OS:           "Android 5.0",
	},
	models.AboutDevice{
		ID:           2,
		Manufacturer: "One Plus",
		Model:        "Three",
		OS:           "Android 6.0",
	},
}

// Devices is a slice of sample Devices.
var Devices = []models.Device{
	models.Device{
		AboutDevice: AboutDevices[0],
		Data: models.TrackedData{
			Battery: BatteryData[:1],
		},
	},
}

// BatteryData is a slice of sample Battery data.
var BatteryData = []models.Battery{
	models.Battery{
		Value: 70,
		Time:  "2016-05-31 11:48:48",
	},
	models.Battery{
		Value: 71,
		Time:  "2016-05-31 11:50:31",
	},
}
