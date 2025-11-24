package main

import (
	"log"
	"os"
)

var DebugMode = os.Getenv("DEBUG") == "1"

func main() {
	if DebugMode {
		log.Println("Debug mode enabled")
	}
	
	// Check if LICHESS_TOKEN environment variable is set
	lichessToken := os.Getenv("LICHESS_TOKEN")
	if lichessToken != "" {
		// Run as Lichess bot
		log.Println("Starting in Lichess bot mode...")
		StartLichessBot(lichessToken)
	} else {
		// Run as local server
		log.Println("Starting in local server mode...")
		log.Println("(Set LICHESS_TOKEN environment variable to run as Lichess bot)")
		StartServer()
	}
	// PerformanceTest()
}
