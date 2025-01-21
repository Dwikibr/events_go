package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

func ParseStrIdToInt64(context *gin.Context, param string) int64 {
	parsedInt, err := strconv.ParseInt(context.Param(param), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}
	return parsedInt
}

func ParseInt64ToStr(id int64) string {
	return strconv.FormatInt(id, 10)
}

func Hasher(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func ValidatePassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	fmt.Println(err)
	return err == nil
}

func GetSignedKey() string {
	return "MySuperSecretKey"
}
