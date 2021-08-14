package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/websocket"
)

const (
	Origin = "https://localhost/"
	Url    = "wss://api.hitbtc.com/api/2/ws"
)

func (c *CurrRes) GetSingleCurrency(key string) (res FinalResult, ok bool) {
	c.mu.Lock()
	res, ok = c.Data[key]
	defer c.mu.Unlock()
	return
}

func (c *CurrRes) AddDataToRes(key string, res FinalResult) {
	c.mu.Lock()
	c.Data[key] = res
	defer c.mu.Unlock()
}

func (c *CurrRes) GetAllCurrencies() (res map[string]FinalResult) {
	c.mu.Lock()
	res = c.Data
	defer c.mu.Unlock()
	return
}

var cur *CurrRes = &CurrRes{}

func WebSocketData(ws *websocket.Conn, req interface{}) (res []byte) {
	b, _ := json.Marshal(&req)
	if _, err := ws.Write(b); err != nil {
		log.Fatal(err)
	}
	size := 50
	var msg = make([]byte, size)
	var n int
	n, err := ws.Read(msg)
	if err != nil {
		log.Fatal(err)
	}
	temp := make([]byte, size)
	m := n
	for n == size {
		n, err = ws.Read(temp)
		msg = append(msg, temp...)
		m += n
	}
	// fmt.Println(temp)
	if err != nil {
		log.Fatal(err)
	}
	return msg[:m]
}

func GetDataSymbol(symbol string) {
	ws, err := websocket.Dial(Url, "", Origin)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	for {
		hitRequest := RequestCurrSymbol{}
		hitRequest.Method = "subscribeTicker"
		hitRequest.Params.Symbol = symbol
		var givenResponse BTCResponse
		_ = json.Unmarshal(WebSocketData(ws, hitRequest), &givenResponse)
		if len(givenResponse.Params.Ask) > 0 {
			if data, ok := cur.GetSingleCurrency(symbol); ok {
				data.Ask = givenResponse.Params.Ask
				data.Bid = givenResponse.Params.Bid
				data.High = givenResponse.Params.High
				data.Last = givenResponse.Params.Last
				data.Low = givenResponse.Params.Low
				data.Open = givenResponse.Params.Open
				cur.AddDataToRes(symbol, data)
			} else {
				fmt.Println("Invalid symbol")
			}
		}
	}
}
func GetData(symbol string) {
	ws, err := websocket.Dial(Url, "", Origin)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	if _, ok := cur.GetSingleCurrency(symbol); !ok {
		var result FinalResult
		getSymbols := RequestCurrSymbol{}

		//get symbols
		getSymbols.Method = "getSymbol"
		getSymbols.Params.Symbol = symbol
		var givenGetSymbolResponse GetSymbolResp
		_ = json.Unmarshal(WebSocketData(ws, getSymbols), &givenGetSymbolResponse)
		result.FeeCurrency = givenGetSymbolResponse.Result.FeeCurrency

		//get base currencies
		getSymbols.Method = "getCurrency"
		getSymbols.Params.Currency = givenGetSymbolResponse.Result.BaseCurrency

		var givenGetCurrencyResponse GetCurrencyResp
		_ = json.Unmarshal(WebSocketData(ws, getSymbols), &givenGetCurrencyResponse)
		result.FullName = givenGetCurrencyResponse.Result.FullName
		result.ID = givenGetCurrencyResponse.Result.ID
		//add resultant data
		cur.AddDataToRes(symbol, result)
		go GetDataSymbol(symbol)
	}

}

func GetAllCurrencyHandler(w http.ResponseWriter, rq *http.Request) {
	response, _ := json.Marshal(cur.GetAllCurrencies())
	io.WriteString(w, string(response))
}

func GetCurrencyHandler(w http.ResponseWriter, rq *http.Request) {
	val, ok := cur.GetSingleCurrency(strings.TrimPrefix(rq.URL.Path, "/currency/"))
	if ok {
		response, _ := json.Marshal(val)
		io.WriteString(w, string(response))
	} else {
		io.WriteString(w, "invalid id")
	}
}

func main() {
	cur.Data = map[string]FinalResult{}
	file, _ := os.Open("config.json")
	currencySymbols := CurrencySymbols{}
	err := json.NewDecoder(file).Decode(&currencySymbols)
	if err != nil {
		fmt.Println("error:", err)
	}
	for _, v := range currencySymbols.Symbolscur {
		go GetData(v)
	}
	http.HandleFunc("/currency/all", GetAllCurrencyHandler)
	http.HandleFunc("/currency/", GetCurrencyHandler)
	log.Fatal(http.ListenAndServe(":8085", nil))
	file.Close()
}
