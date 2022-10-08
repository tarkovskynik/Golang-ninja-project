package main

import "github.com/tarkovskynik/Golang-ninja-project/internal/app"

var configDir = "config"

// @title File Manager App API
// @version 1.0
// @description API Server for File Manager Application
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey JWT
// @in header
// @name Authorization
func main() {
	app.Run(configDir)
}
