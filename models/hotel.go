package models

import (
	"time"

	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	Number string  `json: "number"` // nomor kamar hotel
	Type   string  `json:"type"`    //tipe kamar, misal Deluxe
	Price  float64 `json:"price"`   // harga kamar
}

type Reservation struct {
	gorm.Model
	GuestName string    `json:"guest_name"`
	CheckIn   time.Time `json:"check_in"`
	CheckOut  time.Time `json:"check_out"`

	// Foreign Key (Menghubungkan ke Room)
	RoomID uint `json:"room_id"`
	Room   Room `json:"room" gorm:"foreignKey:RoomID"` // Relasi untuk Preload
}
