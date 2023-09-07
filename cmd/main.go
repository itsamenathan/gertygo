package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"main/internal/discord"
	"main/internal/finnhub"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	discord.Init(os.Getenv("DISCORD_TOKEN"))
	finnhub.Init(os.Getenv("FINNHUB_TOKEN"))

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running...  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discord.Close()
}
