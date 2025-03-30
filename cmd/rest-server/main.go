package main

import (
	"fmt"
	"github.com/Olovets/TradingBot/cmd"
	"github.com/Olovets/TradingBot/cmd/rest-server/handler"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	os.Setenv("API_TOKEN", "hello")
	os.Setenv("APP_ENV", "stage")

	handlers := handler.NewHandler()

	// Your Alpha  API key

	// Define the pairs you want to fetch data for
	//pairs := []string{"EURUSDT"}

	//// Loop through the pairs and fetch the candles
	//for _, pair := range pairs {
	//	candles, err := marketdata.FetchCandlesLast31Days(pair, apiKey)
	//	if err != nil {
	//		log.Printf("Error fetching candles for pair %s: %v", pair, err)
	//		continue
	//	}
	//
	//	// Save the fetched candles to the database
	//	marketdata.SaveCandles(db, pair, candles)
	//
	//}

	srv := new(cmd.Server)

	if err := srv.Run(":8071", handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
	fmt.Println("Server started")
}
