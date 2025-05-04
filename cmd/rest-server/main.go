package main

import (
	"fmt"
	"log"

	"github.com/Olovets/TradingBot/config"
	"github.com/Olovets/TradingBot/internal/entity/blocks"
	"github.com/Olovets/TradingBot/internal/marketdata"
	"github.com/Olovets/TradingBot/internal/models"
	"gorm.io/gorm"
)

func main() {
	var db *gorm.DB
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	fmt.Println("Successfully connected to database")

	db.AutoMigrate(&models.Candle{})

	candles, _ := marketdata.FetchCandlesLast31Days("EURUSDT", "asadasdasdas")

	marketdata.SaveCandles(db, "EURUSDT", candles)

	allCandles, err := models.GetAllCandles(db)
	if err != nil {
		log.Fatalf("Failed to get candles: %v", err)
	}
	allCandles = models.CutNonNeededCandles(allCandles)

	actualPrice := allCandles[len(allCandles)-1].Close

	candles, err = models.GenerateCandles(allCandles, "1d", int64(0), int64(0))
	if err != nil {
		log.Fatalf("Failed to generate candles: %v", err)
	}

	dailyOb := blocks.NewOrderBlock().IdentifyAll(candles)
	dailyFvg := blocks.NewFvgBlock().IdentifyAll(candles)
	rbFvg := blocks.NewRejectionBlock().IdentifyAll(candles)

	dailyLastBlock := models.ReturnLast(dailyOb)

	lastDailyFvg := models.ReturnLast(dailyFvg)
	if lastDailyFvg.Timestamp > dailyLastBlock.Timestamp {
		dailyLastBlock = lastDailyFvg
	}

	lastRbFvg := models.ReturnLast(rbFvg)
	if lastRbFvg.Timestamp > dailyLastBlock.Timestamp {
		dailyLastBlock = lastRbFvg
	}
	allDailyBlocks := models.AggregateBlocks([][]models.Block{dailyOb, dailyFvg, rbFvg})

	dailyLastConterTrend := models.ReturnNearFTA(dailyLastBlock, allDailyBlocks)

	fmt.Println(dailyLastConterTrend, actualPrice)

	fmt.Println("Server started")
}
