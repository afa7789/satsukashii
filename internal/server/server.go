package server

import (
	"fmt"
	"log"
	"strings"

	"github.com/afa7789/satsukashii/internal/bigmac"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

// server client

// multilanguage setup, if added to the first element after / the language, change the multilanguage used to render the HTMLs
// will have to use a template engine for the above of course

// Server is the definition of a REST server based on Gin
type Server struct {
	router *fiber.App
}

// return the .js if needed on request.
func New() *Server {
	server := &Server{
		// reps: si.Reps,
	}

	// https://github.com/gofiber/template
	engine := html.New("./web/templates", "")
	// create functions that will be used in the template
	// those functions are used to do the pagination.
	engine.AddFunc("add", func(a, b int) int {
		return a + b
	})
	engine.AddFunc("sub", func(a, b int) int {
		return a - b
	})

	// Reload the templates on each render, good for development
	engine.Reload(true) // Optional. Default: false
	r := fiber.New(fiber.Config{
		Views:             engine,
		EnablePrintRoutes: false,
	})

	cd, err := bigmac.GenerateChartData(
		500, // height
		700, // width
	)
	if err != nil {
		log.Fatal(err)
	}

	// ================ROUTES====================
	// Static Files
	r.Static("/public", "./web/static")

	r.Get("/", server.chartPage(cd))

	server.router = r
	return server
}

// Start starts the server
func (s *Server) Start(port int) {
	err := s.router.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		// Using this error treatment to try again on next port
		if strings.Contains(err.Error(), "address already in use") {
			fmt.Println("")
			log.Printf("PORT ALREADY IN USE::%d", port)
			port++
			log.Printf("TRYING NEXT PORT:%d\n", port)
			s.Start(port)
		} else {
			panic(err)
		}
	}
}
