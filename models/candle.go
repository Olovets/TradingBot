package models

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

// Generic Candle structure for any trading pair
type Candle struct {
	ID        uint    `gorm:"primarykey" json:"id,omitempty"`
	Pair      string  `gorm:"not null" json:"pair,omitempty"`      // Currency pair (e.g., EURUSD, GBPUSD)
	Timestamp int64   `gorm:"not null" json:"timestamp,omitempty"` // Unix timestamp for the candle's close time
	Open      float64 `gorm:"not null" json:"open,omitempty"`
	High      float64 `gorm:"not null" json:"high,omitempty"`
	Low       float64 `gorm:"not null" json:"low,omitempty"`
	Close     float64 `gorm:"not null" json:"close,omitempty"`
}

// GetAllCandles retrieves all candles from the database
func GetAllCandles(db *gorm.DB) ([]Candle, error) {
	var candles []Candle

	// Query all candles from the database
	err := db.Find(&candles).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching candles from database: %v", err)
	}

	return candles, nil
}

// GenerateCandles aggregates 5-minute candles into the desired period
func GenerateCandles(candles []Candle, period string, startTimestamp, endTimestamp int64) ([]Candle, error) {
	// Period durations in seconds
	var periodSeconds int64
	switch period {
	case "15m":
		periodSeconds = 15 * 60
	case "1h":
		periodSeconds = 60 * 60
	case "4h":
		periodSeconds = 4 * 60 * 60
	case "1d":
		periodSeconds = 24 * 60 * 60
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}

	// Align timestamps to the session start at 21:00 UTC
	sessionStart := alignToSessionStart(startTimestamp)
	sessionEnd := alignToSessionStart(endTimestamp)

	var result []Candle
	var currentCandle *Candle

	for _, candle := range candles {
		// Skip candles outside the start and end range
		if candle.Timestamp < sessionStart || candle.Timestamp > sessionEnd {
			continue
		}

		// Calculate the aligned candle start time for the given period
		candleStart := alignToPeriodStart(candle.Timestamp, periodSeconds)

		// If currentCandle is nil or the timestamp doesn't match the current aggregated candle
		if currentCandle == nil || currentCandle.Timestamp != candleStart {
			// Save the completed candle and start a new one
			if currentCandle != nil {
				result = append(result, *currentCandle)
			}

			// Initialize a new candle
			currentCandle = &Candle{
				Pair:      candle.Pair,
				Timestamp: candleStart,
				Open:      candle.Open,
				High:      candle.High,
				Low:       candle.Low,
				Close:     candle.Close,
			}
		} else {
			// Update the current candle with the aggregated data
			currentCandle.High = max(currentCandle.High, candle.High)
			currentCandle.Low = min(currentCandle.Low, candle.Low)
			currentCandle.Close = candle.Close
		}
	}

	// Append the last candle
	if currentCandle != nil {
		result = append(result, *currentCandle)
	}

	return result, nil
}

// Helper function to align a timestamp to the session start at 21:00 UTC
func alignToSessionStart(timestamp int64) int64 {
	utcTime := time.Unix(timestamp, 0).UTC()
	sessionStart := time.Date(utcTime.Year(), utcTime.Month(), utcTime.Day(), 21, 0, 0, 0, time.UTC)
	if utcTime.Before(sessionStart) {
		sessionStart = sessionStart.AddDate(0, 0, -1) // Move to the previous day's session start
	}
	return sessionStart.Unix()
}

// Helper function to align a timestamp to the nearest period start
func alignToPeriodStart(timestamp, periodSeconds int64) int64 {
	sessionStart := alignToSessionStart(timestamp)
	offset := (timestamp - sessionStart) % periodSeconds
	return timestamp - offset
}

// Helper function to find the maximum of two float64 values
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// Helper function to find the minimum of two float64 values
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
