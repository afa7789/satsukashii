package server

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) indexPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).Render("index.html", fiber.Map{
			"Title":      "afa7789 - Computer Wizzard Tech Blog",
			"MainHeader": true,
		})
	}
}
