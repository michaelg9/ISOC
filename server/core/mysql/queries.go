package mysql

// IDEA: make an array of struct which have fields dataStuct, get and insert so you can loop over them if you want to
//       insert new data
const (
	passwordQuery     = "SELECT passwordHash FROM User WHERE username = ?"
	insertIntoData    = "INSERT INTO Data (device, timestamp) VALUES (?, ?);"
	insertIntoBattery = "INSERT INTO BatteryStatus VALUES (?, ?);"
	getUser           = "SELECT uid, username, passwordHash, apiKey FROM User WHERE username = ?;"
	getBattery        = "SELECT d.timestamp, b.batteryPercentage " +
		"FROM Data d, BatteryStatus b " +
		"WHERE d.device = ? AND d.id = b.id;"
	getDevices = "SELECT dev.id, dev.manufacturer, dev.modelName, dev.osVersion " +
		"FROM Device dev, User u " +
		"WHERE u.apiKey = ? AND u.uid = dev.user;"
)
