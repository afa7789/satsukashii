package bigmac

import (
	"sort"
	"time"

	"github.com/afa7789/satsukashii/pkg/bigmac"
	"github.com/afa7789/satsukashii/pkg/bitcoin_price"
	calc "github.com/afa7789/satsukashii/pkg/calculator"
)

type Prices struct {
	Date         int64   `json:"date"`
	Price        float64 `json:"price"`
	PriceSatoshi float64 `json:"price_satoshi"`
}

type PricesData struct {
	Prices          []Prices `json:"prices"`
	MaxPrice        float64  `json:"max_price"`
	MaxPriceSatoshi float64  `json:"max_price_satoshi"`
	LowestDate      int64    `json:"smallest_date"`
	BiggestDate     int64    `json:"biggest_date"`
}

func generatePricesData() (PricesData, error) {
	bmData, err := bigmac.NewBigMacData("assets/csv/big-mac-source-data-v2.csv")
	if err != nil {
		return PricesData{}, err
	}

	var btcData bitcoin_price.BitcoinPriceFetcher
	btcData, err = bitcoin_price.NewBTCPricesCSV("assets/csv/bitcoin_2010-07-17_2024-12-05.csv")
	if err != nil {
		return PricesData{}, err
	}

	date := "2010-07-10"
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return PricesData{}, err
	}
	historicalData, err := btcData.FetchHistoricalData(parsedDate)
	if err != nil {
		return PricesData{}, err
	}

	// Sort historical data by date
	var dates []time.Time
	for d := range historicalData {
		dates = append(dates, d)
	}

	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	var results []Prices
	maxPrice := 0.0
	maxPriceSatoshi := 0.0

	// Iterate over sorted dates and fetch the corresponding BitcoinPrice
	for _, d := range dates {
		data := historicalData[d]
		bigmacPrice, _, ok := bmData.GetPriceTimestamp("USD", data.Start)
		if ok {
			bitcoinAmount := calc.CalculateBitcoinAmount(data.Close, bigmacPrice)
			satsValue := calc.BTCtoSATS(bitcoinAmount)
			if bigmacPrice > maxPrice {
				maxPrice = bigmacPrice
			}
			if satsValue > maxPriceSatoshi {
				maxPriceSatoshi = satsValue
			}
			results = append(results, Prices{
				Date:         data.Start.Unix(),
				Price:        bigmacPrice,
				PriceSatoshi: satsValue,
			})
		}
	}

	var lowestDateUnix, biggestDateUnix int64
	if len(dates) > 0 {
		lowestDateUnix = dates[0].Unix()
		biggestDateUnix = dates[len(dates)-1].Unix()
	}

	return PricesData{
		Prices:          results,
		MaxPrice:        maxPrice,
		MaxPriceSatoshi: maxPriceSatoshi,
		LowestDate:      lowestDateUnix,
		BiggestDate:     biggestDateUnix,
	}, nil
}

type ChartData struct {
	SizeHeight     int64
	SizeWidth      int64
	X1Array        []float64
	Y1Array        []float64
	Y1ArraySatoshi []float64
	// X2Array        []float64
	// Y2Array        []float64
	// Y2ArraySatoshi []float64
}

func GenerateChartData(h, w int64) (ChartData, error) {
	sizeHeight := h
	sizeWidth := w

	// Get the Big Mac data
	bmData, err := generatePricesData()
	if err != nil {
		return ChartData{}, err
	}

	// array of lineX!
	// Price
	X1Array := make([]float64, 0) // dot and line
	Y1Array := make([]float64, 0)
	// lineX2Array := make([]float64, 0)
	// lineY2Array := make([]float64, 0)

	// Satoshi Lines
	YS1Array := make([]float64, 0)
	// lineYS2Array := make([]float64, 0)

	// fix max value at PriceData for Price and PriceSatoshi ?

	for _, data := range bmData.Prices {
		// fill x1, y1 and ys1
		Y1Array = append(
			Y1Array,
			data.Price*float64(sizeHeight)/bmData.MaxPrice,
		)
		X1Array = append(
			X1Array,
			float64(data.Date-bmData.LowestDate)*float64(sizeWidth)/float64(bmData.BiggestDate-bmData.LowestDate),
		)
		YS1Array = append(
			YS1Array,
			data.PriceSatoshi*float64(sizeHeight)/bmData.MaxPriceSatoshi,
		)
	}

	// // Build the line arrays by shifting the X1Array, Y1Array and YS1Array by one element.
	// if len(X1Array) > 1 {
	// 	lineX2Array = append([]float64{}, X1Array[1:]...)
	// 	lineY2Array = append([]float64{}, Y1Array[1:]...)
	// 	lineYS2Array = append([]float64{}, YS1Array[1:]...)
	// }

	return ChartData{
		SizeHeight:     sizeHeight,
		SizeWidth:      sizeWidth,
		X1Array:        X1Array,
		Y1Array:        Y1Array,
		Y1ArraySatoshi: YS1Array,
		// X2Array:        lineX2Array,
		// Y2Array:        lineY2Array,
		// Y2ArraySatoshi: lineYS2Array,
	}, nil
}
