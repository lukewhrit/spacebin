package middlewares

import (
	"github.com/gofiber/fiber"
)

// SecurityHeaders sets various headers related to security
func SecurityHeaders() func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		// Set some security headers:
		c.Set("X-Download-Options", "noopen")
		c.Set("X-DNS-Prefetch-Control", "off")
		c.Set("X-Frame-Options", "SAMEORIGIN")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("Referrer-Policy", "no-referrer-when-downgrade")
		c.Set("Content-Security-Policy", "default-src 'self'; img-src 'self' data:; object-src 'none'; script-src 'self'; frame-ancestors 'none'; base-uri 'none'; form-action 'none';")
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Set("Cache-Control", "max-age=31536000")

		// Go to next middleware:
		c.Next()
	}
}
