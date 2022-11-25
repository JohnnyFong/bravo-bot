package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/bravo-bot/cmd/bravo/binance"
	"github.com/bravo-bot/cmd/bravo/discord"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var wg sync.WaitGroup
	wg.Add(4)

	go binance.GetListenKey(&wg)
	go binance.ListenWebsocket(&wg)
	go binance.ReceiveWebsocket(&wg)
	fmt.Println("Binance Service Started")
	go discord.InitDiscordBot(&wg)

	wg.Wait()
}
