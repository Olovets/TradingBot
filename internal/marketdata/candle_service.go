// internal/marketdata/candle_service.go
package marketdata

import (
	"TradingBot/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
)

// FetchAllCandlesByPair retrieves all candles for a given trading pair from the database,
// adjusts the timestamp by adding 7 hours, and sorts them from new to old.
func FetchAllCandlesByPair(db *gorm.DB, pair string) ([]models.Candle, error) {
	var candles []models.Candle

	// Query the database for all candles matching the trading pair and sort by timestamp (descending)
	err := db.Where("pair = ?", pair).Order("timestamp DESC").Find(&candles).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching candles from database: %v", err)
	}

	return candles, nil
}

// SaveCandles will save new candles for different pairs to the database if they don't already exist
func SaveCandles(db *gorm.DB, pair string, candles []models.Candle) {
	for _, candle := range candles {
		// Check if the candle already exists based on the Pair and Timestamp
		var existingCandle models.Candle
		err := db.Where("pair = ? AND timestamp = ?", pair, candle.Timestamp).First(&existingCandle).Error

		if err != nil && err != gorm.ErrRecordNotFound {
			log.Printf("Error checking for existing candle: %v", err)
			continue
		}

		// If no existing candle found, save the new one
		if err == gorm.ErrRecordNotFound {
			result := db.Create(&candle)
			if result.Error != nil {
				log.Printf("Failed to save candle: %v", result.Error)
				continue
			}

			log.Printf("Candle for pair %s and timestamp %s saved successfully!", pair, candle.Timestamp)
		} else {
			log.Printf("Candle for pair %s and timestamp %s already exists, skipping.", pair, candle.Timestamp)
		}
	}
}

// GetLatestCandleTimestamp retrieves the latest timestamp for a given trading pair from the database
func GetLatestCandleTimestamp(db *gorm.DB, pair string) (int64, error) {
	var latestCandle models.Candle

	// Query the database for the latest candle by pair, ordered by timestamp in descending order
	err := db.Where("pair = ?", pair).Order("timestamp DESC").First(&latestCandle).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No records found for the given pair
			return 0, nil
		}
		return 0, fmt.Errorf("error fetching latest candle timestamp: %v", err)
	}

	return latestCandle.Timestamp, nil
}
