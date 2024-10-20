package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func LogRequests(c *fiber.Ctx) error {
	method := c.Method()             // Get the HTTP method
	url := c.OriginalURL()           // Get the request URL
	userAgent := c.Get("User-Agent") // Get the User-Agent header

	forwardedIP := c.Get("X-Forwarded-For")
	if forwardedIP == "" {
		forwardedIP = c.IP() // Fallback to direct IP
	} else {
		forwardedIP = c.IP() + " (forwarded for " + forwardedIP + ")"
	}

	// If the route is /rpc, concat the bitcoin method to the url
	if url == "/rpc" {
		var rpcReq RpcRequest
		if err := c.BodyParser(&rpcReq); err == nil {
			url += " " + rpcReq.Method
		}
	}

	log.Printf("Request from %s - %s %s (User-Agent: %s)", forwardedIP, method, url, userAgent)

	err := c.Next()

	statusCode := c.Response().StatusCode()

	if err != nil {
		log.Printf("Response to %s - %d %s", forwardedIP, statusCode, err.Error())
	} else {
		log.Printf("Response to %s - %d", forwardedIP, statusCode)
	}

	return err
}

func RateLimiterMiddleware() fiber.Handler {
	return limiter.New()
}

func ApiKeyProtection(c *fiber.Ctx) error {
	if os.Getenv("RPC_ENV") == "development" {
		return c.Next()
	}

	apiKey := c.Get("x-api-key")
	if apiKey != os.Getenv("RPC_API_KEY") {
		log.Printf("Unauthorized access with wrong API key")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	return c.Next()
}
