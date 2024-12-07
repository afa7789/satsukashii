package main

import (
	"database/sql"
	"log"
	"sort"
	"time"

	"github.com/afa7789/satsukashii/pkg/bigmac"
	"github.com/afa7789/satsukashii/pkg/bitcoin_price"
	calc "github.com/afa7789/satsukashii/pkg/calculator"
)

func main() {
	testerfunction()
	bmData, err := bigmac.NewBigMacData("assets/csv/big-mac-source-data-v2.csv")
	if err != nil {
		panic(err)
	}

	var btcData bitcoin_price.BitcoinPriceFetcher

	// Connect to the SQLite database
	db, err := sql.Open("sqlite3", "assets/database/satsukashii.db")
	if err != nil {
		log.Fatal(err)
		return
	}
	btcData = bitcoin_price.NewBitcoinPriceDB(db)

	date := "2010-07-10"
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(err)
	}
	historicalData, err := btcData.FetchHistoricalData(parsedDate)
	if err != nil {
		panic(err)
	}
	// sort historical data by date
	// Step 1: Extract keys (dates) from the map
	var dates []time.Time
	for date := range historicalData {
		dates = append(dates, date)
	}

	// Step 2: Sort the dates slice
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j]) // Compare dates
	})

	// Step 3: Iterate over sorted dates and fetch the corresponding BitcoinPrice
	for _, date := range dates {
		data := historicalData[date]
		// Simulating price fetching (as `bmData` and `calc` are not defined)
		bigmacPrice, _, ok := bmData.GetPriceTimestamp("USD", data.Start)
		calculateBitcoinPriceBigmac := calc.CalculateBitcoinAmount(data.Close, bigmacPrice)
		if ok {
			log.SetFlags(0)
			log.Printf("Date: %s, Price: %f, SATS: %f", data.Start.Format("2006-01-02"), bigmacPrice, calc.BTCtoSATS(calculateBitcoinPriceBigmac))
		}
	}

	// for _, data := range historicalData {
	// 	// log.Printf("Date: %s, Price: %f", data.Start, data.Close)
	// 	bigmacPrice, _, ok := bmData.GetPriceTimestamp("USD", data.Start)
	// 	if ok {
	// 		log.SetFlags(0)
	// 		log.Printf("Date: %s, Price: %f, SATS: %f", data.Start.Format("2006-01-02"), bigmacPrice, calc.BTCtoSATS(data.Close))
	// 	}
	// }
}

func testerfunction() {
	// big mac data from file
	bmd, err := bigmac.NewBigMacData("assets/csv/big-mac-source-data-v2.csv")
	if err != nil {
		panic(err)
	}
	date := "2020-01-14"
	bigmac_price, ok := bmd.GetPrice("USD", date)
	if ok {
		log.Printf("Price for USD on %s: %f", date, bigmac_price)
		bpcsv, err := bitcoin_price.NewBTCPricesCSV("assets/csv/bitcoin_2010-07-17_2024-12-05.csv")
		if err != nil {
			panic(err)
		}
		parsedDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			panic(err)
		}
		btc_price, err := bpcsv.FetchPriceByDate(parsedDate)
		if err != nil {
			panic(err)
		}
		log.Printf("Price in BTC on %s: %f", date, btc_price.Close)
		bigmac_in_btc := calc.CalculateBitcoinAmount(btc_price.Close, bigmac_price)
		log.Printf(
			"Price in BTC on %s: %f, sats: %f", date,
			bigmac_in_btc,
			calc.BTCtoSATS(bigmac_in_btc),
		)

	}
}
