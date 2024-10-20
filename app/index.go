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
	"getblockchaininfo":    true,
	"getblockhash":         true,
	"getblock":             true,
	"getblockheader":       true,
	"getchaintips":         true,
	"getdifficulty":        true,
	"getmempoolinfo":       true,
	"getrawmempool":        true,
	"getmempoolentry":      true,
	"getrawtransaction":    true,
	"decoderawtransaction": true,
	"getnetworkinfo":       true,
	"getconnectioncount":   true,
	"getpeerinfo":          true,
	"estimatesmartfee":     true,
	"getmininginfo":        true,
	"getblocktemplate":     true,
	"getblockstats":        true,
	"getblocksubsidy":      true,
	"getchaintxstats":      true,
	"gettxoutsetinfo":      true,
	"gettxout":             true,
	"gettxoutproof":        true,
	"getblockfilter":       true,
	"getbestblockhash":     true,
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

	app.Get("/allowed-methods", func(c *fiber.Ctx) error {
		return c.JSON(allowedMethods)
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
	fmt.Printf("Server is running on port %d\n", port)

	if os.Getenv("RPC_ENV") == "development" {
		// Use http in development
		log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))

	} else {
		// Use https in production
		fullchain := os.Getenv("SSL_CERTIFICATES_PATH") + "/fullchain.pem"
		privkey := os.Getenv("SSL_CERTIFICATES_PATH") + "/privkey.pem"

		// Check if the files exist
		_, err := os.Stat(fullchain)
		if err != nil {
			log.Fatal("Fullchain file not found at ", fullchain)
		}

		_, err = os.Stat(privkey)
		if err != nil {
			log.Fatal("Privkey file not found at ", privkey)
		}

		log.Fatal(app.ListenTLS(fmt.Sprintf(":%d", port), fullchain, privkey))
	}
}
