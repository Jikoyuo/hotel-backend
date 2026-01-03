package controllers

import (
	"hotel-backend/config"
	"hotel-backend/models"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	TotalPage   int   `json:"total_page"`
	PageSize    int   `json:"page_size"`
	TotalItems  int64 `json:"total_items"`
}
type PaginatedResponse struct {
	Meta PaginationMeta `json:"meta"`
	Data interface{}    `json:"data"`
}

func CreateRoom(c *fiber.Ctx) error {
	room := new(models.Room)

	if err := c.BodyParser(room); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	result := config.DB.Create(&room)

	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": result.Error.Error(),
		})
	}
	return c.JSON(room)
}

func GetReservations(c *fiber.Ctx) error {
	var reservations []models.Reservation
	config.DB.Preload("Room").Find(&reservations)
	return c.JSON(reservations)
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

func GetRoomById(c *fiber.Ctx) error {
	id := c.Params("id")

	var room models.Room
	if err := config.DB.First(&room, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Kamar tidak ditemukan!"})
	}
	return c.JSON(room)
}

func GetRooms(c *fiber.Ctx) error {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	var rooms []models.Room
	var total int64

	config.DB.Model(&models.Room{}).Count(&total)

	if pageStr == "" && limitStr == "" {
		result := config.DB.Find(&rooms)

		if result.Error != nil {
			return c.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
		}

		response := PaginatedResponse{
			Meta: PaginationMeta{
				CurrentPage: 1,
				TotalPage:   1,
				PageSize:    int(total),
				TotalItems:  total,
			},
			Data: rooms,
		}
		return c.JSON(response)
	}

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	result := config.DB.Offset(offset).Limit(limit).Find(&rooms)

	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
	}

	totalPage := math.Ceil(float64(total) / float64(limit))

	response := PaginatedResponse{
		Meta: PaginationMeta{
			CurrentPage: page,
			TotalPage:   int(totalPage),
			PageSize:    limit,
			TotalItems:  total,
		},
		Data: rooms,
	}

	return c.JSON(response)
}
