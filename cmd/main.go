package main

import "github.com/ptsypyshev/http-ratelimiter/internal/app"

func main() {
	a := app.NewApp() // Create App instance
	a.Run()           // Run App
}
