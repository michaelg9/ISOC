package models_test

import (
	"database/sql"
	"testing"

	"github.com/michaelg9/ISOC/server/services/models"
)

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

var insertionUser = `
INSERT INTO User VALUES (1,'user@usermail.com','$2a$10$539nT.CNbxpyyqrL9mro3OQEKuAjhTD3UjEa8JYPbZMZEM/HizvxK','37e72ff927f511e688adb827ebf7e157');
`

var destroyUser = `
DROP TABLE User;
`

func setup() (*models.DB, error) {
	db, err := models.NewDB("treigerm:Hip$terSWAG@/test_db")
	if err != nil {
		return &models.DB{}, err
	}
	db.MustExec(schemaUser)
	db.MustExec(insertionUser)
	return db, nil
}

func cleanUp(db *models.DB) {
	db.MustExec(destroyUser)
}

func TestGetUser(t *testing.T) {
	db, err := setup()
	if err != nil {
		t.Errorf("\n...error on db setup = %v", err.Error())
	}

	user := models.User{Email: "user@usermail.com"}
	expected := models.User{
		ID:           1,
		Email:        "user@usermail.com",
		PasswordHash: "$2a$10$539nT.CNbxpyyqrL9mro3OQEKuAjhTD3UjEa8JYPbZMZEM/HizvxK",
		APIKey:       "37e72ff927f511e688adb827ebf7e157",
	}
	result, err := db.GetUser(user)
	if err != nil {
		t.Errorf("\n...error = %v", err.Error())
	} else if expected != result {
		t.Errorf("\n...exptected = %v\n...obtained = %v", expected, result)
	}

	cleanUp(db)
}

func TestDeleteUser(t *testing.T) {
	db, err := setup()
	if err != nil {
		t.Errorf("\n...error on db setup = %v", err.Error())
	}

	user := models.User{
		ID:           1,
		Email:        "user@usermail.com",
		PasswordHash: "$2a$10$539nT.CNbxpyyqrL9mro3OQEKuAjhTD3UjEa8JYPbZMZEM/HizvxK",
		APIKey:       "37e72ff927f511e688adb827ebf7e157"}
	err = db.DeleteUser(user)
	if err != nil {
		t.Errorf("\n...error = %v", err.Error())
	}

	expectedErr := sql.ErrNoRows
	_, err = db.GetUser(user)
	if err != expectedErr {
		t.Errorf("\n...error = %v", err.Error())
	}

	cleanUp(db)
}

func TestCreateUser(t *testing.T) {
	// TODO: Implement
}
