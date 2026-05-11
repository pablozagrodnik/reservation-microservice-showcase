package pg

import (
	"fmt"
	"github.com/pablozagrodnik/reservation-microservice-showcase/internal/models"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	var count int64
	db.Model(&models.Movie{}).Count(&count)
	if count > 0 {
		return nil
	}

	// 1. sale
	var rooms []models.Room
	for i := 1; i <= 5; i++ {
		rooms = append(rooms, models.Room{Name: fmt.Sprintf("Sala %d", i)})
	}
	db.Create(&rooms)

	// 2. miejsca w sali - różne ilości
	var seats []models.Seat
	for i, room := range rooms {
		rows, cols := 5, 10

		if i == 0 {
			rows, cols = 4, 5
		} else if i >= 3 {
			rows, cols = 8, 15
		}

		for r := 1; r <= rows; r++ {
			for c := 1; c <= cols; c++ {
				seats = append(seats, models.Seat{RoomID: room.ID, Row: r, Col: c})
			}
		}
	}
	db.Create(&seats)

	// 3. filmy
	movies := []models.Movie{
		{Title: "Diuna: Część Druga", Poster: "https://fwcdn.pl/fpo/34/81/10003481/8115126_1.10.webp"},
		{Title: "Deadpool & Wolverine", Poster: "https://fwcdn.pl/fpo/85/13/868513/8133073.10.webp"},
		{Title: "Furioza", Poster: "https://fwcdn.pl/fpo/38/10/833810/7972235_1.10.webp"},
		{Title: "Friz & Wersow. Miłość w czasach online", Poster: "https://fwcdn.pl/fpo/21/22/10092122/8199797.10.webps"},
		{Title: "Interstellar", Poster: "https://fwcdn.pl/fpo/56/29/375629/7670122_2.10.webp"},
		{Title: "8 Mila", Poster: "https://fwcdn.pl/fpo/23/58/32358/7537321_1.10.webp"},
		{Title: "Królestwo Planety Małp", Poster: "https://fwcdn.pl/fpo/88/22/848822/8128230.10.webp"},
		{Title: "Jesteś Bogiem", Poster: "https://fwcdn.pl/fpo/45/45/544545/7478694_1.10.webp"},
		{Title: "Miś", Poster: "https://fwcdn.pl/fpo/08/96/896/8029196_1.10.webp"},
		{Title: "Viy", Poster: "https://fwcdn.pl/fpo/25/98/102598/6918815_1.10.webp"},
	}
	db.Create(&movies)

	// 4. seanse
	timeslots := []string{"10:00", "13:00", "16:00", "19:00", "22:00"}
	var screenings []models.Screening

	movieIndex := 0
	for _, room := range rooms {
		for _, time := range timeslots {
			screenings = append(screenings, models.Screening{
				MovieID: movies[movieIndex%len(movies)].ID,
				Time:    time,
				RoomID:  room.ID,
			})
			movieIndex++
		}
	}
	db.Create(&screenings)

	return nil
}
