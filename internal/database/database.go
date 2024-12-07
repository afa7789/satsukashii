package database

import (
	"database/sql"
	"log"
	"strings"
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

func InsertBitcoinPricesBatch(db *sql.DB, prices []bitcoin_price.BitcoinPrice) error {
	if len(prices) == 0 {
		return nil // Nothing to insert
	}

	// Start the transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Build the query with placeholders for batch insert
	query := `
    INSERT INTO bitcoin_prices (start, open, high, low, close, volume, market_cap, currency_code)
    VALUES `

	// Slice to store query placeholders and values
	var placeholders []string
	var values []interface{}

	// Build placeholders and values dynamically
	for _, price := range prices {
		placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?, ?, ?)")
		values = append(values, price.Start, price.Open, price.High, price.Low, price.Close, price.Volume, price.MarketCap, price.CurrencyCode)
	}

	// Join placeholders to complete the query
	query += strings.Join(placeholders, ", ")

	// Execute the batch insert
	_, err = tx.Exec(query, values...)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Print("rollback failed")
		}
		return err
	}

	// Commit the transaction
	return tx.Commit()
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
