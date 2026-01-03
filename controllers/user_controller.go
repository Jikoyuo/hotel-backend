package controllers

import (
	"hotel-backend/config"
	"hotel-backend/models"
	"hotel-backend/utils"
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// Struct khusus untuk menangkap Request Body
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUser(c *fiber.Ctx) error {
	// 1. Pakai struct request khusus
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// 2. Hash Password
	// Cost 10-14 adalah standar yang baik.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal memproses password"})
	}

	// 3. Masukkan ke Struct Model User
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword), // Simpan versi acak, bukan mentah
		Role:     models.RoleGuest,
	}

	// 4. Simpan ke DB
	result := config.DB.Create(&user)

	if result.Error != nil {
		// Cek apakah errornya karena duplikat
		if strings.Contains(result.Error.Error(), "duplicate key") {
			return c.Status(409).JSON(fiber.Map{"error": "Username atau Email sudah terdaftar!"})
		}
		return c.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
	}

	// 5. Return JSON (Password field otomatis hilang karena json:"-")
	return c.Status(201).JSON(fiber.Map{
		"message": "User registered successfully",
		"data":    user,
	})
}

func GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	if err := config.DB.First(&user, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan!"})
	}
	return c.JSON(user)
}

func GetUsers(c *fiber.Ctx) error {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	var users []models.User
	var total int64

	config.DB.Model(&models.User{}).Count(&total)

	if pageStr == "" && limitStr == "" {
		result := config.DB.Find(&users)

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
			Data: users,
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
	result := config.DB.Offset(offset).Limit(limit).Find(&users)

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
		Data: users,
	}

	return c.JSON(response)
}
