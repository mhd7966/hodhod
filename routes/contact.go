package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mhd7966/hodhod/controllers"
	"github.com/mhd7966/hodhod/log"
)

func ContactRouter(app fiber.Router) {

	api := app.Group("/contact")

	api.Get("/", controllers.GetUserContacts)
	api.Get("/:id", controllers.GetContact)
	api.Post("/", controllers.NewContact)
	api.Put("/:id", controllers.UpdateContact)
	api.Delete("/:id", controllers.DeleteContact)

	log.Log.Info("Contact routes created :)")

}
