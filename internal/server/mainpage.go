package server

import (
	"net/http"

	"github.com/afa7789/satsukashii/internal/bigmac"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) chartPage(cd bigmac.ChartData) fiber.Handler {
	return func(c *fiber.Ctx) error {

		return c.Status(http.StatusOK).Render("chart.html", fiber.Map{
			"SizeSVGH": cd.SizeHeight,
			"SizeSVGW": cd.SizeWidth,
			// normal price lines
			"X1Array": cd.X1Array,
			"Y1Array": cd.Y1Array,
			// price in satoshi lines
			"Y1ArraySatoshi": cd.Y1ArraySatoshi,
		})
	}
}
