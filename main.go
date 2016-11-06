package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
)

var db *sql.DB // Connection pool

func main() {

	// Setup DB connection
	initDatabaseConnection()

	// Initialize the router
	router := initRouter()

	// Serve it like it's hot
	router.Run(standard.New(":8080"))
}

func initDatabaseConnection() {
	// Get DB connection details
	host := os.Getenv("MYSQL_HOST")
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	database := os.Getenv("MYSQL_DATABASE")

	// Initialize DB connection
	sqlConnString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", user, pass, host, database)
	conn, err := sql.Open("mysql", sqlConnString)
	if err != nil {
		log.Fatal("Could not initialize DB connection:", err)
	}
	db = conn // Assign the pool to the global one :)
}

func initRouter() *echo.Echo {
	e := echo.New()

	e.POST("/shorten", createShortURL)
	e.GET("/:shortcode", redirectTo)
	e.GET("/:shortcode/stats", getShortcodeStats)

	return e
}
