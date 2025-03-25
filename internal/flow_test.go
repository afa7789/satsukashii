package flow_test

import (
	"database/sql"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/afa7789/satsukashii/pkg/bigmac"
	"github.com/afa7789/satsukashii/pkg/bitcoin_price"
	calc "github.com/afa7789/satsukashii/pkg/calculator"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// TestCriticalFlow tests the critical flow of integrating Big Mac pricing data
// and Bitcoin historical price data.
//
// It performs the following steps:
// 1. Loads Big Mac pricing data from a CSV file.
// 2. Creates a temporary SQLite database, populates it with sample data, and connects to it.
// 3. Fetches historical Bitcoin price data for a specific date ("2010-07-10").
// 4. Sorts the historical data by date.
// 5. For each date, retrieves the Big Mac price and calculates the equivalent Bitcoin amount.
// 6. Tests CSV-based Bitcoin price fetching and calculations for a different date ("2020-01-14").
// 7. Cleans up by deleting the temporary SQLite file.
func TestCriticalFlow(t *testing.T) {
	// Step 1: Load Big Mac data
	bmData, err := bigmac.NewBigMacData("assets/csv/big-mac-source-data-v2.csv")
	if err != nil {
		t.Fatalf("Failed to load Big Mac data: %v", err)
	}

	// Step 2: Create a temporary SQLite database
	tmpFile, err := os.CreateTemp("", "satsukashii_test_*.db")
	if err != nil {
		t.Fatalf("Failed to create temporary SQLite file: %v", err)
	}
	tmpDBPath := tmpFile.Name()
	tmpFile.Close()            // Close the file so SQLite can open it
	defer os.Remove(tmpDBPath) // Ensure the file is deleted after the test

	// Connect to the temporary database
	db, err := sql.Open("sqlite3", tmpDBPath)
	if err != nil {
		t.Fatalf("Failed to connect to temporary SQLite database: %v", err)
	}
	defer db.Close()

	// Create a sample table and insert test data
	_, err = db.Exec(`
		CREATE TABLE bitcoin_prices (
			start_date TEXT PRIMARY KEY,
			close REAL
		);
		INSERT INTO bitcoin_prices (start_date, close) VALUES
			('2010-07-10', 0.05),
			('2010-07-11', 0.06),
			('2010-07-12', 0.07);
	`)
	if err != nil {
		t.Fatalf("Failed to initialize temporary SQLite database: %v", err)
	}

	// Initialize BitcoinPriceFetcher with the temporary database
	btcData := bitcoin_price.NewBitcoinPriceDB(db)

	// Step 3: Parse test date and fetch historical Bitcoin data
	date := "2010-07-10"
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		t.Fatalf("Failed to parse date %q: %v", date, err)
	}

	historicalData, err := btcData.FetchHistoricalData(parsedDate)
	if err != nil {
		t.Fatalf("Failed to fetch historical Bitcoin data for %v: %v", parsedDate, err)
	}

	// Step 4: Sort historical data by date
	var dates []time.Time
	for date := range historicalData {
		dates = append(dates, date)
	}
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	// Step 5: Verify sorted dates and process data
	if len(dates) == 0 {
		t.Fatal("No historical data dates found")
	}
	t.Logf("Number of historical data dates: %d", len(dates))

	// Process each date and calculate Bitcoin amounts
	for _, date := range dates {
		data := historicalData[date]
		bigmacPrice, _, ok := bmData.GetPriceTimestamp("USD", data.Start)
		if !ok {
			t.Logf("No Big Mac price found for date %v", data.Start)
			continue
		}

		calculateBitcoinPriceBigmac := calc.CalculateBitcoinAmount(data.Close, bigmacPrice)
		sats := calc.BTCtoSATS(calculateBitcoinPriceBigmac)
		t.Logf("Date: %s, Big Mac Price: %f, SATS: %f", data.Start.Format("2006-01-02"), bigmacPrice, sats)
	}

	// Step 6: Additional CSV-based Bitcoin price test
	dateCSV := "2020-01-14"
	bigmacPrice, ok := bmData.GetPrice("USD", dateCSV)
	if !ok {
		t.Fatalf("Failed to fetch Big Mac price for USD on %s", dateCSV)
	}
	t.Logf("Price for USD on %s: %f", dateCSV, bigmacPrice)

	bpcsv, err := bitcoin_price.NewBTCPricesCSV("assets/csv/bitcoin_2010-07-17_2025-04-25.csv")
	if err != nil {
		t.Fatalf("Failed to load Bitcoin price CSV: %v", err)
	}

	parsedDateCSV, err := time.Parse("2006-01-02", dateCSV)
	if err != nil {
		t.Fatalf("Failed to parse date %q: %v", dateCSV, err)
	}

	btcPrice, err := bpcsv.FetchPriceByDate(parsedDateCSV)
	if err != nil {
		t.Fatalf("Failed to fetch Bitcoin price for %v: %v", parsedDateCSV, err)
	}
	t.Logf("Price in BTC on %s: %f", dateCSV, btcPrice.Close)

	bigmacInBTC := calc.CalculateBitcoinAmount(btcPrice.Close, bigmacPrice)
	sats := calc.BTCtoSATS(bigmacInBTC)
	t.Logf("Price in BTC on %s: %f, SATS: %f", dateCSV, bigmacInBTC, sats)

	// Step 7: Cleanup is handled by defer os.Remove(tmpDBPath)
}
