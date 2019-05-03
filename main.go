package main

import (
	"database/sql"
  "os"
  "log"
  "fmt"
  "time"

  "net/http"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"

  _ "github.com/go-sql-driver/mysql"
  _ "github.com/joho/godotenv/autoload"
)

func GetCount(gauge Gauge, db *sql.DB) float64 {
  var cnt float64

  start := time.Now()

  err := db.QueryRow(gauge.Query).Scan(&cnt)

  elapsed := time.Since(start)
  log.Printf("Query %s_%s took %s", gauge.Subsystem, gauge.Name, elapsed)

  if err != nil {
    log.Fatal(err)
  }
  return cnt 
}

func main() {    
    log.Println("mysql user count started")
   
		var (
			username = os.Getenv("DB_USER")
      password = os.Getenv("DB_PASSWORD")
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

    metrics, err := LoadJson();
    if err != nil {
      log.Fatal(fmt.Printf("%+v\n", err))
    }

    for i := 0; i < len(metrics.MetricTypes.Gauges); i++ {
      index := i;
      if err := prometheus.Register(prometheus.NewGaugeFunc(
        prometheus.GaugeOpts{
            Subsystem: metrics.MetricTypes.Gauges[index].Subsystem,
            Name:      metrics.MetricTypes.Gauges[index].Name,
            Help:      metrics.MetricTypes.Gauges[index].Help,
        },
        func() float64 { return GetCount(metrics.MetricTypes.Gauges[index], db) },
      )); 
      err == nil {
        log.Println(fmt.Printf("GaugeFunc %s registered.", metrics.MetricTypes.Gauges[i].Name))
      }
      if err != nil {
        fmt.Printf("%+v\n", err)
      }
    }

    http.Handle("/metrics", promhttp.Handler())
	  log.Fatal(http.ListenAndServe(":8080", nil))     
}

