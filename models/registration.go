package models

import (
	"RestApi/db"
	"errors"
)

type Registration struct {
	ID      int64 `json:"id"`
	UserID  int64 `json:"user_id"`
	EventID int64 `json:"event_id"`
}

func CreateRegistration(UserId, EventID int64) *Registration {
	return &Registration{
		UserID:  UserId,
		EventID: EventID,
	}
}

func (regis *Registration) Validate() error {
	err := ValidateUserId(regis.UserID)
	if err != nil {
		return err
	}
	err = ValidateEventId(regis.EventID)
	if err != nil {
		return err
	}
	return nil
}

func (regis *Registration) Save() error {
	query := `Insert into registrations (user_id, event_id) values (?, ?)`
	res, err := db.DB.Exec(query, regis.UserID, regis.EventID)
	if err != nil {
		return errors.New("failed to Save Registration")
	}

	regis.ID, err = res.LastInsertId()
	if err != nil {
		return errors.New("failed to Save Registration")
	}
	return nil
}

func (regis *Registration) Cancel() error {
	query := `Delete from registrations where id = ?`
	_, err := db.DB.Exec(query, regis.ID)
	if err != nil {
		return errors.New("failed to Cancel Registration")
	}
	return nil
}
