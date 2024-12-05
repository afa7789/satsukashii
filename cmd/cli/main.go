package main

import (
	"log"

	"github.com/afa7789/satsukashii/pkg/bigmac"
)

func main() {
	// big mac data from file
	bmd, err := bigmac.NewBigMacData("resources/csv/big-mac-source-data-v2.csv")
	if err != nil {
		panic(err)
	}
	price, ok := bmd.GetPrice("USD", "2020-01-01")
	if ok {
		log.Printf("Price for USD on 2020-01-01: %f", price)
	}
}
