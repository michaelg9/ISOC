package mocks

// IDEA: Save data with key as error message

import "github.com/michaelg9/ISOC/server/models"

var (
	AccessToken  = "123"
	RefreshToken = "1234"
)

var Users = []models.User{
	models.User{
		ID:           1,
		Email:        "user@usermail.com",
		PasswordHash: "$2a$10$539nT.CNbxpyyqrL9mro3OQEKuAjhTD3UjEa8JYPbZMZEM/HizvxK", // Passord: 123456
		APIKey:       "37e72ff927f511e688adb827ebf7e157",
		Admin:        true,
	},
	models.User{
		ID:    2,
		Email: "user@mail.com",
	},
}

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

var Devices = []models.Device{
	models.Device{
		AboutDevice: AboutDevices[0],
		Data: models.Features{
			Battery: BatteryData[:1],
		},
	},
}

var SavedFeatures = models.Features{
	Battery:    BatteryData[:1],
	Call:       CallData[:1],
	App:        AppData[:1],
	Runservice: RunserviceData[:1],
}

var features = models.Features{
	Battery:    BatteryData,
	Call:       CallData,
	App:        AppData,
	Runservice: RunserviceData,
}

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

var CallData = []models.Call{
	models.Call{
		Type:    "Outgoing",
		Start:   "2016-05-31 11:50:00",
		End:     "2016-05-31 11:51:00",
		Contact: "43A",
	},
	models.Call{
		Type:    "Ingoing",
		Start:   "2016-06-30 11:50:00",
		End:     "2016-06-30 11:51:00",
		Contact: "43A",
	},
}

var AppData = []models.App{
	models.App{
		Name:      "com.isoc.Monitor",
		UID:       123,
		Version:   "2.3",
		Installed: "2016-05-31 11:50:00",
		Label:     "Monitor",
	},
	models.App{
		Name:      "com.isoc.MobileApp",
		UID:       124,
		Version:   "2.0",
		Installed: "2016-05-31 11:51:00",
		Label:     "MobileApp",
	},
}

var RunserviceData = []models.Runservice{
	models.Runservice{
		AppName: "com.isoc.Monitor",
		RX:      12,
		TX:      10,
		Start:   "2016-05-31 11:50:00",
		End:     "2016-05-31 13:50:00",
	},
	models.Runservice{
		AppName: "com.isoc.Monitor",
		RX:      15,
		TX:      0,
		Start:   "2016-05-31 14:50:00",
		End:     "2016-05-31 15:50:00",
	},
}

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
