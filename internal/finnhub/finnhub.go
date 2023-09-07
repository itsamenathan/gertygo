package finnhub

import (
	"context"
	"errors"

	"github.com/Finnhub-Stock-API/finnhub-go/v2"
)

var finnhubClient *finnhub.DefaultApiService

func Init(token string) error {
	cfg := finnhub.NewConfiguration()
	cfg.AddDefaultHeader("X-Finnhub-Token", token)
	finnhubClient = finnhub.NewAPIClient(cfg).DefaultApi

	return nil
}

// btc to usd use symbol BINANCE:BTCUSDT
func GetStockQuote(symbol string) (finnhub.Quote, error) {
	res, _, err := finnhubClient.Quote(context.Background()).Symbol(symbol).Execute()
	if err != nil {
		return res, err
	}
	if res.GetC() == 0 {
		err = errors.New("couldn't get stock quote data returned 0")
		return res, err
	}
	return res, err
}
