	package Handler

	import (
		"assaultrifle/Database"
		"assaultrifle/Error"
		"assaultrifle/Form"
		"fmt"

		"github.com/gofiber/fiber/v2"
	)

	func CreateProduct(c *fiber.Ctx) error {
		var reqbody Form.ProductBody

		if err := c.BodyParser(&reqbody); err != nil {
			return Error.StatusBadRequest(c)
		}

		if reqbody.Content == "" || reqbody.Name == "" || reqbody.Features == "" || reqbody.User == "" {
			return Error.CustomError(c, "Boş veri bulunmaması gerekiyor")
		}

		err, token := Database.CreateProduct(reqbody.Name, reqbody.Content, reqbody.Features, reqbody.User)
		if err != nil {
			return Error.CustomError(c, err.Error())
		}

		return Error.CustomSuccessAuth(c, token, "Başarıyla eklendi")
	}

	func CreateSeller(c *fiber.Ctx) error {
		var reqbody Form.SellerBody

		if err := c.BodyParser(&reqbody); err != nil {
			return Error.StatusBadRequest(c)
		}

		if reqbody.Name == "" || reqbody.User == "" || reqbody.	Website == "" || reqbody.Logo == "" { 
			return Error.CustomError(c, "Boş veri bulunmaması gerekiyor")
		}

		err, token := Database.CreateSeller(reqbody.Name, reqbody.Website, reqbody.Logo,  reqbody.User ,reqbody.ProductID)
		if err != nil {
			return Error.CustomError(c, err.Error())
		}

		return Error.CustomSuccessAuth(c, token, "Başarıyla eklendi")
	}

	func ProductByToken(c *fiber.Ctx) error {
		token := c.Params("token")
		if token == "" {
			return Error.CustomError(c, "Ürün token'ı sağlanmalıdır.")
		}

		productWithPrices, err := Database.GetProductWithAllPrices(token)
		if err != nil {
			if err.Error() == fmt.Sprintf("Ürün bulunamadı: %s", token) {
				return Error.NotFound(c, "Ürün bulunamadı.")
			}

			return Error.CustomError(c, err.Error())
		}

		return c.JSON(fiber.Map{"status": "success", "data": productWithPrices})
	}
	func AllProducts(c *fiber.Ctx) error {
		products, err := Database.GetAllProducts()
		if err != nil {
			return Error.CustomError(c, err.Error())
		}

		if len(products) == 0 {
			return Error.NotFound(c, "Hiç ürün bulunamadı.")
		}

		return Error.CustomSuccessProduct(c,products, "İşlem başarılı oldu")
	}


