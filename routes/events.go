package routes

import (
	"RestApi/models"
	"RestApi/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to Fetch Events Data", "err": err})
		return
	}
	context.JSON(http.StatusOK, events)
}

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

func getEvent(context *gin.Context) {

	eventId := utils.ParseStrIdToInt64(context, "id")

	row, err := models.GetEventById(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Event with ID %d not found", eventId)})
		return
	}

	context.JSON(http.StatusOK, row)
}

func updateEvent(context *gin.Context) {
	eventID := utils.ParseStrIdToInt64(context, "id")

	event, err := models.GetEventById(eventID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Event with ID %d not found", eventID)})
		return
	}

	currSession := context.GetInt("id")
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

func deleteEvent(context *gin.Context) {
	eventID := utils.ParseStrIdToInt64(context, "id")
	event, err := models.GetEventById(eventID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Event with ID %d not found", eventID)})
		return
	}

	if context.GetInt("id") != event.UserID {
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

func GetEventDetail(context *gin.Context) {
	eventId := utils.ParseStrIdToInt64(context, "id")
	event, err := models.GetEventWithRegistration(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Event with ID %d not found", eventId), "error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, event)
}
