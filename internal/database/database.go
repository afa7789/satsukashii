package database

import (
	"database/sql"
	"time"

	"github.com/afa7789/satsukashii/pkg/bitcoin_price"
)

// wrapper of a database
// if we need
// not sure if needed.

func InsertBitcoinPrice(db *sql.DB, price bitcoin_price.BitcoinPrice) error {
	query := `
    INSERT INTO bitcoin_prices (start, open, high, low, close, volume, market_cap, currency_code)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?);
    `
	_, err := db.Exec(query, price.Start, price.Open, price.High, price.Low, price.Close, price.Volume, price.MarketCap, price.CurrencyCode)
	return err
}

func FetchBitcoinPrices(db *sql.DB) ([]bitcoin_price.BitcoinPrice, error) {
	rows, err := db.Query("SELECT start, open, high, low, close, volume, market_cap, currency_code FROM bitcoin_prices")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prices []bitcoin_price.BitcoinPrice
	for rows.Next() {
		var price bitcoin_price.BitcoinPrice
		var start string
		if err := rows.Scan(&start, &price.Open, &price.High, &price.Low, &price.Close, &price.Volume, &price.MarketCap, &price.CurrencyCode); err != nil {
			return nil, err
		}
		price.Start, _ = time.Parse(time.RFC3339, start)
		prices = append(prices, price)
	}
	return prices, nil
}
