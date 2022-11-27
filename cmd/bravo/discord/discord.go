package discord

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"

	"github.com/bravo-bot/pkg/binance"
	"github.com/bravo-bot/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

var Pnl = false

func InitDiscordBot(wg *sync.WaitGroup) {
	defer wg.Done()
	discord.Init(messageCreate)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!pos" && m.ChannelID == os.Getenv("BOT_CHANNEL") {
		response, err := binance.GetPositionRisk(os.Getenv("API_KEY"))
		if err != nil {
			fmt.Printf("Failed to get positionRisk: %s\n", err)
			discord.SendMessage("No position", os.Getenv("BOT_CHANNEL"))
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			var rs []binance.PositionRisk
			err := json.Unmarshal(data, &rs)
			if err != nil {
				fmt.Printf("Failed to transform positionRisk: %s\n", err)
				discord.SendMessage("No position", os.Getenv("BOT_CHANNEL"))
			} else {
				f := []*discordgo.MessageEmbedField{}
				for _, ps := range rs {
					if ps.EntryPrice != "0.0" {
						f = append(f, []*discordgo.MessageEmbedField{
							{
								Name:  "\u200B",
								Value: "\u200B",
							},
							{
								Name:   "Symbol",
								Value:  ps.Symbol,
								Inline: true,
							},
							{
								Name:   "Position Side",
								Value:  ps.PositionSide,
								Inline: true,
							},
							{
								Name:   "Leverage",
								Value:  ps.Leverage,
								Inline: true,
							},
							{
								Name:   "Entry Price",
								Value:  ps.EntryPrice,
								Inline: true,
							},
							{
								Name:   "Market Price",
								Value:  ps.MarkPrice,
								Inline: true,
							},
						}...)
						if Pnl {
							f = append(f, []*discordgo.MessageEmbedField{
								{
									Name:  "PNL",
									Value: ps.UnRealizedProfit,
								},
							}...)
						}
					}
				}

				if len(f) > 0 {
					embed := discordgo.MessageEmbed{
						Type:   discordgo.EmbedTypeRich,
						Title:  "Current Positions",
						Color:  0x0099FF,
						Fields: f,
					}
					discord.SendEmbedMessage(&embed, os.Getenv("BOT_CHANNEL"))
				} else {
					discord.SendMessage("No position", os.Getenv("BOT_CHANNEL"))
				}
			}
		}
	}

	if m.Content == "!pnl" && m.ChannelID == os.Getenv("CMD_CHANNEL") {
		Pnl = !Pnl
		discord.SendMessage("PNL mode - "+strconv.FormatBool(Pnl), os.Getenv("CMD_CHANNEL"))
	}

	if m.Content == "!order" && m.ChannelID == os.Getenv("BOT_CHANNEL") {
		response, err := binance.GetOpenOrder(os.Getenv("API_KEY"))
		if err != nil {
			fmt.Printf("Failed to get openOrder: %s\n", err)
			discord.SendMessage("No position", os.Getenv("BOT_CHANNEL"))
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			var rs []binance.OpenOrder
			err := json.Unmarshal(data, &rs)
			if err != nil {
				fmt.Printf("Failed to transform openOrder: %s\n", err)
				discord.SendMessage("No open order", os.Getenv("BOT_CHANNEL"))
			} else {
				f := []*discordgo.MessageEmbedField{}
				for _, ps := range rs {
					f = append(f, []*discordgo.MessageEmbedField{
						{
							Name:  "\u200B",
							Value: "\u200B",
						},
						{
							Name:   "Symbol",
							Value:  ps.Symbol,
							Inline: true,
						},
						{
							Name:   "Type",
							Value:  ps.Type,
							Inline: true,
						},
						{
							Name:   "Position Side",
							Value:  ps.Side,
							Inline: true,
						},
						{
							Name:   "Entry Price",
							Value:  ps.Price,
							Inline: true,
						},
						{
							Name:   "Stop Price",
							Value:  ps.StopPrice,
							Inline: true,
						},
					}...)
				}

				if len(f) > 0 {
					embed := discordgo.MessageEmbed{
						Type:   discordgo.EmbedTypeRich,
						Title:  "Current Open Orders",
						Color:  0x0099FF,
						Fields: f,
					}
					discord.SendEmbedMessage(&embed, os.Getenv("BOT_CHANNEL"))
				} else {
					discord.SendMessage("No open order", os.Getenv("BOT_CHANNEL"))
				}
			}
		}
	}
}
