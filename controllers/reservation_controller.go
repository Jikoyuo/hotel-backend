package controllers

import (
	"hotel-backend/config"
	"hotel-backend/models"

	"github.com/gofiber/fiber/v2"
)

func GetReservations(c *fiber.Ctx) error {
	var reservations []models.Reservation
	config.DB.Preload("Room").Find(&reservations)
	return c.JSON(reservations)
}

func GetReservationById(c *fiber.Ctx) error {
	id := c.Params("id")
	var reservation models.Reservation
	if err := config.DB.First(&reservation, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Reservasi tidak ditemukan!"})
	}
	return c.JSON(reservation)
}

func CreateReservation(c *fiber.Ctx) error {
	res := new(models.Reservation)

	if err := c.BodyParser(res); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var room models.Room
	if err := config.DB.First(&room, "id = ?", res.RoomID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Kamar tidak ditemukan!"})
	}
	var count int64

	config.DB.Model(&models.Reservation{}).Where(
		"room_id = ? AND check_in < ? AND check_out > ?",
		res.RoomID, res.CheckOut, res.CheckIn,
	).Count(&count)

	if count > 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Kamar sudah terisi di tanggal tersebut (Double Booking)",
		})
	}

	if result := config.DB.Create(&res); result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Booking Success!",
		"data":    res,
	})
}
