package models

import (
	"RestApi/db"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Event struct {
	ID           int64                `json:"id"`
	Title        string               `json:"title" binding:"required"`
	Description  string               `json:"description" binding:"required"`
	Location     string               `json:"location" binding:"required"`
	Datetime     time.Time            `json:"datetime" binding:"required"`
	UserID       int                  `json:"user_id"`
	Registration []registrationFormat `json:"registration"`
}

type registrationFormat map[string]interface{}

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

func GetEventWithRegistration(eventId int64) (*Event, error) {
	var allRegistration []registrationFormat
	var event *Event
	query := `
		SELECT e.id, e.name, e.description, e.location, e.datetime, e.user_id,
		       r.user_id AS registration_user_id, 
		       r.registration_date AS registration_date
		FROM events e
		LEFT JOIN registrations r ON e.id = r.event_id
		WHERE e.id = ?
		`

	rows, err := db.DB.Query(query, eventId)
	if err != nil {
		return nil, fmt.Errorf("error getting event with registration: %v", err)
	}

	for rows.Next() {
		var (
			eventId, ownerId            int64
			name, description, location string
			eventTime                   time.Time
			registerUsr                 sql.NullInt64
			registerDate                sql.NullTime
		)

		err = rows.Scan(&eventId, &name, &description, &location, &eventTime, &ownerId, &registerUsr, &registerDate)
		if err != nil {
			panic(err)
		}

		if event == nil {
			event = &Event{
				ID:           eventId,
				Title:        name,
				Description:  description,
				Location:     location,
				Datetime:     eventTime,
				UserID:       int(ownerId),
				Registration: allRegistration,
			}
		}

		if registerUsr.Valid && registerDate.Valid {
			allRegistration = append(allRegistration, registrationFormat{
				"user_id":       registerUsr.Int64,
				"register_date": registerDate.Time,
			})
		}
	}

	if event == nil {
		return nil, errors.New("event not found")
	}

	event.Registration = allRegistration
	return event, nil
}
