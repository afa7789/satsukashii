package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDB(filePath string) {
	// Connect to the SQLite database
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// empty table and drop it if it exists
	droptable := `DROP TABLE IF EXISTS bitcoin_prices;`
	_, err = db.Exec(droptable)
	if err != nil {
		log.Println("Error dropping table:", err)
		return
	}

	// Create the table if it doesn't exist
	createTable := `CREATE TABLE IF NOT EXISTS bitcoin_prices (
	
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        start DATETIME NOT NULL,
        open REAL NOT NULL,
        high REAL NOT NULL,
        low REAL NOT NULL,
        close REAL NOT NULL,
        volume REAL NOT NULL,
        market_cap REAL NOT NULL,
        currency_code TEXT NOT NULL
    );`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Println("Error on creating table(s):", err)
		return
	}

	log.Println("Database and table set up successfully!")
}
