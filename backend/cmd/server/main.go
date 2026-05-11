package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pablozagrodnik/reservation-microservice-showcase/internal/models"
	"github.com/pablozagrodnik/reservation-microservice-showcase/internal/repository/pg"
)

// odpowiedzi json
type SeatResponse struct {
	ID      uint `json:"id"`
	Row     int  `json:"row"`
	Col     int  `json:"col"`
	IsTaken bool `json:"is_taken"`
}

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

	r.GET("/api/screenings/:id/seats", func(c *gin.Context) {
		screeningID := c.Param("id")

		// pobranie danych o seansie
		var screening models.Screening
		if err := db.Preload("Room").First(&screening, screeningID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Seans nie istnieje"})
			return
		}

		// pobranie miejsc sali
		var seats []models.Seat
		db.Where("room_id = ?", screening.RoomID).Order("row asc, col asc").Find(&seats)

		// pobranie istniejących rezerwacji dla tego seansu
		var reservations []models.Reservation
		db.Where("screening_id = ?", screeningID).Find(&reservations)

		// mapa zajętych miejsc
		takenMap := make(map[uint]bool)
		for _, res := range reservations {
			takenMap[res.SeatID] = true
		}

		// mapowanie miejsc na odpowiedzi json
		var response []SeatResponse
		for _, seat := range seats {
			response = append(response, SeatResponse{
				ID:      seat.ID,
				Row:     seat.Row,
				Col:     seat.Col,
				IsTaken: takenMap[seat.ID],
			})
		}

		c.JSON(http.StatusOK, response)
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
