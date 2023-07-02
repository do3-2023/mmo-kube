package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func runMigration(dbURL, migrationPath string) error {
	m, err := migrate.New("file://"+migrationPath, dbURL)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func main() {
	godotenv.Load()

	port := 3000

	requiredEnvVars := []string{"PG_USER", "PG_PASSWORD", "PG_DATABASE", "PG_HOSTNAME"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("Missing environment variable: %s", envVar)
		}
	}

	pgUser := os.Getenv("PG_USER")
	pgPassword := os.Getenv("PG_PASSWORD")
	pgDatabase := os.Getenv("PG_DATABASE")
	pgHostname := os.Getenv("PG_HOSTNAME")

	if os.Getenv("ENV") == "migrate" || os.Getenv("ENV") == "dev" {
		log.Print("Running migrations")

		dbURL := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", pgUser, pgPassword, pgHostname, pgPassword)
		migrationPath := "migrations"

		err := runMigration(dbURL, migrationPath)
		if err != nil {
			log.Fatal(err)
		}

		if os.Getenv("ENV") == "migrate" {
			os.Exit(0)
		}
	}

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=5432 sslmode=disable",
		pgUser,
		pgPassword,
		pgDatabase,
		pgHostname,
	))

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		var count int
		var lastCallRaw sql.NullTime

		err := db.QueryRow("SELECT COUNT(*), MAX(call_time) FROM calls").Scan(&count, &lastCallRaw)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO calls (call_time) VALUES ($1)", time.Now())
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		response := gin.H{
			"count":     count,
			"last_call": nil,
			"timestamp": time.Now().Format(time.RFC3339),
		}

		if lastCallRaw.Valid {
			response["last_call"] = lastCallRaw.Time.Format(time.RFC3339)
		}

		c.JSON(http.StatusOK, response)
	})

	router.GET("/healthz", func(c *gin.Context) {
		var result int
		err := db.QueryRow("SELECT 1").Scan(&result)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusNoContent)
	})

	log.Printf("Server is running on http://localhost:%d", port)

	err = router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
}
