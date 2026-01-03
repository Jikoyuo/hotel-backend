package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Protected(c *fiber.Ctx) error {
	// 1. Ambil Token dari Header "Authorization"
	// Format standar: "Bearer eyJhbGci..."
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized: Token wajib ada!"})
	}

	// 2. Buang kata "Bearer " agar sisa token-nya saja
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	// 3. Parse & Validasi Token
	secretKey := "rahasia-dapur-hotel-backend" // Harus SAMA dengan controller tadi

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Pastikan metode signing-nya HMAC (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return []byte(secretKey), nil
	})

	// Jika token rusak atau tanda tangan salah
	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized: Token tidak valid!"})
	}

	// 4. (Opsional) Ambil data dari token (User ID)
	// claims := token.Claims.(jwt.MapClaims)
	// userID := claims["user_id"]
	// c.Locals("user_id", userID) // Simpan ID user di Context biar bisa dipakai controller lain

	// 5. Lanjut ke Controller tujuan
	return c.Next()
}

// Fungsi ini menerima variadic parameter (bisa banyak role)
func Authorize(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil User dari Locals (Hasil dari middleware Protected sebelumnya)
		// ATAU kita parse ulang tokennya disini untuk ambil claim 'role'

		userToken := c.Get("Authorization")
		if userToken == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
		}

		tokenString := strings.Replace(userToken, "Bearer ", "", 1)

		// Parse Token (Sama kayak Protected, tapi kita butuh datanya)
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("rahasia-dapur-hotel-backend"), nil
		})

		// 2. Ambil Role dari dalam Token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
		}

		userRole := claims["role"].(string) // Casting ke string

		// 3. Logika "Dewa" (Master Admin & Manager selalu boleh)
		if userRole == "MASTER ADMIN" || userRole == "MANAGER" {
			return c.Next()
		}

		// 4. Cek apakah role user ada di daftar allowedRoles
		isAllowed := false
		for _, role := range allowedRoles {
			if role == userRole {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			return c.Status(403).JSON(fiber.Map{
				"error": "Forbidden: Anda tidak punya akses (" + userRole + ")",
			})
		}

		return c.Next()
	}
}
