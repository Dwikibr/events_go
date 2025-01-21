package models

import (
	"RestApi/db"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Event struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title" binding:"required"`
	Description  string    `json:"description" binding:"required"`
	Location     string    `json:"location" binding:"required"`
	Datetime     time.Time `json:"datetime" binding:"required"`
	UserID       int       `json:"user_id"`
	Registration []int64   `json:"registration"`
}

func (e *Event) Save() error {
	query := `
		INSERT INTO events (name, description, location, dateTime, user_id)
		VALUES (?, ?, ?, ?, ?)
	`
	//db.DB.Exec(query, e.Title, e.Description, e.Location, e.Datetime, e.UserID)
	//db.DB.Query(query, e.Title, e.Description, e.Location, e.Datetime, e.UserID)
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			fmt.Printf("Error closing statement: %v\n", err)
		}
	}(stmt)

	result, err := stmt.Exec(e.Title, e.Description, e.Location, e.Datetime, e.UserID)
	if err != nil {
		return err
	}
	e.ID, err = result.LastInsertId()

	return err
}

func GetAllEvents() ([]Event, error) {
	query := "SELECT * FROM events"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Printf("Error closing rows: %v\n", err)
		}
	}(rows)

	var events []Event
	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.Location, &e.Datetime, &e.UserID)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func GetEventById(id int64) (*Event, error) {
	query := "Select * from events where id = ?"
	row := db.DB.QueryRow(query, id)
	var event Event
	err := row.Scan(&event.ID, &event.Title, &event.Description, &event.Location, &event.Datetime, &event.UserID)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (e *Event) Update() error {
	query := `
		Update events
		Set name = ?, description = ?, location = ?, dateTime = ?
		Where id = ?
	`
	_, err := db.DB.Exec(query, e.Title, e.Description, e.Location, e.Datetime, e.ID)
	if err != nil {
		return err
	}
	return nil
}

func (e *Event) Delete() error {
	query := "Delete from events where id = ?"
	_, err := db.DB.Exec(query, e.ID)
	if err != nil {
		return err
	}
	return nil
}

func ValidateEventId(id int64) error {
	query := `Select id from events where id = ?`
	row := db.DB.QueryRow(query, id)

	var event Event
	err := row.Scan(&event.ID)
	if err != nil {
		return err
	}
	return nil
}

func GetEventWithRegistration(eventId int64) ([]int64, error) {
	var allRegistration []int64
	query := ``
	res, err := db.DB.Query(query, eventId)
	if err != nil {
		return allRegistration, errors.New("failed to get all registrations")
	}

	for res.Next() {
		var id int64
		scanErr := res.Scan(&id)
		if scanErr != nil {
			return allRegistration, errors.New("failed to scan registration")
		}
		allRegistration = append(allRegistration, id)
	}

	return allRegistration, nil
}
