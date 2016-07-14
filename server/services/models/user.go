package models

// User is the struct of the stored user data
type User struct {
	ID           int    `db:"uid"`
	Email        string `db:"email"`
	PasswordHash string `db:"passwordHash"`
	APIKey       string `db:"apiKey"`
}

// GetUser gets a user with the specified email
func (db *DB) GetUser(user User) (fullUser User, err error) {
	getUserQuery := `SELECT uid, email, passwordHash, COALESCE(apiKey, '') AS apiKey FROM User
		             WHERE email = :email OR apiKey = :apiKey;`
	stmt, err := db.PrepareNamed(getUserQuery)
	if err != nil {
		return
	}

	err = stmt.Get(&fullUser, user)
	return
}

// CreateUser creates a new user from the given struct
func (db *DB) CreateUser(user User) error {
	insertUserQuery := `INSERT INTO User (email, passwordHash) VALUES (:email, :passwordHash);`
	_, err := db.NamedExec(insertUserQuery, user)
	return err
}

// UpdateUser update the specified user
func (db *DB) UpdateUser(user User, field string) error {
	queries := map[string]string{
		"Email":        `UPDATE User SET email = :email WHERE uid = :uid;`,
		"PasswordHash": `UPDATE User SET passwordHash = :passwordHash WHERE uid = :uid;`,
		"APIKey":       `UPDATE User SET apiKey = :apiKey WHERE uid = :uid;`,
	}

	_, err := db.NamedExec(queries[field], user)
	return err
}

// DeleteUser deletes the user with the information in the user struct.
// Right now also deletes all the devices and its data from the user.
func (db *DB) DeleteUser(user User) error {
	deleteUserQuery := `DELETE FROM User WHERE email = :email OR apiKey = :apiKey;`
	_, err := db.NamedExec(deleteUserQuery, user)
	return err
}