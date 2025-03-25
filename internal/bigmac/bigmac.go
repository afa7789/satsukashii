package bigmac

import (
	"math"
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

func generatePricesData(btcData bitcoin_price.BitcoinPriceFetcher) (PricesData, error) {
	bmData, err := bigmac.NewBigMacData("assets/csv/big-mac-source-data-v2.csv")
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
	SpaceDiff      int64
	SizeHeight     int64
	SizeWidth      int64
	X1Array        []float64
	Y1Array        []float64
	Y1ArraySatoshi []float64
	// X2Array        []float64
	// Y2Array        []float64
	// Y2ArraySatoshi []float64
	MaxPrice        float64
	MaxPriceSatoshi float64
}

func GenerateChartData(btcData bitcoin_price.BitcoinPriceFetcher, h, w, spaceDiff int64) (ChartData, error) {
	sizeHeight := h
	sizeWidth := w

	bmData, err := generatePricesData(btcData)
	if err != nil {
		return ChartData{}, err
	}

	X1Array := make([]float64, 0)
	Y1Array := make([]float64, 0)
	YS1Array := make([]float64, 0)

	for _, data := range bmData.Prices {
		X1Array = append(
			X1Array,
			float64(spaceDiff)+(float64(data.Date-bmData.LowestDate)*float64(sizeWidth-(2*spaceDiff))/float64(bmData.BiggestDate-bmData.LowestDate)),
		)
		Y1Array = append(
			Y1Array,
			(float64(sizeHeight))-(data.Price*float64(sizeHeight-spaceDiff)/bmData.MaxPrice)+float64(spaceDiff),
		)

		satoshiPriceCalculate := (math.Log(data.PriceSatoshi) / math.Log(bmData.MaxPriceSatoshi)) * float64(sizeHeight-spaceDiff)
		satoshiEntry := (float64(sizeHeight)) - satoshiPriceCalculate + float64(spaceDiff)
		YS1Array = append(
			YS1Array,
			satoshiEntry,
		)
	}

	return ChartData{
		SpaceDiff:       spaceDiff,
		SizeHeight:      sizeHeight,
		SizeWidth:       sizeWidth,
		X1Array:         X1Array,
		Y1Array:         Y1Array,
		Y1ArraySatoshi:  YS1Array,
		MaxPrice:        bmData.MaxPrice,
		MaxPriceSatoshi: bmData.MaxPriceSatoshi,
	}, nil
}
