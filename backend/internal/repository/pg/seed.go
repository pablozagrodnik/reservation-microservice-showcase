package pg

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/pablozagrodnik/reservation-microservice-showcase/internal/models"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {

	db.Exec("TRUNCATE TABLE reservations, seats, screenings, movies, rooms RESTART IDENTITY CASCADE")

	// 1. sale i miejsca - jeśli baza jest pusta
	var roomCount int64
	db.Model(&models.Room{}).Count(&roomCount)

	var rooms []models.Room
	for i := 1; i <= 5; i++ {
		rooms = append(rooms, models.Room{Name: fmt.Sprintf("Sala %d", i)})
	}
	db.Create(&rooms)

	for i, room := range rooms {
		rows, cols := 5, 10

		if i == 0 {
			rows, cols = 4, 5
		} else if i >= 3 {
			rows, cols = 8, 15
		}

		var seatsInRoom []models.Seat
		for r := 1; r <= rows; r++ {
			for c := 1; c <= cols; c++ {
				seat := models.Seat{RoomID: room.ID, Row: r, Col: c}
				seatsInRoom = append(seatsInRoom, seat)
			}
		}
		db.Create(&seatsInRoom)
	}

	// 2. filmy - jeśli baza jest pusta
	var movieCount int64
	db.Model(&models.Movie{}).Count(&movieCount)

	movies := []models.Movie{
		{
			Title:       "Diuna: Część Druga",
			Description: "Paul Atryda jednoczy się z Chani i Fremenami, by zemścić się na spiskowcach, którzy zniszczyli jego rodzinę.",
			Poster:      "https://fwcdn.pl/fpo/34/81/10003481/8115126_1.10.webp",
		},
		{
			Title:       "Deadpool & Wolverine",
			Description: "Wolverine leczy rany, gdy jego drogi krzyżują się z pyskatym Deadpoolem. Wspólnie łączą siły, by pokonać wspólnego wroga.",
			Poster:      "https://fwcdn.pl/fpo/85/13/868513/8133073.10.webp",
		},
		{
			Title:       "Furioza",
			Description: "Policjantka składa swojemu byłemu chłopakowi propozycję: albo przeniknie do grupy kibiców, albo jego brat trafi do więzienia.",
			Poster:      "https://fwcdn.pl/fpo/38/10/833810/7972235_1.10.webp",
		},
		{
			Title:       "Friz & Wersow. Miłość w czasach online",
			Description: "Dokumentalna podróż przez życie najpopularniejszej pary influencerów w Polsce, pokazująca kulisy ich związku i kariery.",
			Poster:      "https://fwcdn.pl/fpo/21/22/10092122/8199797.10.webp",
		},
		{
			Title:       "Interstellar",
			Description: "Grupa astronautów podróżuje przez tunel czasoprzestrzenny w poszukiwaniu nowego domu dla ludzkości.",
			Poster:      "https://fwcdn.pl/fpo/56/29/375629/7670122_2.10.webp",
		},
		{
			Title:       "8 Mila",
			Description: "Młody raper z Detroit próbuje swoich sił w bitwach freestyle'owych, walcząc o uznanie i lepszą przyszłość.",
			Poster:      "https://fwcdn.pl/fpo/23/58/32358/7537321_1.10.webp",
		},
		{
			Title:       "Królestwo Planety Małp",
			Description: "Lata po panowaniu Cezara, młoda małpa wyrusza w podróż, która zmieni przyszłość obu gatunków.",
			Poster:      "https://fwcdn.pl/fpo/88/22/848822/8128230.10.webp",
		},
		{
			Title:       "Jesteś Bogiem",
			Description: "Oparta na faktach historia powstania i tragicznych losów legendarnej polskiej grupy hip-hopowej Paktofonika.",
			Poster:      "https://fwcdn.pl/fpo/45/45/544545/7478694_1.10.webp",
		},
		{
			Title:       "Miś",
			Description: "Prezes klubu sportowego Tęcza próbuje wyjechać do Londynu, mierząc się z absurdami codzienności w czasach PRL.",
			Poster:      "https://fwcdn.pl/fpo/08/96/896/8029196_1.10.webp",
		},
		{
			Title:       "Viy",
			Description: "XVIII wiek. Kartograf Jonathan Green wyrusza w podróż na Wschód, trafiając do tajemniczej i mrocznej wioski.",
			Poster:      "https://fwcdn.pl/fpo/25/98/102598/6918815_1.10.webp",
		},
	}
	db.Create(&movies)

	// odświeżanie seansów
	var futureScreenings int64
	db.Model(&models.Screening{}).Where("start_time >= ?", time.Now()).Count(&futureScreenings)

	// generowanie nowej puli
	if futureScreenings == 0 {
		var rooms []models.Room
		db.Find(&rooms)

		var movies []models.Movie
		db.Find(&movies)

		hours := []int{10, 13, 16, 19, 22}
		now := time.Now()
		movieIndex := 0

		for d := 0; d < 10; d++ {
			for h, hour := range hours {
				startTime := time.Date(now.Year(), now.Month(), now.Day()+d, hour, 0, 0, 0, now.Location())

				// pomijanie odbytych
				if startTime.Before(now) {
					continue
				}

				movieIdx := (d + h) % len(movies)
				roomIdx := (d + h) % len(rooms)

				movie := movies[movieIdx]
				room := rooms[roomIdx]

				screening := models.Screening{
					MovieID:   movie.ID,
					StartTime: startTime,
					RoomID:    room.ID,
				}
				db.Create(&screening)

				// losowanie rezerwacji
				occupancyRate := 0.8 - (float64(d) * 0.07)
				if occupancyRate < 0.1 {
					occupancyRate = 0.1
				}

				var seatsInRoom []models.Seat
				db.Where("room_id = ?", room.ID).Find(&seatsInRoom)

				numToOccupy := int(float64(len(seatsInRoom)) * occupancyRate)

				rand.Shuffle(len(seatsInRoom), func(i, j int) {
					seatsInRoom[i], seatsInRoom[j] = seatsInRoom[j], seatsInRoom[i]
				})

				for i := 0; i < numToOccupy; i++ {
					db.Create(&models.Reservation{
						ScreeningID: screening.ID,
						SeatID:      seatsInRoom[i].ID,
					})
				}

				movieIndex++
			}
		}
	}

	return nil
}
