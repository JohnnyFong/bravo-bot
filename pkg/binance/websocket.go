package binance

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type BinanceOrderEvent struct {
	EventType       string                 `json:"e"`
	EventTime       int64                  `json:"E"`
	TransactionTime int64                  `json:"T"`
	Object          map[string]interface{} `json:"o"`
}

type BinanceOrderObject struct {
	Symbol                         string `json:"s"`
	ClientOrderID                  string `json:"c"`
	Side                           string `json:"S"`
	OrderType                      string `json:"o"`
	TimeInForce                    string `json:"f"`
	OriginalQuantity               string `json:"q"`
	OriginalPrice                  string `json:"p"`
	AveragePrice                   string `json:"ap"`
	StopPrice                      string `json:"sp"`
	ExecutionType                  string `json:"x"`
	OrderStatus                    string `json:"X"`
	OrderID                        int64  `json:"i"`
	OrderLastFilledQuantity        string `json:"l"`
	OrderFilledAccumulatedQuantity string `json:"z"`
	LastFilledPrice                string `json:"L"`
	OrderTradeTime                 int64  `json:"T"`
	TradeID                        int64  `json:"t"`
	BidsNational                   string `json:"b"`
	MakerSide                      bool   `json:"m"`
	ReduceOnly                     bool   `json:"R"`
	StopPriceWorkingType           string `json:"wt"`
	OriginalOrderType              string `json:"ot"`
	PositionSide                   string `json:"ps"`
	ConditionalOrder               bool   `json:"cp"`
	RealizedProfit                 string `json:"rp"`
}

func ListenWebsocket(listenKey string, ch chan BinanceOrderObject) {
	fmt.Println("connecting to websocket")
	c, _, err := websocket.DefaultDialer.Dial("wss://fstream.binance.com/ws/"+listenKey, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	// Clean up on exit from this goroutine
	defer c.Close()
	// Loop reading messages. Send each message to the channel.
	fmt.Println("connected! listening to message")
	for {
		_, m, err := c.ReadMessage()
		if err != nil {
			fmt.Println("error listening", err)
		}
		var res BinanceOrderEvent
		err = json.Unmarshal(m, &res)
		if err != nil {
			fmt.Println("unmarshal failed", err)
		}
		if res.EventType == "ORDER_TRADE_UPDATE" {
			jsonString, _ := json.Marshal(res.Object)
			var obj BinanceOrderObject
			json.Unmarshal(jsonString, &obj)
			if (obj.AveragePrice != "0" || obj.OriginalPrice != "0" || obj.StopPrice != "0") && obj.ExecutionType != "EXPIRED" {
				ch <- obj
			}
		}
	}
}
