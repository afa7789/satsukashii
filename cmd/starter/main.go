package main

import "github.com/afa7789/satsukashii/internal/database"

func main() {
	// start the DB
	database.CreateDB("assets/database/satsukashii.db")
	// starter or ingester, name may change, will run the first tasks needed to fill the DB.
	// panic("not implemented")
}
