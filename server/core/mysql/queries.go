package mysql

const (
	passwordQuery     = "SELECT passwordHash FROM User WHERE username = ?"
	insertIntoData    = "INSERT INTO Data (device, timestamp) VALUES (?, ?);"
	insertIntoBattery = "INSERT INTO BatteryStatus VALUES (?, ?);"
	GetBattery        = "SELECT d.timestamp, b.batteryPercentage " +
		"FROM Data d, BatteryStatus b " +
		"WHERE d.device = ? AND d.id = b.id;"
	GetDevices = "SELECT dev.id, dev.manufacturer, dev.modelName, dev.osVersion " +
		"FROM Device dev, User u " +
		"WHERE u.apiKey = ? AND u.uid = dev.user;"
)
