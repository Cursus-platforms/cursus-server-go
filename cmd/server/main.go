// File: cmd/api/main.go
package main

import (
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/Cursus-platforms/cursus-server-go/internal/db"
)

func main() {
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatalf("FATAL: Could not connect to database: %v", err)
	}
	log.Println("✅ Successfully connected to the database!")
	defer dbConn.Close()

	mainRouter := router.NewRouter(productHandler)

	port := ":8080" 
	log.Printf("Starting REST API server on http://localhost%s\n", port)

	err = http.ListenAndServe(port, mainRouter)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}