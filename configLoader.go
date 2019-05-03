package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type MetricTypes struct {
	MetricTypes MetricType `json:"metrictypes"`
}

type MetricType struct {
	Gauges []Gauge `json:"gauges"`
}

type Gauge struct {
	Name   string `json:"name"`
	Subsystem   string `json:"subsystem"`
	Help    string    `json:"help"`
	Query string `json:"query"`
}

func LoadJson() (MetricTypes, error) {
	jsonFile, err := os.Open("/config/test.json")
	// if we os.Open returns an error then handle it
	if err != nil {
			fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var metrics MetricTypes

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &metrics)

	// we iterate through every user within our users array and
	// print out the user Type, their name, and their facebook url
	// as just an example
	for i := 0; i < len(metrics.MetricTypes.Gauges); i++ {
			fmt.Println("Name: " + metrics.MetricTypes.Gauges[i].Name)
			fmt.Println("Subsystem: " + metrics.MetricTypes.Gauges[i].Subsystem)
			fmt.Println("Help: " + metrics.MetricTypes.Gauges[i].Help)
			fmt.Println("Query: " + metrics.MetricTypes.Gauges[i].Query)
	}

	return metrics, nil
}