package main

import (
	"hotel-backend/config"
	"hotel-backend/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	config.ConnectDB()
	app := fiber.New()
	app.Use(cors.New())
	api := app.Group("/api")
	api.Post("/rooms", controllers.CreateRoom)
	api.Get("/rooms", controllers.GetRooms)
	api.Get("/room/:id", controllers.GetRoomById)

	api.Post("/reservation", controllers.CreateReservation)
	api.Get("/reservation", controllers.GetReservations)
	app.Listen(":8080")
}
