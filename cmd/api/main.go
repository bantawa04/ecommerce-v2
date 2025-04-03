package main

import (
	"beautyessentials.com/internal/bootstrap"
)

func main() {
	// Build and start the application using Fx
	app := bootstrap.BuildApp()
	app.Run()
}
