package routes

import (
	"github.com/abr-ooo/go-pkgs"
	"github.com/gofiber/fiber/v2"
	"github.com/abr-ooo/hodhod/log"
)


func MainRouter(app fiber.Router) {
	api := app.Group("/v0", gopkgs.Auth)

	LogRouter(api)
	HooRouter(api)
	AlertRouter(api)
	ContactRouter(api)
	AlertGroupRouter(api)

	log.Log.Info("All routes created successfully :)")

}
