package main

import (
	"RestApi/db"
	"RestApi/routes"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Start")
	db.InitDB()

	server := gin.Default()
	routes.SetRoutes(server)

	err := server.Run(":8181")
	if err != nil {
		fmt.Println(err)
	}
}
