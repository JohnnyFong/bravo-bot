package binance

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	d "github.com/bravo-bot/cmd/bravo/discord"
	"github.com/bravo-bot/pkg/binance"
	"github.com/bravo-bot/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

var listenKey string
var binanceChannel = make(chan binance.BinanceOrderObject)

var loc, _ = time.LoadLocation("Asia/Singapore")

func GetListenKey(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		listenKey = binance.GetListenKey(os.Getenv("API_KEY"))
		time.Sleep(50 * time.Minute)
	}
}

func ListenWebsocket(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		if listenKey != "" {
			binance.ListenWebsocket(listenKey, binanceChannel)
		} else {
			fmt.Println("listenKey is empty, wait for 10 second")
			time.Sleep(10 * time.Second)
		}
	}
}

func ReceiveWebsocket(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		obj := <-binanceChannel
		go sendEmbedMessage(obj)
	}
}

func sendEmbedMessage(o binance.BinanceOrderObject) {

	var eType string
	priceType := "Entry Price"
	rp, _ := strconv.ParseFloat(strings.TrimSpace(o.RealizedProfit), 64)
	if o.RealizedProfit == "0" {
		eType = o.ExecutionType
	} else if rp > 0 {
		eType = "TAKE PROFIT"
	} else {
		eType = "STOP LOSS"
	}

	var c int
	if o.Side == "SELL" {
		c = 0xFF0000
	} else {
		c = 0x00FF00
	}

	if eType == "TAKE PROFIT" {
		priceType = "Exit Price"
		c = 0x0099FF
	} else if eType == "CANCELED" {
		c = 0x808080
	} else if eType == "STOP LOSS" {
		priceType = "Exit Price"
		c = 0xFFFF00
	}
	var price string
	if o.AveragePrice != "0" {
		price = o.AveragePrice
	} else if o.OriginalPrice != "0" {
		price = o.OriginalPrice
	} else if o.StopPrice != "0" {
		priceType = "Exit Price"
		price = o.StopPrice
	} else {
		price = "0"
	}
	f := []*discordgo.MessageEmbedField{}
	f = append(f, []*discordgo.MessageEmbedField{
		{
			Name:  priceType,
			Value: price,
		},
		{
			Name:  "Time",
			Value: time.Now().In(loc).Format("2006-01-02 15:04:05"),
		},
	}...)
	if d.Pnl {
		f = append(f, []*discordgo.MessageEmbedField{
			{
				Name:  "PNL",
				Value: o.RealizedProfit,
			},
		}...)
	}
	embed := discordgo.MessageEmbed{
		Type:   discordgo.EmbedTypeRich,
		Title:  eType + " " + o.OrderType + " " + o.Side + " - " + o.Symbol,
		Color:  c,
		Fields: f,
	}
	discord.SendEmbedMessage(&embed, os.Getenv("POS_CHANNEL"))
}
