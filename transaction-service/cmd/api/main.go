package main

import (
	"github.com/joho/godotenv"
	"github.com/grayDorian1/Entain/internal/app"
)

// @title           Entain Transaction API
// @version         1.0
// @description     Transaction service
// @host            localhost:8080
// @BasePath        /
func main() {
	godotenv.Load()
	app.Run()
}