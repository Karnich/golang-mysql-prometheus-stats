package main

import (
  "fmt"
	"database/sql"
  "os"
  "log"

  "net/http"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"

  _ "github.com/go-sql-driver/mysql"
)

func GetCount(schemadottablename string, column string, db *sql.DB) float64 {
  var cnt float64
  _ = db.QueryRow(`select count(` + column + `) from ` + schemadottablename).Scan(&cnt)
  return cnt 
}

func main() {    
    log.Println("mysql user count started")
   
    // Open up our database connection.
    // I've set up a database on my local machine using phpmyadmin.
		// The database is called testDb
		var (
			username = os.Getenv("DB_USER")
      password = os.Getenv("DB_PASSWORD")
		)
		
    db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/dbname", username, password))
    
    // if there is an error opening the connection, handle it
    if err != nil {
        panic(err.Error())
    }
    
    // defer the close till after the main function has finished
    // executing 
    defer db.Close()

    if err := prometheus.Register(prometheus.NewGaugeFunc(
      prometheus.GaugeOpts{
          Subsystem: "runtime",
          Name:      "user_count",
          Help:      "Number of Users that currently exist.",
      },
      func() float64 { return GetCount("mysql.user", "users", db) },
    )); err == nil {
        log.Println("GaugeFunc 'user count' registered.")
    }
    
    http.Handle("/metrics", promhttp.Handler())
	  log.Fatal(http.ListenAndServe("8080", nil))     
}
