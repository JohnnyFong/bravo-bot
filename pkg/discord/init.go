package discord

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

var dg *discordgo.Session

func Init(fn func(*discordgo.Session, *discordgo.MessageCreate)) {
	var err error
	dg, err = discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(fn)
	fmt.Println("Connected to discord!")
	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	// fmt.Println("Bot is now running. Press CTRL-C to exit.")
	// sc := make(chan os.Signal, 1)
	// signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	// <-sc

	// // Cleanly close down the Discord session.
	// dg.Close()
}
