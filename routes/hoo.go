package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mhd7966/hodhod/controllers"
	"github.com/mhd7966/hodhod/log"
)

func HooRouter(app fiber.Router) {

	api := app.Group("/hoo")

	api.Post("", controllers.Hoo)

	log.Log.Info("Hoo routes created :)")

}
