package main

import (
	"RestApi/db"
	"RestApi/routes"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Start...")

	seed := flag.Bool("seed", false, "Run database seeder")
	flag.Parse()
	fmt.Println("-seed:", *seed)

	db.InitDB()

	if *seed {
		fmt.Println("Running database seeder...")
		db.SeedDB() // Run the seeder
		fmt.Println("Seeding completed!")
	}

	server := gin.Default()
	routes.SetRoutes(server)

	err := server.Run(":8181")
	if err != nil {
		fmt.Println(err)
	}
}
