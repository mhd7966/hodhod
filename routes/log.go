package routes

import (
	"github.com/abr-ooo/hodhod/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/abr-ooo/hodhod/log"
)

func LogRouter(app fiber.Router) {

	api := app.Group("/log")

	api.Get("", controllers.Log)

	log.Log.Info("Hoo routes created :)")

	
}