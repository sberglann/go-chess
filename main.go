package main

import (
	"flag"
	"log"
	"os"
)

var DebugMode = os.Getenv("DEBUG") == "1"

func main() {
	if DebugMode {
		log.Println("Debug mode enabled")
	}

	// Define all flags
	localFlag := flag.Bool("local", false, "Run local server")
	benchFlag := flag.Bool("bench", false, "Run performance test")
	cpuprofileFlag := flag.String("cpuprofile", "", "write cpu profile to file")
	
	// Parse flags once
	flag.Parse()

	if *benchFlag {
		log.Println("Running performance test...")
		PerformanceTest(*cpuprofileFlag)
		return
	}

	if *localFlag {
		log.Println("Starting in local server mode...")
		StartServer()
		return
	}

	// Default: Check if LICHESS_TOKEN environment variable is set
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
}
