package routes

import (
	"github.com/abr-ooo/hodhod/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/abr-ooo/hodhod/log"
)

func HooRouter(app fiber.Router) {

	api := app.Group("/hoo")

	api.Post("", controllers.Hoo)

	log.Log.Info("Hoo routes created :)")

	
}