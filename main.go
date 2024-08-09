package main

import (
	"log"

	"avito-backend-bootcamp/database"
	"avito-backend-bootcamp/routers"
)

func main() {
	err := database.InitDB()

	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	routes := routers.ApiHandleFunctions{}
	log.Printf("Server started")

	router := routers.NewRouter(routes)
	log.Fatal(router.Run(":8080"))
}
