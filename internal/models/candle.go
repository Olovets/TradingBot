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

func GetAllCandlesAfterTime(candles []Candle, timestamp int64) []Candle {
	var result []Candle
	for _, candle := range candles {
		if candle.Timestamp > timestamp {
			result = append(result, candle)
		}
	}
	return result
}

func (c *Candle) IsBullish() bool {
	return c.Close > c.Open
}

func (c *Candle) IsBearish() bool {
	return c.Close < c.Open
}

func (c *Candle) IsDoji() bool {
	return c.Close == c.Open
}

func (c *Candle) BullOrBearCandle() string {
	if c.Open < c.Close {
		return "bull"
	} else if c.Open > c.Close {
		return "bear"

	}
	return "doji"
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

func CutNonNeededCandles(candles []Candle) []Candle {
	firstCandleIndex := 0

	for i := 0; i < len(candles); i++ {

		// Convert timestamp to time.Time in UTC
		t := time.Unix(candles[i].Timestamp, 0).UTC()

		// Print the time
		fmt.Println("Time in UTC:", t)

		// Check if time is 20:45 (8:45 PM) UTC
		if t.Hour() == 23 && t.Minute() == 55 {
			firstCandleIndex = i + 1
			break
		} else {
			fmt.Println("❌ Timestamp is not at 20:45 UTC")
		}
	}

	return candles[firstCandleIndex:]

}

// GenerateCandles aggregates 5-minute candles into the desired period
func GenerateCandles(candles []Candle, period string, startTimestamp, endTimestamp int64) ([]Candle, error) {
	if endTimestamp == 0 {
		endTimestamp = 99999999999999
	}

	// Period durations in seconds
	var countOfCandles int

	switch period {
	case "15m":
		countOfCandles = 3
	case "1h":
		countOfCandles = 12
	case "4h":
		countOfCandles = 48
	case "1d":
		countOfCandles = 288
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}

	result := []Candle{}

	openPrice := 0.0
	closePrice := 0.0
	highPrice := 0.0
	lowPrice := 0.0
	lastTimestamp := int64(0)

	for i, candle := range candles {
		if i == 299 {
			fmt.Println((i + 1) % countOfCandles)
		}

		if (i+1)%countOfCandles == 1 {
			openPrice = candle.Open
			lowPrice = candle.Low
			highPrice = candle.High
			lastTimestamp = candle.Timestamp
		} else if (i+1)%countOfCandles == 0 {
			closePrice = candle.Close

			if candle.Low < lowPrice {
				lowPrice = candle.Low
			}

			if candle.High > highPrice {
				highPrice = candle.High
			}

			// Create a new candle with the aggregated data
			newCandle := Candle{
				Pair:      candle.Pair,
				Timestamp: lastTimestamp,
				Open:      openPrice,
				High:      highPrice,
				Low:       lowPrice,
				Close:     closePrice,
			}
			result = append(result, newCandle)

			openPrice = 0.0
			closePrice = 0.0
			highPrice = 0.0
			lowPrice = 0.0
			lastTimestamp = 0
		} else {
			if candle.High > highPrice {
				highPrice = candle.High
			}
			if candle.Low < lowPrice {
				lowPrice = candle.Low
			}

		}

	}

	filteredCandles := []Candle{}

	for _, candle := range result {
		if candle.Timestamp >= startTimestamp && candle.Timestamp <= endTimestamp {
			filteredCandles = append(filteredCandles, candle)
		}
	}

	return filteredCandles, nil
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
