package server

import (
	"net/http"

	"github.com/afa7789/satsukashii/internal/bigmac"
	"github.com/afa7789/satsukashii/pkg/bitcoin_price"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) chartPage(bpr bitcoin_price.BTCPricesRanged, cd bigmac.ChartData) fiber.Handler {
	return func(c *fiber.Ctx) error {

		return c.Status(http.StatusOK).Render("chart.html", fiber.Map{
			"SizeSVGH":  cd.SizeHeight,
			"SizeSVGW":  cd.SizeWidth,
			"SpaceDiff": cd.SpaceDiff,
			// normal price lines
			"X1Array": cd.X1Array,
			"Y1Array": cd.Y1Array,
			// price in satoshi lines
			"Y1ArraySatoshi": cd.Y1ArraySatoshi,
			// max price
			"MaxPrice": cd.MaxPrice,
			// max price in satoshi
			"MaxPriceSatoshi": cd.MaxPriceSatoshi,
			"BTC_PRICES":      bpr.Prices,
			"BTC_DATES":       bpr.Dates,
		})
	}
}
