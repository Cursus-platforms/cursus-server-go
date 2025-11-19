package main

import (
	"log"
	"net/http"

	_ "github.com/Cursus-platforms/cursus-server-go/docs"
	"github.com/Cursus-platforms/cursus-server-go/internal/infrastructure/mailer"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/Cursus-platforms/cursus-server-go/internal/features/auth"
	"github.com/Cursus-platforms/cursus-server-go/internal/features/role"
	"github.com/Cursus-platforms/cursus-server-go/internal/features/user"
	"github.com/Cursus-platforms/cursus-server-go/internal/infrastructure/db"
	"github.com/Cursus-platforms/cursus-server-go/internal/infrastructure/redis"
	"github.com/Cursus-platforms/cursus-server-go/internal/infrastructure/router"
)

// @title           E-Learning Platform API
// @version         1.0
// @description     This is the API documentation for the E-Learning project.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Your Name
// @contact.url    http://www.your-website.com
// @contact.email  your.email@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

var dbConn *sqlx.DB

func main() {
	//Connect Postgres db
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatalf("FATAL: Could not connect to database: %v", err)
	}
	log.Println("Successfully connected to the database!")
	defer dbConn.Close()

	//Connect Redis db
	rdb, err := redis.Connect()
	if err != nil {
		log.Fatal("FATAL: Could not connected to Redis!")
	}
	defer rdb.Close()

	//Init mailer service
	mailSvc := mailer.NewSMTPMailer()

	//Init repositories
	userRepo := user.NewRepository(dbConn)
	roleRepo := role.NewRepository(dbConn)

	//Init services
	authSvc := auth.NewService(userRepo, roleRepo, rdb, mailSvc)

	//Init handlers
	authHandler := auth.NewHandler(authSvc)

	allHandler := router.HandleDependencies{
		AuthHandler: authHandler,
	}

	mainRouter := router.NewRouter(allHandler)

	port := ":8080"
	log.Printf("Starting REST API server on http://localhost%s\n", port)
	log.Printf("Swagger docs available at http://localhost%s/swagger/index.html", port)

	err = http.ListenAndServe(port, mainRouter)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
