package main

import (
	"os"

	"github.com/Rafii271/UTS-EAI/config"
	"github.com/Rafii271/UTS-EAI/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {

	// Load ENV Files
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	// Connect to database
	config.Connect(&config.Config{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
		SSLMODE:  os.Getenv("DB_SSLMODE"),
	})

	// Start Fiber Server
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Token, Type",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// Routes
	routes.Setup(app)

	err = app.Listen(":" + os.Getenv("PORT"))

	if err != nil {
		panic(err)
	}
}
