package main

import (
	"BE-6/src/api/route"
	"BE-6/src/config/database"
	"BE-6/src/config/env"
	"BE-6/src/config/storage"

	"github.com/joho/godotenv"
)

func main() {
	// load env file
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// migrate gorm
	database.MigrateMySQL()

	storage.InitGDrive()

	app := route.InitServer()
	app.Logger.Fatal(app.Start(":" + env.GetServerEnv()))
}
