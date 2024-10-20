package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

type RpcRequest struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

var allowedMethods = map[string]bool{
	"getblockchaininfo": true,
	"getblockhash":      true,
	"getblock":          true,
	// Add other allowed methods here
}

func main() {
	LoadEnv()
	app := fiber.New()

	// Middleware
	app.Use(LogRequests)
	app.Use(RateLimiterMiddleware())
	app.Use(ApiKeyProtection)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!")
	})

	app.Post("/rpc", func(c *fiber.Ctx) error {
		// Parse the incoming request
		var rpcReq RpcRequest
		if err := c.BodyParser(&rpcReq); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		// Check if the method is allowed
		if !allowedMethods[rpcReq.Method] {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Method not allowed"})
		}

		// Call the Rpc function
		result, err := Rpc(rpcReq.Method, rpcReq.Params)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// Return the result
		return c.JSON(result)
	})

	// Start the server
	port := 3000

	fullchain := "/etc/letsencrypt/live/home.tomx.sh/fullchain.pem"
	privkey := "/etc/letsencrypt/live/home.tomx.sh/privkey.pem"

	// Check if the files exist
	_, err := os.Stat(fullchain)
	if err != nil {
		log.Fatal("Fullchain file not found at", fullchain)
	}

	_, err = os.Stat(privkey)
	if err != nil {
		log.Fatal("Privkey file not found at", privkey)
	}

	fmt.Printf("Server is running on port %d\n", port)

	log.Fatal(app.ListenTLS(fmt.Sprintf(":%d", port), fullchain, privkey))
	// TODO: use http for development
}
