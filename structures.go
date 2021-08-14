package main

import "sync"

type ResponseCur struct {
	ID          string `json:"id"`
	FullName    string `json:"fullName"`
	FeeCurrency string `json:"feeCurrency"`
	Ask         string `json:"ask"`
	Bid         string `json:"bid"`
	Last        string `json:"last"`
	Open        string `json:"open"`
	Low         string `json:"low"`
	High        string `json:"high"`
}

type BTCResponse struct {
	Params struct {
		Ask  string `json:"ask"`
		Bid  string `json:"bid"`
		Last string `json:"last"`
		Open string `json:"open"`
		Low  string `json:"low"`
		High string `json:"high"`
	} `json:"params"`
}

type RequestCurrSymbol struct {
	Params struct {
		Currency string `json:"currency"`
		Symbol   string `json:"symbol"`
	} `json:"params"`
	Method string `json:"method"`
}
type GetSymbolResp struct {
	Result struct {
		BaseCurrency string `json:"baseCurrency"`
		FeeCurrency  string `json:"feeCurrency"`
	} `json:"result"`
}
type GetCurrencyResp struct {
	Result struct {
		ID       string `json:"id"`
		FullName string `json:"fullName"`
	} `json:"result"`
}

type FinalResult struct {
	ID          string `json:"id"`
	FullName    string `json:"fullName"`
	Ask         string `json:"ask"`
	Bid         string `json:"bid"`
	Last        string `json:"last"`
	Open        string `json:"open"`
	Low         string `json:"low"`
	High        string `json:"high"`
	FeeCurrency string `json:"feeCurrency"`
}

type CurrRes struct {
	Data map[string]FinalResult
	mu   sync.Mutex
}

type CurrencySymbols struct {
	Symbolscur []string
}
