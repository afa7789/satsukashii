// DISCLAIMER
// this is just a placeholder, as I was brainstorming making it a service that would be kept in the air for a long time
// I don't care enough to make it work, so I'm just going to leave it here as a placeholder
// I would need an API for it.

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/afa7789/satsukashii/internal/database"
	"github.com/afa7789/satsukashii/pkg/bitcoin_price"
)

func main() {
	// start the DB
	db := database.CreateDB("assets/database/satsukashii.db")
	if db != nil {
		defer db.Close()
	} else {
		panic("Error creating DB")
	}

	// start the ingester
	// fill db with data
	var btcData bitcoin_price.BitcoinPriceFetcher
	btcData, err := bitcoin_price.NewBTCPricesCSV("assets/csv/bitcoin_2010-07-17_2024-12-05.csv")
	if err != nil {
		panic(err)
	}

	firstTime, err := time.Parse("2006-01-02", "2010-07-10")
	if err != nil {
		panic(err)
	}
	mapAllData, err := btcData.FetchHistoricalData(firstTime)
	if err != nil {
		panic(err)
	}

	log.Printf("Fetched %d data points\n", len(mapAllData))

	//convert map to array
	allData := make([]bitcoin_price.BitcoinPrice, 0, len(mapAllData))
	for _, v := range mapAllData {
		allData = append(allData, v)
	}

	// insert data into the DB in batches of 1000
	batchSize := 1000
	for i := 0; i < len(allData); i += batchSize {
		end := i + batchSize
		if end > len(allData) {
			end = len(allData)
		}
		err = database.InsertBitcoinPricesBatch(db, allData[i:end])
		if err != nil {
			panic(err)
		}
		fmt.Printf("Inserted %d data points, batch: %d\n", end-i, i/batchSize+1)
	}
	// starter or ingester, name may change, will run the first tasks needed to fill the DB.
	// panic("not implemented")
}
