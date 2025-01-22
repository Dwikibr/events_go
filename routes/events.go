package routes

import (
	"RestApi/models"
	"RestApi/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// getEvents handles the GET request to fetch all events.
// It returns a JSON response with the list of events or an error message.
func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to Fetch Events Data", "err": err})
		return
	}
	context.JSON(http.StatusOK, events)
}

// createEvent handles the POST request to create a new event.
// It binds the JSON payload to the Event model, sets the UserID from the token info, and saves the event.
// It returns a JSON response with the created event or an error message.
func createEvent(context *gin.Context) {
	var event models.Event
	tokenInfo := context.MustGet("tokenInfo").(map[string]interface{})

	err := context.ShouldBindJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event.UserID = int(tokenInfo["id"].(float64))
	err = event.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to Save Event"})
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"message": "Event Created",
		"event":   event,
	})
}

// getEvent handles the GET request to fetch a specific event by its ID.
// It returns a JSON response with the event details or an error message.
func getEvent(context *gin.Context) {
	eventId := utils.ParseStrIdToInt64(context, "id")

	row, err := models.GetEventById(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Event with ID %d not found", eventId)})
		return
	}

	context.JSON(http.StatusOK, row)
}

// updateEvent handles the PUT request to update an existing event by its ID.
// It checks if the current user is authorized to update the event, binds the JSON payload, and updates the event.
// It returns a JSON response with the updated event or an error message.
func updateEvent(context *gin.Context) {
	eventID := utils.ParseStrIdToInt64(context, "id")

	event, err := models.GetEventById(eventID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Event with ID %d not found", eventID)})
		return
	}

	currSession := context.GetInt("userId")
	if currSession != event.UserID {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	err = context.ShouldBindJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = event.Update()
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to Update Event"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event Updated", "event": event})
}

// deleteEvent handles the DELETE request to delete an existing event by its ID.
// It checks if the current user is authorized to delete the event and deletes the event.
// It returns a JSON response with a success message or an error message.
func deleteEvent(context *gin.Context) {
	eventID := utils.ParseStrIdToInt64(context, "id")
	event, err := models.GetEventById(eventID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Event with ID %d not found", eventID)})
		return
	}

	if context.GetInt("userId") != event.UserID {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	err = event.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to Delete Event"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Event Deleted"})
}

// GetEventDetail handles the GET request to fetch the details of a specific event by its ID.
// It checks if the current user is authorized to view the event details and fetches the event with its registrations.
// It returns a JSON response with the event details or an error message.
func GetEventDetail(context *gin.Context) {
	eventId := utils.ParseStrIdToInt64(context, "id")
	event, err := models.GetEventById(eventId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("event with id %v not found", eventId)})
		return
	}

	curSession := context.GetInt("userId")
	fmt.Printf("Session: %v OwnerId: %v", curSession, event.UserID)
	if event.UserID != curSession {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized to get event detail"})
		return
	}

	err = event.GetEventWithRegistration()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Event with ID %d not found", eventId), "error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, event)
}
