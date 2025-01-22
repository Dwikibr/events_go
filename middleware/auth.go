package middleware

import (
	"RestApi/utils"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthenticateUser(context *gin.Context) {
	token := context.GetHeader("Authorization")
	parsedToken, err := VerifyToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
	}

	tokenInfo, err := TokenExtractor(parsedToken)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
	}
	context.Set("tokenInfo", tokenInfo)
	context.Set("currentUser", tokenInfo["username"])

	idFloat := tokenInfo["id"].(float64)
	userId := int(idFloat)
	context.Set("userId", userId)
	context.Next()
}

func VerifyToken(token string) (*jwt.Token, error) {
	//Verify token signing method
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(utils.GetSignedKey()), nil
	})
	if err != nil {
		return nil, errors.New("failed to parse token")
	}

	// Verify token is valid
	if !parsedToken.Valid {
		return nil, errors.New("invalid Token")
	}

	return parsedToken, nil
}

func TokenExtractor(parsedToken *jwt.Token) (map[string]interface{}, error) {
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to parse claims")
	}
	email := claims["username"]
	id := claims["id"]
	return map[string]interface{}{"username": email, "id": id}, nil
}
