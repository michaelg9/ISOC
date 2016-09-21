package models

// User is the struct of the stored user data
type User struct {
	ID           int    `db:"uid" json:"id" xml:"id"`
	Email        string `db:"email" json:"email" xml:"email"`
	PasswordHash string `db:"passwordHash" json:"-" xml:"-"`
	Admin        bool   `db:"admin" json:"admin" xml:"admin"`
}

// GetAllUsers gets all registered users.
func (db *DB) GetAllUsers() (users []User, err error) {
	getUsersQuery := `SELECT uid, email, passwordHash, admin FROM User;`
	err = db.Select(&users, getUsersQuery)
	return
}

// GetUser gets a user with the specified email
func (db *DB) GetUser(user User) (fullUser User, err error) {
	getUserQuery := `SELECT uid, email, passwordHash, admin FROM User
		             WHERE email = :email OR uid = :uid;`
	stmt, err := db.PrepareNamed(getUserQuery)
	if err != nil {
		return
	}

	err = stmt.Get(&fullUser, user)
	return
}

// CreateUser creates a new user from the given struct
func (db *DB) CreateUser(user User) (insertedID int, err error) {
	insertUserQuery := `INSERT INTO User (email, passwordHash, admin) VALUES (:email, :passwordHash, :admin);`
	result, err := db.NamedExec(insertUserQuery, user)
	if err != nil {
		return
	}

	id, err := result.LastInsertId()
	return int(id), err
}

// UpdateUser update the specified user
func (db *DB) UpdateUser(user User) error {
	// TODO: Decide if user can be made Admin over the API, if yes we need
	//       to make Admin a pointer since otherwise the update mechanism always
	//       updates Admin to false (since false is the boolean default value)
	queries := map[string]string{
		"Email":        `UPDATE User SET email = :email WHERE uid = :uid;`,
		"PasswordHash": `UPDATE User SET passwordHash = :passwordHash WHERE uid = :uid;`,
	}

	// Use custom update method
	return db.update(queries, user)
}

// DeleteUser deletes the user with the information in the user struct.
// Right now also deletes all the devices and its data from the user.
func (db *DB) DeleteUser(user User) error {
	deleteUserQuery := `DELETE FROM User WHERE email = :email OR uid = :uid;`
	_, err := db.NamedExec(deleteUserQuery, user)
	return err
}
