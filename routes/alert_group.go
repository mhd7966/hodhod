package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mhd7966/hodhod/controllers"
	"github.com/mhd7966/hodhod/log"
)

func AlertGroupRouter(app fiber.Router) {

	api := app.Group("/alert-group")

	api.Get("", controllers.GetAlertGroups)
	api.Get("/:id", controllers.GetAlertGroup)
	api.Post("", controllers.NewAlertGroup)
	api.Put("/:id", controllers.UpdateAlertGroup)
	api.Delete("/:id", controllers.DeleteAlertGroup)

	api.Put("/:id/alert", controllers.UpdateAlertGroupAlerts)

	api.Get("/:id/project", controllers.GetProjectsAlertGroup)
	api.Get("/:id/project/:pid", controllers.GetProjectAlertGroup)
	api.Post("/:id/project", controllers.NewProjectAlertGroup)
	api.Delete("/:id/project/:pid", controllers.DeleteProjectAlertGroup)

	api.Get("/:id/project/:pid/contact", controllers.GetContactsProjectAlertGroup)
	api.Get("/:id/project/:pid/contact/:cpid", controllers.GetContactProjectAlertGroup)
	api.Post("/:id/project/:pid/contact", controllers.NewContactProjectAlertGroup)
	api.Put("/:id/project/:pid/contact/:cpid", controllers.UpdateContactProjectAlertGroup)
	api.Delete("/:id/project/:pid/contact/:cpid", controllers.DeleteContactProjectAlertGroup)

	log.Log.Info("AlertGroup routes created :)")

}
