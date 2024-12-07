package bitcoin_price

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type BitcoinPriceDB struct {
	DB *sql.DB
}

func NewBitcoinPriceDB(db *sql.DB) *BitcoinPriceDB {
	return &BitcoinPriceDB{DB: db}
}

func (bpdb *BitcoinPriceDB) FetchHistoricalData(startDate time.Time) (map[time.Time]BitcoinPrice, error) {
	db := bpdb.DB
	rows, err := db.Query(`
        SELECT start, open, high, low, close, volume, market_cap, currency_code 
        FROM bitcoin_prices 
        WHERE start >= ?
    `, startDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	historicalData := make(map[time.Time]BitcoinPrice)
	for rows.Next() {
		var price BitcoinPrice
		var start string
		if err := rows.Scan(&start, &price.Open, &price.High, &price.Low, &price.Close, &price.Volume, &price.MarketCap, &price.CurrencyCode); err != nil {
			return nil, err
		}
		parsedTime, _ := time.Parse(time.RFC3339, start)
		price.Start = parsedTime
		historicalData[price.Start] = price
	}

	return historicalData, nil
}

func (bpdb *BitcoinPriceDB) FetchPriceByDate(date time.Time) (BitcoinPrice, error) {
	db := bpdb.DB
	query := `
        SELECT start, open, high, low, close, volume, market_cap, currency_code 
        FROM bitcoin_prices 
        WHERE start = ?
    `

	var price BitcoinPrice
	var start string
	err := db.QueryRow(query, date.Format(time.RFC3339)).Scan(&start, &price.Open, &price.High, &price.Low, &price.Close, &price.Volume, &price.MarketCap, &price.CurrencyCode)
	if err != nil {
		return BitcoinPrice{}, err
	}

	price.Start, _ = time.Parse(time.RFC3339, start)
	return price, nil
}
