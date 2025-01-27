package main

import (
	"log"
)

func main() {
	app := NewFileBundlerApp()
	if err := app.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
