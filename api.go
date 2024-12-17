package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Name string          `json:"name"`
	Data []ResponseEntry `json:"data"`
}

type ResponseEntry struct {
	Time  string  `json:"time"`
	Value float64 `json:"value"`
}

func (r *Response) sort() {
	sort.Slice(r.Data, func(i, j int) bool {
		t1, err1 := time.Parse("2006-01-02", r.Data[i].Time)
		t2, err2 := time.Parse("2006-01-02", r.Data[j].Time)
		if err1 != nil || err2 != nil {
			return false
		}
		return t1.Before(t2)
	})
}

func (r *Response) add(date time.Time, amount int) {
	formatted := date.Format("2006-01-02")
	re := ResponseEntry{
		Time:  formatted,
		Value: float64(amount),
	}

	r.Data = append(r.Data, re)
}

func q(db *sql.DB, table string, start, end int) (*Response, error) {
	endDate := time.Now().AddDate(0, 0, 1-end).Format("2006-01-02")
	startDate := time.Now().AddDate(0, 0, -start).Format("2006-01-02")

	query := fmt.Sprintf(`
        SELECT a.date, a.amount
        FROM %s a
        JOIN (
            SELECT MIN(date) as min_date
            FROM %s
            WHERE DATE(date) BETWEEN $1 AND $2
            GROUP BY DATE(date)
        ) b ON a.date = b.min_date
    `, table, table)

	rows, err := db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := &Response{
		Name: table,
		Data: make([]ResponseEntry, 0),
	}

	for rows.Next() {
		var date time.Time
		var amount int
		if err := rows.Scan(&date, &amount); err != nil {
			return nil, err
		}
		res.add(date, amount)
		fmt.Printf("Date: %s | Amount: %d\n", date.String(), amount)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	res.sort()

	return res, nil
}

func StartAPI(db *sql.DB) error {
	paths := []string{"active_users", "purchased", "total_users", "used", "protocol_balance", "total_files", "available_space"}

	r := gin.Default()

	for _, path := range paths {
		r.GET(fmt.Sprintf("/%s", path), func(c *gin.Context) {
			startString := c.DefaultQuery("start", "30")
			endString := c.DefaultQuery("end", "0")

			start, err := strconv.ParseInt(startString, 10, 64)
			if err != nil {
				err := c.Error(err)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
			end, err := strconv.ParseInt(endString, 10, 64)
			if err != nil {
				err := c.Error(err)
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			res, err := q(db, path, int(start), int(end))
			if err != nil {
				err := c.Error(err)
				if err != nil {
					fmt.Println(err)
				}
				return
			}

			c.JSON(http.StatusOK, res)
		})
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, paths)
	})

	err := r.Run("0.0.0.0:5756")
	if err != nil {
		return err
	}

	return nil
}
