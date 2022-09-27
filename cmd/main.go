package main

import "github.com/tarkovskynik/Golang-ninja-project/internal/app"

var configDir = "config"

func main() {
	app.Run(configDir)
}
