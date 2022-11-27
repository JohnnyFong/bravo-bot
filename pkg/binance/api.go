package binance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/bravo-bot/pkg/binance/utils"
)

type ListenKeyRes struct {
	ListenKey string `json:"listenKey"`
}

type PositionRisk struct {
	Symbol           string `json:"symbol"`
	PositionAmount   string `json:"positionAmt"`
	EntryPrice       string `json:"entryPrice"`
	MarkPrice        string `json:"markPrice"`
	UnRealizedProfit string `json:"unRealizedProfit"`
	LiquidationPrice string `json:"liquidationPrice"`
	Leverage         string `json:"leverage"`
	MaxNotionalValue string `json:"maxNotionalValue"`
	MarginType       string `json:"marginType"`
	IsolatedMargin   string `json:"isolatedMargin"`
	IsAutoAddMargin  string `json:"isAutoAddMargin"`
	PositionSide     string `json:"positionSide"`
	Notional         string `json:"notional"`
	IsolatedWallet   string `json:"isolatedWallet"`
	UpdateTime       int64  `json:"updateTime"`
}

type OpenOrder struct {
	AvgPrice      string `json:"avgPrice"`
	ClientOrderId string `json:"clientOrderId"`
	CumQuote      string `json:"cumQuote"`
	ExecutedQty   string `json:"executedQty"`
	OrderId       int64  `json:"orderId"`
	OrigQty       string `json:"origQty"`
	OrigType      string `json:"origType"`
	Price         string `json:"price"`
	ReduceOnly    bool   `json:"reduceOnly"`
	Side          string `json:"side"`
	PositionSide  string `json:"positionSide"`
	Status        string `json:"status"`
	StopPrice     string `json:"stopPrice"`
	ClosePosition bool   `json:"closePosition"`
	Symbol        string `json:"symbol"`
	Time          int64  `json:"time"`
	TimeInForce   string `json:"timeInForce"`
	Type          string `json:"type"`
	ActivePrice   string `json:"activatePrice"`
	PriceRate     string `json:"priceRate"`
	UpdateTime    int64  `json:"updateTime"`
	WorkingType   string `json:"workingType"`
	PriceProtect  bool   `json:"priceProtect"`
}

func GetListenKey(apiKey string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", os.Getenv("BASE_URL")+"/fapi/v1/listenKey", nil)
	req.Header.Set("X-MBX-APIKEY", apiKey)
	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to get listenKey: %s\n", err)
		return ""
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var rs ListenKeyRes
		err := json.Unmarshal(data, &rs)
		if err != nil {
			fmt.Printf("Failed to transform listenKey: %s\n", err)
		}
		fmt.Printf("ListenKey: %s\n", rs.ListenKey)
		return rs.ListenKey
	}
}

func GetPositionRisk(apiKey string) (*http.Response, error) {
	query := "recWindow=60000"
	q := utils.ApiQuery(query, os.Getenv("SECRET_KEY"))
	client := &http.Client{}
	req, _ := http.NewRequest("GET", os.Getenv("BASE_URL")+"/fapi/v1/positionRisk?"+q, nil)
	req.Header.Set("X-MBX-APIKEY", apiKey)
	return client.Do(req)
}

func GetOpenOrder(apiKey string) (*http.Response, error) {
	query := "recWindow=60000"
	q := utils.ApiQuery(query, os.Getenv("SECRET_KEY"))
	client := &http.Client{}
	req, _ := http.NewRequest("GET", os.Getenv("BASE_URL")+"/fapi/v1/openOrders?"+q, nil)
	req.Header.Set("X-MBX-APIKEY", apiKey)
	return client.Do(req)
}
