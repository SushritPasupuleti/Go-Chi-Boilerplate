package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"server/db"
	"server/models"
	"server/routes"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

type Application struct {
	Config Config
	Models models.Models
}

func (app *Application) Serve() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return err
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
		return err
	}

	log.Printf("Listening on port %s", port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: routes.Routes(),
	}

	return srv.ListenAndServe()
}
// @title Swagger Example API
// @version 1.0
// @description This is a sample server.
// @termsOfService http://localhost:6000/terms/

// @contact.name API Support
// @contact.url http://www.localhost:6000/support
// @contact.email support@localhost:6000

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:5000
// @BasePath /
func main() {
	fmt.Println("Hello there!")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	dbConn, err := db.Connect(dsn)
	if err != nil {
		log.Fatal("Error connecting to database", err)
	}

	defer dbConn.DB.Close()

	app := Application{
		Config: Config{
			Port: os.Getenv("PORT"),
		},
		Models: models.New(dbConn.DB),
	}

	err = app.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
