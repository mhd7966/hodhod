package routes

import (
	"github.com/abr-ooo/hodhod/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/abr-ooo/hodhod/log"

)

func AlertRouter(app fiber.Router) {

	api := app.Group("/alert")

	api.Get("", controllers.GetAlerts)
	api.Get("/:id", controllers.GetAlert)
	api.Post("", controllers.NewAlert)
	api.Put("/:id", controllers.UpdateAlert)
	api.Delete("/:id", controllers.DeleteAlert)

	log.Log.Info("Alert routes created :)")

}
