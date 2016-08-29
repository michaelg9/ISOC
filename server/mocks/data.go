package mocks

import "github.com/michaelg9/ISOC/server/models"

// JWT is a sample JWT.
var (
	AccessToken  = "123"
	RefreshToken = "1234"
)

// Users is a slice of sample users.
var Users = []models.User{
	models.User{
		ID:           1,
		Email:        "user@usermail.com",
		PasswordHash: "$2a$10$539nT.CNbxpyyqrL9mro3OQEKuAjhTD3UjEa8JYPbZMZEM/HizvxK", // Passord: 123456
		APIKey:       "37e72ff927f511e688adb827ebf7e157",
		Admin:        true,
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

// Devices is a slice of sample devices.
var Devices = []models.Device{
	models.Device{
		AboutDevice: AboutDevices[0],
		Data: models.Features{
			Battery: BatteryData[:1],
		},
	},
}

// BatteryData is a slice of sample battery data.
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

// Uploads is a slice of sample upload structs.
var Uploads = []models.Upload{
	// Should work
	models.Upload{
		Meta:     models.AboutDevice{ID: 1},
		Features: features,
	},
	// Should fail
	models.Upload{
		Features: features,
	},
}

var features = models.Features{
	Battery: BatteryData,
}
