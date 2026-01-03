package models

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomType string

const (
	TypeStandard     RoomType = "STANDARD"
	TypeDeluxe       RoomType = "DELUXE"
	TypePresidential RoomType = "PRESIDENTIAL"
)

type Base struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type Room struct {
	Base
	Number string   `json:"number"`
	Type   RoomType `json:"type"`
	Price  float64  `json:"price"`
}

func (r *Room) BeforeCreate(tx *gorm.DB) (err error) {
	if len(r.ID) == 0 {
		r.ID = uuid.New().String()
	}
	r.Type = RoomType(strings.ToUpper(string(r.Type)))
	switch r.Type {
	case TypeStandard, TypeDeluxe, TypePresidential:
		return nil
	default:
		return errors.New("tipe kamar tidak valid (Harus STANDARD, DELUXE, atau PRESIDENTIAL)")
	}
}

type Reservation struct {
	Base
	GuestName string    `json:"guest_name"`
	CheckIn   time.Time `json:"check_in"`
	CheckOut  time.Time `json:"check_out"`
	RoomID    string    `json:"room_id"`
	Room      Room      `json:"room" gorm:"foreignKey:RoomID"`
}

func (res *Reservation) BeforeCreate(tx *gorm.DB) (err error) {
	if len(res.ID) == 0 {
		res.ID = uuid.New().String()
	}
	return
}
