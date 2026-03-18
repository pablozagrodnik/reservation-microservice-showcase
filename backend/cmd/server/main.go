package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
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
