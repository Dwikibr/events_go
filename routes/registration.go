package routes

import (
	"RestApi/models"
	"RestApi/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterEvent(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId := utils.ParseStrIdToInt64(context, "id")
	registration := models.CreateRegistration(userId, eventId)
	err := registration.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = registration.Save()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Registration Successful", "registration": registration})
}

func CancelRegistration(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId := context.Param("id")

	registration := models.CreateRegistration(userId, utils.ParseStrIdToInt64(context, eventId))
	err := registration.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = registration.Cancel()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Registration Cancelled"})
}
