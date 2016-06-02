package mysql

const (
	passwordQuery     = "SELECT passwordHash FROM User WHERE username = ?"
	insertIntoData    = "INSERT INTO Data (device, timestamp) VALUES (?, ?);"
	insertIntoBattery = "INSERT INTO BatteryStatus VALUES (?, ?);"
	getAllBattery     = "SELECT d.timestamp, b.batteryPercentage " +
		"FROM Data d, BatteryStatus b, Device dev, User u " +
		"WHERE u.apiKey = ? AND u.uid = dev.user AND dev.id = d.device " +
		"AND d.id = b.id;"
	getAllDevices = "SELECT dev.id, dev.manufacturer, dev.modelName, dev.osVersion " +
		"FROM Device dev, User u " +
		"WHERE u.apiKey = ? AND u.uid = dev.user;"
)
