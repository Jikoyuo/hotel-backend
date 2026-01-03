package main

import (
	"hotel-backend/config"
	"hotel-backend/controllers"
	"hotel-backend/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	config.ConnectDB()
	app := fiber.New()
	app.Use(cors.New())
	api := app.Group("/api")

	api.Post("/login", controllers.Login)      // Login dulu buat dapat token
	api.Post("/users", controllers.CreateUser) // Register user baru
	api.Get("/rooms", controllers.GetRooms)    // Get daftar kamar (bebas)

	api.Use(middleware.Protected)
	api.Post("/rooms", middleware.Authorize("MANAGER", "MASTER ADMIN"), controllers.CreateRoom)
	api.Get("/room/:id", middleware.Authorize("MANAGER", "MASTER ADMIN", "RECEPTIONIST"), controllers.GetRoomById)

	api.Post("/reservation", middleware.Authorize("GUEST", "RECEPTIONIST"), controllers.CreateReservation)
	api.Get("/reservation", middleware.Authorize("RECEPTIONIST"), controllers.GetReservations)

	api.Get("/users", middleware.Authorize("MASTER ADMIN"), controllers.GetUsers)
	api.Get("/users/:id", middleware.Authorize("MASTER ADMIN"), controllers.GetUserById)
	app.Listen(":8080")
}
