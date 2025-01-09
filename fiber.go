package myfx

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func NewFiberServer() *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler:          handleFiberError,
	})
	return app
}

func handleFiberError(c *fiber.Ctx, err error) error {
	code := http.StatusInternalServerError
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{"message": err.Error()})
}
