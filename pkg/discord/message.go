package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func SendEmbedMessage(embed *discordgo.MessageEmbed, ch string) {
	_, err := dg.ChannelMessageSendEmbed(ch, embed)
	if err != nil {
		fmt.Println(err)
	}
}

func SendMessage(msg string, ch string) {
	_, err := dg.ChannelMessageSend(ch, msg)
	if err != nil {
		fmt.Println(err)
	}
}
