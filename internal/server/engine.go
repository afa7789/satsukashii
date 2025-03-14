package server

import (
	"fmt"

	"github.com/gofiber/template/html/v2"
)

func addEngineFunctions(engine *html.Engine) *html.Engine {

	// Basic math functions with consistent types
	engine.AddFunc("add", func(a, b interface{}) int {
		var numA, numB int

		switch v := a.(type) {
		case int:
			numA = v
		case int64:
			numA = int(v)
		case float64:
			numA = int(v)
		default:
			numA = 0
		}

		switch v := b.(type) {
		case int:
			numB = v
		case int64:
			numB = int(v)
		case float64:
			numB = int(v)
		default:
			numB = 0
		}

		return numA + numB
	})

	engine.AddFunc("iterate", func(count int) []int {
		items := make([]int, count)
		for i := range items {
			items[i] = i
		}
		return items
	})

	engine.AddFunc("multiply", func(a, b interface{}) float64 {
		var numA, numB float64

		switch v := a.(type) {
		case int:
			numA = float64(v)
		case int64:
			numA = float64(v)
		case float64:
			numA = v
		default:
			numA = 0
		}

		switch v := b.(type) {
		case int:
			numB = float64(v)
		case int64:
			numB = float64(v)
		case float64:
			numB = v
		default:
			numB = 0
		}

		return numA * numB
	})

	engine.AddFunc("divide", func(a, b interface{}) float64 {
		var numA, numB float64

		switch v := a.(type) {
		case int:
			numA = float64(v)
		case int64:
			numA = float64(v)
		case float64:
			numA = v
		default:
			numA = 0
		}

		switch v := b.(type) {
		case int:
			numB = float64(v)
		case int64:
			numB = float64(v)
		case float64:
			numB = v
		default:
			numB = 0
		}

		if numB == 0 {
			return 0
		}
		return numA / numB
	})

	engine.AddFunc("subtract", func(a, b interface{}) int {
		var numA, numB int

		switch v := a.(type) {
		case int:
			numA = v
		case int64:
			numA = int(v)
		case float64:
			numA = int(v)
		default:
			numA = 0
		}

		switch v := b.(type) {
		case int:
			numB = v
		case int64:
			numB = int(v)
		case float64:
			numB = int(v)
		default:
			numB = 0
		}

		return numA - numB
	})

	engine.AddFunc("formatPrice", func(price float64) string {
		return fmt.Sprintf("%.2f", price)
	})

	engine.AddFunc("formatSats", func(sats float64) string {
		return fmt.Sprint(int(sats + 0.5))
	})

	engine.AddFunc("toInt", func(f interface{}) int {
		switch v := f.(type) {
		case int:
			return v
		case int64:
			return int(v)
		case float64:
			return int(v)
		default:
			return 0
		}
	})

	// Specialized helpers for SVG calculations - all return int for consistency
	engine.AddFunc("calculateGridY", func(height interface{}, index int) int {
		var h float64
		switch v := height.(type) {
		case int:
			h = float64(v)
		case int64:
			h = float64(v)
		case float64:
			h = v
		default:
			h = 0
		}

		spacing := h / 4.0
		return 50 + int(float64(index)*spacing)
	})

	engine.AddFunc("calculateGridX", func(width interface{}, index int, divisions int) int {
		var w float64
		switch v := width.(type) {
		case int:
			w = float64(v)
		case int64:
			w = float64(v)
		case float64:
			w = v
		default:
			w = 0
		}

		spacing := w / float64(divisions)
		return 50 + int(float64(index)*spacing)
	})

	engine.AddFunc("calculateYAxisLabel", func(svgHeight interface{}, index int, divisions int) int {
		var h float64
		switch v := svgHeight.(type) {
		case int:
			h = float64(v)
		case int64:
			h = float64(v)
		case float64:
			h = v
		default:
			h = 0
		}

		height := h - 100 // Same as $height = subtract .SizeSVGH 100
		spacing := height / float64(divisions)
		return int(float64(int(h)-50) - (float64(index) * spacing))
	})

	// Helper for comparing values in templates
	engine.AddFunc("lt", func(a, b int) bool {
		return a < b
	})

	engine.AddFunc("len", func(s interface{}) int {
		switch v := s.(type) {
		case []string:
			return len(v)
		case []float64:
			return len(v)
		case []int64:
			return len(v)
		default:
			return 0
		}
	})
	return engine
}
