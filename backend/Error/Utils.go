package Error

import (
	"assaultrifle/Form"

	"github.com/gofiber/fiber/v2"
)

func StatusBadRequest(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"status": "error",
		"error":  "Geçersiz giriş",
	})
}

func InvalidPerm(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "error",
		"message": "Kullanıcı admin yetkisine sahip değil",
	})
}

func NotFound(c *fiber.Ctx, text string) error {
	return c.JSON(fiber.Map{
		"status":  "error",
		"message": text,
	})
}



func CustomError(c *fiber.Ctx, text string) error {
	return c.JSON(fiber.Map{
		"status":  "error",
		"message": text,
	})
}

func CustomSuccess(c *fiber.Ctx, text string) error {
	return c.JSON(fiber.Map{
		"status":  "OK",
		"message": text,
	})
}


func CustomSuccessAuth(c *fiber.Ctx, data, text string) error {
	return c.JSON(fiber.Map{
		"status":  "OK",
		"message": text,
		"data": data,
	})
}

func CustomSuccessUser(c *fiber.Ctx, data[]Form.User, text string) error {
	return c.JSON(fiber.Map{
		"status":  "OK",
		"message": text,
		"data": data,
	})
}

func CustomSuccessProduct(c *fiber.Ctx, data[]Form.Product, text string) error {
	return c.JSON(fiber.Map{
		"status":  "OK",
		"message": text,
		"data": data,
	})
}


func CustomSuccessSeller(c *fiber.Ctx, data Form.Seller, text string) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "OK",
		"message": text,
		"data":    data,
	})
}



func CustomSuccessStatus(c *fiber.Ctx,  data interface{}, text string) error {
	return c.JSON(fiber.Map{
		"status":  "OK",
		"message": text,
		"data": data,
	})
}
