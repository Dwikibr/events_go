package routes

import (
	"RestApi/models"
	"RestApi/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SignUp(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = user.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to SignUp"})
		return
	}
	returnedUserdata := map[string]string{"id": utils.ParseInt64ToStr(user.ID), "username": user.Username}
	context.JSON(http.StatusCreated, gin.H{"message": "User Created", "user": returnedUserdata})
}

func Login(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Username and Password are required"})
		return
	}
	err = user.ValidateCredential()
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Credential"})
		return
	}
	jwtToken, tokenErr := utils.GenerateToken(user.Username, user.ID)
	if tokenErr != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to Generate Token"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login Success", "token": jwtToken})
}
