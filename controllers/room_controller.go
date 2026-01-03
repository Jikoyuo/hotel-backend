package controllers

import (
	"hotel-backend/config"
	"hotel-backend/models"
	"hotel-backend/utils"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

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

		response := utils.PaginatedResponse{
			Meta: utils.PaginationMeta{
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

	response := utils.PaginatedResponse{
		Meta: utils.PaginationMeta{
			CurrentPage: page,
			TotalPage:   int(totalPage),
			PageSize:    limit,
			TotalItems:  total,
		},
		Data: rooms,
	}

	return c.JSON(response)
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

func GetRoomById(c *fiber.Ctx) error {
	id := c.Params("id")

	var room models.Room
	if err := config.DB.First(&room, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Kamar tidak ditemukan!"})
	}
	return c.JSON(room)
}
