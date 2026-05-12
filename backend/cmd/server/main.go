package main

import (
	"log"
	"net/http"
	"os"
	"time"

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

type CreateRoomRequest struct {
	Name string `json:"name" binding:"required"`
	Rows int    `json:"rows" binding:"required"`
	Cols int    `json:"cols" binding:"required"`
}

type CreateScreeningRequest struct {
	MovieID   uint      `json:"movie_id" binding:"required"`
	RoomID    uint      `json:"room_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
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

	r.Use(CORSMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/movies", func(c *gin.Context) {
		var movies []models.Movie

		// pobranie filmów z dołączeniem seansów
		result := db.Preload("Screenings.Room").Find(&movies)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd pobierania danych z bazy"})
			return
		}

		c.JSON(http.StatusOK, movies)
	})

	r.GET("/screenings/:id/seats", func(c *gin.Context) {
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

	r.POST("/reservations", func(c *gin.Context) {
		var res models.Reservation
		if err := c.ShouldBindJSON(&res); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Błędne dane rezerwacji"})
			return
		}

		// sprawdzenie zajętości miejsca
		var existing models.Reservation
		err := db.Where("screening_id = ? AND seat_id = ?", res.ScreeningID, res.SeatID).First(&existing).Error
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "To miejsce zostało właśnie zajęte przez kogoś innego!"})
			return
		}

		// zapis rezerwacji
		if err := db.Create(&res).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd podczas zapisu w bazie"})
			return
		}

		log.Printf("Rezerwacja nr %d zapisana dla: %s", res.ID, res.UserEmail)

		c.JSON(http.StatusCreated, res)
	})

	// panel admina (pełny crud)
	admin := r.Group("/admin")
	{
		// --- SALE ---
		admin.GET("/rooms", func(c *gin.Context) {
			var rooms []models.Room
			db.Find(&rooms)
			c.JSON(http.StatusOK, rooms)
		})

		admin.POST("/rooms", func(c *gin.Context) {
			var req CreateRoomRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			room := models.Room{Name: req.Name}
			if err := db.Create(&room).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd tworzenia sali"})
				return
			}

			var seats []models.Seat
			for r := 1; r <= req.Rows; r++ {
				for col := 1; col <= req.Cols; col++ {
					seats = append(seats, models.Seat{
						RoomID: room.ID,
						Row:    r,
						Col:    col,
					})
				}
			}
			db.Create(&seats)

			c.JSON(http.StatusCreated, gin.H{"message": "Sala i miejsca wygenerowane", "room": room})
		})

		admin.DELETE("/rooms/:id", func(c *gin.Context) {
			id := c.Param("id")
			if err := db.Unscoped().Delete(&models.Room{}, id).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd podczas usuwania sali"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Sala i wszystkie jej zależności zostały usunięte"})
		})

		// --- SEANSE ---
		admin.GET("/screenings", func(c *gin.Context) {
			var screenings []models.Screening
			db.Preload("Movie").Preload("Room").Find(&screenings)
			c.JSON(http.StatusOK, screenings)
		})

		admin.POST("/screenings", func(c *gin.Context) {
			var req CreateScreeningRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			screening := models.Screening{
				MovieID:   req.MovieID,
				RoomID:    req.RoomID,
				StartTime: req.StartTime,
			}
			db.Create(&screening)
			c.JSON(http.StatusCreated, screening)
		})

		admin.PATCH("/screenings/:id", func(c *gin.Context) {
			id := c.Param("id")
			var screening models.Screening

			if err := db.First(&screening, id).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Nie znaleziono seansu"})
				return
			}

			if err := c.ShouldBindJSON(&screening); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			db.Save(&screening)
			c.JSON(http.StatusOK, screening)
		})

		admin.DELETE("/screenings/:id", func(c *gin.Context) {
			id := c.Param("id")
			if err := db.Unscoped().Delete(&models.Screening{}, id).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd podczas usuwania seansu"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Seans usunięto pomyślnie"})
		})

		// --- FILMY ---
		admin.POST("/movies", func(c *gin.Context) {
			var movie models.Movie
			if err := c.ShouldBindJSON(&movie); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			db.Create(&movie)
			c.JSON(http.StatusCreated, movie)
		})

		admin.PATCH("/movies/:id", func(c *gin.Context) {
			id := c.Param("id")
			var movie models.Movie

			if err := db.First(&movie, id).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Nie znaleziono filmu"})
				return
			}

			if err := c.ShouldBindJSON(&movie); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			db.Save(&movie)
			c.JSON(http.StatusOK, movie)
		})

		admin.DELETE("/movies/:id", func(c *gin.Context) {
			id := c.Param("id")
			if err := db.Unscoped().Delete(&models.Movie{}, id).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd podczas usuwania filmu"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Film usunięto pomyślnie"})
		})

		// --- REZERWACJE ---
		admin.GET("/reservations", func(c *gin.Context) {
			var reservations []models.Reservation
			db.Preload("Screening.Movie").Preload("Seat").Find(&reservations)
			c.JSON(http.StatusOK, reservations)
		})

		admin.DELETE("/reservations/:id", func(c *gin.Context) {
			id := c.Param("id")
			if err := db.Unscoped().Delete(&models.Reservation{}, id).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd podczas usuwania rezerwacji"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Rezerwacja anulowana pomyślnie"})
		})
	}

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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
