package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mhd7966/hodhod/controllers"
	"github.com/mhd7966/hodhod/log"
)

func LogRouter(app fiber.Router) {

	api := app.Group("/log")

	api.Get("", controllers.Log)

	log.Log.Info("Hoo routes created :)")

}
