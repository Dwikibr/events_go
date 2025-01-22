package db

import (
	"RestApi/filemanager"
	"RestApi/utils"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type event struct {
	Name        string
	Description string
	Location    string
	Datetime    time.Time
	UserID      int
}

type user struct {
	Username string

	Password string
}

type registration struct {
	EventID          int
	UserID           int
	RegistrationDate time.Time
}

func SeedDB() {
	// Seed Users
	users := readUsers("db/seederData/user.csv")
	insertUsers(users)

	// Seed Events

	events := readEvents("db/seederData/events.csv")
	insertEvents(events)

	// Seed Registrations
	registrations := readRegistrations("db/seederData/registration.csv")
	insertRegistrations(registrations)

	log.Println("Database seeding completed successfully.")
}

func readUsers(filePath string) []user {
	var users []user
	rows := filemanager.ReadCsvFile(filePath)
	for _, row := range rows {
		users = append(users, user{
			Username: row["username"].(string),
			Password: utils.Hasher(row["password"].(string)),
		})
	}
	return users
}

func insertUsers(users []user) {
	query := `INSERT INTO users (username, password) VALUES`
	var placeholders []string
	var values []interface{}
	for _, user := range users {
		placeholders = append(placeholders, "(?, ?)")
		values = append(values, user.Username, user.Password)
	}
	finalQuery := query + strings.Join(placeholders, ",")
	if _, err := DB.Exec(finalQuery, values...); err != nil {
		fmt.Printf("Error seeding user")
	}
}

func readEvents(filePath string) []event {
	var events []event
	rows := filemanager.ReadCsvFile(filePath)
	for _, row := range rows {
		datetime, err := time.Parse("2006-01-02T15:04:05Z", row["Datetime"].(string))
		if err != nil {
			fmt.Printf("Error parsing datetime: %v\n", err)
			return nil
		}
		uid, err := strconv.ParseInt(row["UserID"].(string), 10, 64)
		if err != nil {
			fmt.Printf("Error parsing UserID: %v\n", err)
			return nil
		}
		events = append(events, event{
			Name:        row["Name"].(string),
			Description: row["Description"].(string),
			Location:    row["Location"].(string),
			Datetime:    datetime,
			UserID:      int(uid),
		})
	}
	return events
}

func insertEvents(events []event) {
	query := `INSERT INTO events (name, description, location, dateTime, user_id) VALUES`
	var placeholders []string
	var values []interface{}
	for _, event := range events {
		placeholders = append(placeholders, "(?, ?, ?, ?, ?)")
		values = append(values, event.Name, event.Description, event.Location, event.Datetime, event.UserID)
	}
	finalQuery := query + strings.Join(placeholders, ",")
	if _, err := DB.Exec(finalQuery, values...); err != nil {
		fmt.Printf("Error seeding event")
	}
}

func readRegistrations(filePath string) []registration {
	var registrations []registration
	rows := filemanager.ReadCsvFile(filePath)
	for _, row := range rows {
		eid, err := strconv.ParseInt(row["EventID"].(string), 10, 64)
		if err != nil {
			fmt.Printf("Error parsing EventID: %v\n", err)
			return nil
		}
		uid, err := strconv.ParseInt(row["UserID"].(string), 10, 64)
		if err != nil {
			fmt.Printf("Error parsing UserID: %v\n", err)
			return nil
		}

		datetime, err := time.Parse("2006-01-02T15:04:05Z", row["RegisDate"].(string))
		if err != nil {
			fmt.Printf("Error parsing datetime: %v\n", err)
			return nil
		}
		registrations = append(registrations, registration{
			EventID:          int(eid),
			UserID:           int(uid),
			RegistrationDate: datetime,
		})
	}
	return registrations
}

func insertRegistrations(registrations []registration) {
	query := `INSERT INTO registrations (event_id, user_id, registration_date) VALUES`
	var placeholders []string
	var values []interface{}
	for _, regis := range registrations {
		placeholders = append(placeholders, "(?, ?, ?)")
		values = append(values, regis.EventID, regis.UserID, regis.RegistrationDate)
	}
	finalQuery := query + strings.Join(placeholders, ",")
	if _, err := DB.Exec(finalQuery, values...); err != nil {
		fmt.Printf("Error seeding registration")
	}
}
