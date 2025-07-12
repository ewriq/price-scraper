package Routes

import (
	Handler "assaultrifle/Handler"

	"github.com/gofiber/fiber/v2"
)

func Admin(app fiber.Router) {
	app.Get("/product/list", Handler.AllProducts)
	app.Get("/product/:token", Handler.ProductByToken)

	app.Post("/product", Handler.CreateProduct) 
	app.Post("/seller", Handler.CreateSeller)   

}