package main

import (
    "github.com/suk-chanthea/ezra/bootstrap"
)

func main() {
	// Load config
	config := bootstrap.LoadConfig()

	// Connect database
	db := bootstrap.NewDatabase(config)

	// Setup router with DB and secret key
	router := bootstrap.SetupRouter(db, config.SecretKey)

	// Start server
	router.Run(":" + config.Port)
}
