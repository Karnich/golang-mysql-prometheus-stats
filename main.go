package main

import (
	"database/sql"
  "os"
  "log"
  "fmt"

  "net/http"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"

  _ "github.com/go-sql-driver/mysql"
  _ "github.com/joho/godotenv/autoload"
)

func GetCount(query string, db *sql.DB) float64 {
  var cnt float64

  err := db.QueryRow(query).Scan(&cnt)

  if err != nil {
    log.Fatal(err)
  }
  return cnt 
}

func main() {    
    log.Println("mysql user count started")
   
		var (
      subsystem = os.Getenv("SUBSYSTEM")
      name = os.Getenv("NAME")
      help = os.Getenv("HELP")
			username = os.Getenv("DB_USER")
      password = os.Getenv("DB_PASSWORD")
      query = os.Getenv("QUERY")
      host = os.Getenv("DB_HOST")
      port = os.Getenv("DB_PORT")
      database = os.Getenv("DB_NAME")
		)
		
    db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database))
    
    // if there is an error opening the connection, handle it
    if err != nil {
        panic(err.Error())
    }
    
    // defer the close till after the main function has finished
    // executing 
    defer db.Close()

    if err := prometheus.Register(prometheus.NewGaugeFunc(
      prometheus.GaugeOpts{
          Subsystem: subsystem,
          Name:      name,
          Help:      help,
      },
      func() float64 { return GetCount(query, db) },
    )); err == nil {
        log.Println("GaugeFunc 'user count' registered.")
    }
    
    http.Handle("/metrics", promhttp.Handler())
	  log.Fatal(http.ListenAndServe(":8080", nil))     
}
