package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
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
	t := time.Now().UnixNano() / int64(time.Millisecond)
	query := "recvWindow=60000&timestamp=" + strconv.FormatInt(t, 10)
	key := []byte(os.Getenv("SECERT_KEY"))

	sig := hmac.New(sha256.New, key)
	sig.Write([]byte(query))

	client := &http.Client{}
	req, _ := http.NewRequest("GET", os.Getenv("BASE_URL")+"/fapi/v1/positionRisk?recvWindow=60000&timestamp="+strconv.FormatInt(t, 10)+"&signature="+fmt.Sprintf("%x", (sig.Sum(nil))), nil)
	req.Header.Set("X-MBX-APIKEY", apiKey)

	return client.Do(req)
}
