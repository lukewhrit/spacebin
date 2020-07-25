package middlewares

import (
	"github.com/gofiber/fiber"
	"github.com/spacebin-org/curiosity/config"
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
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Set("Cache-Control", "max-age=31536000")

		if config.GetUseCSP() == true {
			c.Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; base-uri 'none'; form-action 'none';")
		}

		// Go to next middleware:
		c.Next()
	}
}
