package models

import (
	"RestApi/db"
	"RestApi/utils"
	"errors"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *User) Save() error {
	query := `Insert into users (username, password) values (?, ?)`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	hashedPassword := utils.Hasher(u.Password)
	res, execErr := stmt.Exec(u.Username, hashedPassword)
	if execErr != nil {
		return execErr
	}
	u.ID, err = res.LastInsertId()
	return nil
}

func (u *User) ValidateCredential() error {
	query := `Select id, password from users where username = ?`
	row := db.DB.QueryRow(query, u.Username)

	var hashedPassword string
	err := row.Scan(&u.ID, &hashedPassword)
	if err != nil {
		return err
	}

	if !utils.ValidatePassword(u.Password, hashedPassword) {
		return errors.New("invalid Credential")
	}

	return nil
}

func ValidateUserId(id int) error {
	query := `Select id from users where id = ?`
	row := db.DB.QueryRow(query, id)

	var userId int
	err := row.Scan(&userId)
	if err != nil {
		return errors.New("Invalid User")
	}
	return nil
}
