package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func insert(db *sql.DB, sr *StatsResponse) error {
	err := insertPurchased(db, sr)
	if err != nil {
		return err
	}

	err = insertActiveUsers(db, sr)
	if err != nil {
		return err
	}

	err = insertTotalUsers(db, sr)
	if err != nil {
		return err
	}

	err = insertUsed(db, sr)
	if err != nil {
		return err
	}

	return nil
}

func createTables(db *sql.DB) error {
	err := createPurchasedTable(db)
	if err != nil {
		return err
	}

	err = createActiveUsersTable(db)
	if err != nil {
		return err
	}

	err = createTotalUsersTable(db)
	if err != nil {
		return err
	}

	err = createUsedTable(db)
	if err != nil {
		return err
	}

	err = createBalanceTable(db)
	if err != nil {
		return err
	}

	err = createTotalFilesTable(db)
	if err != nil {
		return err
	}

	return nil
}

func updateBalance(db *sql.DB, api string) {
	client := http.DefaultClient
	res, err := client.Get(api)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer res.Body.Close()

	var sr BalanceResponse
	err = json.NewDecoder(res.Body).Decode(&sr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bals := sr.Balances
	for _, bal := range bals {
		if bal.Denom == "ujkl" {
			err = insertBalance(db, bal.Amount)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}
}

func updateFiles(db *sql.DB, api string) {
	client := http.DefaultClient
	res, err := client.Get(api)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer res.Body.Close()

	var sr PageResponse
	err = json.NewDecoder(res.Body).Decode(&sr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = insertTotalFiles(db, sr.Pagination.Total)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func update(db *sql.DB, api string) {
	client := http.DefaultClient
	res, err := client.Get(api)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer res.Body.Close()

	var sr StatsResponse
	err = json.NewDecoder(res.Body).Decode(&sr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = insert(db, &sr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found.")
	}

	host := os.Getenv("STATS_HOST")
	portString := os.Getenv("STATS_PORT")
	user := os.Getenv("STATS_USER")
	password := os.Getenv("STATS_PASSWORD")
	dbname := os.Getenv("STATS_DB_NAME")

	api := os.Getenv("STATS_API")
	balanceAPI := os.Getenv("BALANCE_API")
	filesAPI := os.Getenv("FILES_API")

	port, err := strconv.ParseInt(portString, 10, 64)
	if err != nil {
		panic("cannot parse port")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = createTables(db)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	go func() {
		err := StartAPI(db)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}()

	for {
		fmt.Println(time.Now().String())
		update(db, api)
		updateBalance(db, balanceAPI)
		updateFiles(db, filesAPI)
		time.Sleep(time.Minute * 60)
	}
}
