package main

import (
	"TradingBot/cmd"
	"TradingBot/cmd/rest-server/handler"
	"TradingBot/config"
	"TradingBot/internal/marketdata"
	"TradingBot/models"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

func main() {
	os.Setenv("API_TOKEN", "hello")
	os.Setenv("APP_ENV", "stage")
	// Connect to the database
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	handlers := handler.NewHandler(db)

	// Auto-migrate the models
	err = db.AutoMigrate(&models.Candle{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Your Alpha Vantage API key
	apiKey := "your_alpha_vantage_api_key_here"

	// Define the pairs you want to fetch data for
	pairs := []string{"EURUSDT"}

	// Loop through the pairs and fetch the candles
	for _, pair := range pairs {
		candles, err := marketdata.FetchCandlesLast31Days(pair, apiKey)
		if err != nil {
			log.Printf("Error fetching candles for pair %s: %v", pair, err)
			continue
		}

		// Save the fetched candles to the database
		marketdata.SaveCandles(db, pair, candles)

	}

	srv := new(cmd.Server)

	if err := srv.Run(":8071", handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
}
