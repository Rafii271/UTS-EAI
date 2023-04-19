package services

import (
	"github.com/gofiber/fiber/v2"
)

type UserService interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	ShowProfile(c *fiber.Ctx) error
}

