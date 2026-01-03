package controllers

import (
	"hotel-backend/config"
	"hotel-backend/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Request body khusus Login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx) error {
	// 1. Ambil input user
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// 2. Cek apakah User ada di DB?
	var user models.User
	// Cari berdasarkan email
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Email atau Password salah"}) // Jangan bilang "Email tidak ditemukan" (Security reason)
	}

	// 3. Cek Password (Bandingkan Hash)
	// Ingat: user.Password itu Hash, req.Password itu Polos
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Email atau Password salah"})
	}

	// 4. Generate JWT Token
	// Membuat Claims (Isi data di dalam token)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	// Buat token dengan metode signing HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tanda tangani token dengan Secret Key (Harus rahasia!)
	// Nanti ini wajib dipindah ke .env
	secretKey := "rahasia-dapur-hotel-backend"
	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal generate token"})
	}

	// 5. Kirim Token ke User
	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   t,
	})
}
