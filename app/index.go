package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func main() {
	LoadEnv()

	app := fiber.New()

	// Middleware for API key protection
	app.Use(func(c *fiber.Ctx) error {
		if os.Getenv("ENV") == "local" {
			return c.Next()
		}

		apiKey := c.Get("x-api-key")
		if apiKey != os.Getenv("API_KEY") {
			log.Printf("Unauthorized access with wrong API key")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}
		return c.Next()
	})

	// Define routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!")
	})

	app.Get("/get-blockchain-info", func(c *fiber.Ctx) error {
		result, err := Rpc("getblockchaininfo", nil)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(result)
	})

	app.Get("/get-network-info", func(c *fiber.Ctx) error {
		result, err := Rpc("getnetworkinfo", nil)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(result)
	})

	app.Get("/get-block-count", func(c *fiber.Ctx) error {
		result, err := Rpc("getblockcount", nil)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(result)
	})

	app.Get("/get-block-hash", func(c *fiber.Ctx) error {
		blockHeight := c.Query("blockHeight")
		if blockHeight == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Block height is required"})
		}
		height, _ := strconv.Atoi(blockHeight)
		result, err := Rpc("getblockhash", []interface{}{height})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(result)
	})

	app.Get("/get-block", func(c *fiber.Ctx) error {
		blockHash := c.Query("blockHash")
		if blockHash == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Block hash is required"})
		}
		result, err := Rpc("getblock", []interface{}{blockHash})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(result)
	})

	// Add more routes based on your original code...

	// Start the server
	port := 3000
	fmt.Printf("Server is running on port %d\n", port)
	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
