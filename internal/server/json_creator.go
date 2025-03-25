package server

import (
	"github.com/afa7789/satsukashii/internal/bigmac"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) GetChartJSON(data bigmac.ChartData) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() != fiber.MethodGet {
			return c.Status(fiber.StatusMethodNotAllowed).SendString("Method not allowed")
		}
		return c.JSON(data)
	}
}
