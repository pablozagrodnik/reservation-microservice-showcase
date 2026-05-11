package models

import "time"

type Movie struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Poster      string      `json:"poster"`
	Screenings  []Screening `json:"screenings"`
}

type Room struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	Name       string      `json:"name"`
	Seats      []Seat      `json:"-"`
	Screenings []Screening `json:"-"`
}

type Screening struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MovieID   uint      `json:"-"`
	StartTime time.Time `json:"start_time"`
	RoomID    uint      `json:"room_id"`
	Room      Room      `gorm:"foreignKey:RoomID" json:"room"`
}

type Seat struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	RoomID uint `json:"-"`
	Row    int  `json:"row"`
	Col    int  `json:"col"`
}

type Reservation struct {
	ID          uint `gorm:"primaryKey" json:"id"`
	ScreeningID uint `json:"screening_id"`
	SeatID      uint `json:"seat_id"`
}
