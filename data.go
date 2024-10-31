package main

import "database/sql"

func createPurchasedTable(db *sql.DB) error {
	createTable := `
	CREATE TABLE IF NOT EXISTS purchased (
		date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		amount BIGINT NOT NULL,
		PRIMARY KEY (date)
	);`

	_, err := db.Exec(createTable)
	return err
}

func insertPurchased(db *sql.DB, response *StatsResponse) error {
	_, err := db.Exec(`INSERT INTO purchased (amount) VALUES ($1)`, response.Purchased)
	return err
}

func createUsedTable(db *sql.DB) error {
	createTable := `
	CREATE TABLE IF NOT EXISTS used (
		date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		amount BIGINT NOT NULL,
		PRIMARY KEY (date)
	);`

	_, err := db.Exec(createTable)
	return err
}

func insertUsed(db *sql.DB, response *StatsResponse) error {
	_, err := db.Exec(`INSERT INTO used (amount) VALUES ($1)`, response.Used)
	return err
}

func createActiveUsersTable(db *sql.DB) error {
	createTable := `
	CREATE TABLE IF NOT EXISTS active_users (
		date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		amount BIGINT NOT NULL,
		PRIMARY KEY (date)
	);`

	_, err := db.Exec(createTable)
	return err
}

func insertActiveUsers(db *sql.DB, response *StatsResponse) error {
	_, err := db.Exec(`INSERT INTO active_users (amount) VALUES ($1)`, response.ActiveUsers)
	return err
}

func createTotalUsersTable(db *sql.DB) error {
	createTable := `
	CREATE TABLE IF NOT EXISTS total_users (
		date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		amount BIGINT NOT NULL,
		PRIMARY KEY (date)
	);`

	_, err := db.Exec(createTable)
	return err
}

func insertTotalUsers(db *sql.DB, response *StatsResponse) error {
	_, err := db.Exec(`INSERT INTO total_users (amount) VALUES ($1)`, response.TotalUsers)
	return err
}

func createBalanceTable(db *sql.DB) error {
	createTable := `
	CREATE TABLE IF NOT EXISTS protocol_balance (
		date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		amount BIGINT NOT NULL,
		PRIMARY KEY (date)
	);`

	_, err := db.Exec(createTable)
	return err
}

func insertBalance(db *sql.DB, amount uint64) error {
	_, err := db.Exec(`INSERT INTO protocol_balance (amount) VALUES ($1)`, amount)
	return err
}

func createTotalFilesTable(db *sql.DB) error {
	createTable := `
	CREATE TABLE IF NOT EXISTS total_files (
		date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		amount BIGINT NOT NULL,
		PRIMARY KEY (date)
	);`

	_, err := db.Exec(createTable)
	return err
}

func insertTotalFiles(db *sql.DB, amount uint64) error {
	_, err := db.Exec(`INSERT INTO total_files (amount) VALUES ($1)`, amount)
	return err
}

func createAvailableSpace(db *sql.DB) error {
	createTable := `
	CREATE TABLE IF NOT EXISTS available_space (
		date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		amount BIGINT NOT NULL,
		PRIMARY KEY (date)
	);`

	_, err := db.Exec(createTable)
	return err
}

func insertAvailableSpace(db *sql.DB, amount uint64) error {
	_, err := db.Exec(`INSERT INTO available_space (amount) VALUES ($1)`, amount)
	return err
}
