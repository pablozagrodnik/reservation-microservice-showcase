package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pablozagrodnik/reservation-microservice-showcase/internal/models"
	"github.com/pablozagrodnik/reservation-microservice-showcase/internal/repository/pg"
)

func main() {

	db, err := pg.NewDB()
	if err != nil {
		log.Fatalf("Nie udało się połączyć z bazą: %v", err)
	}

	if err := pg.Migrate(db); err != nil {
		log.Fatalf("Błąd migracji: %v", err)
	}

	if err := pg.Seed(db); err != nil {
		log.Fatalf("Błąd seedowania danych: %v", err)
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/api/movies", func(c *gin.Context) {
		var movies []models.Movie

		// pobranie filmów z dołączeniem seansów
		result := db.Preload("Screenings.Room").Find(&movies)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd pobierania danych z bazy"})
			return
		}

		c.JSON(http.StatusOK, movies)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Start na porcie: %s", port)
	err = r.Run(":" + port)
	if err != nil {
		return
	}
}
