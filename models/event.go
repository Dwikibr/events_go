package models

import (
	"RestApi/db"
	"database/sql"
	"fmt"
	"time"
)

// Event represents an event with its details and registrations.
type Event struct {
	ID           int64                `json:"id"`
	Title        string               `json:"title" binding:"required"`
	Description  string               `json:"description" binding:"required"`
	Location     string               `json:"location" binding:"required"`
	Datetime     time.Time            `json:"datetime" binding:"required"`
	UserID       int                  `json:"user_id"`
	Registration []registrationFormat `json:"registration"`
}

// registrationFormat represents the format of a registration.
type registrationFormat map[string]interface{}

// Save inserts a new event into the database.
func (e *Event) Save() error {
	query := `
		INSERT INTO events (name, description, location, dateTime, user_id)
		VALUES (?, ?, ?, ?, ?)
	`
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

// GetAllEvents retrieves all events from the database.
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

// GetEventById retrieves an event by its ID from the database.
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

// Update modifies an existing event in the database.
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

// Delete removes an event from the database.
func (e *Event) Delete() error {
	query := "Delete from events where id = ?"
	_, err := db.DB.Exec(query, e.ID)
	if err != nil {
		return err
	}
	return nil
}

// ValidateEventId checks if an event ID exists in the database.
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

// GetEventWithRegistration retrieves an event along with its registrations from the database.
func (e *Event) GetEventWithRegistration() error {
	var allRegistration []registrationFormat
	query := `
		SELECT e.id, e.name, e.description, e.location, e.datetime, e.user_id,
		       r.user_id AS registration_user_id, 
		       r.registration_date AS registration_date
		FROM events e
		LEFT JOIN registrations r ON e.id = r.event_id
		WHERE e.id = ?
		`

	rows, err := db.DB.Query(query, e.ID)
	if err != nil {
		return fmt.Errorf("error getting event with registration: %v", err)
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

		if registerUsr.Valid && registerDate.Valid {
			allRegistration = append(allRegistration, registrationFormat{
				"user_id":       registerUsr.Int64,
				"register_date": registerDate.Time,
			})
		}
	}

	e.Registration = allRegistration
	return nil
}
