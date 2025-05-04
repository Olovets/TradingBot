package blocks

import (
	"fmt"
	"github.com/Olovets/TradingBot/internal/models"
)

type RejectionBlock struct {
}

func NewRejectionBlock() *RejectionBlock {
	return &RejectionBlock{}
}

// Function to identify Fair Value Gaps (FVGs) and return them as a list of Blocks
func (b *RejectionBlock) IdentifyAll(candles []models.Candle) []models.Block {
	var blockList []models.Block

	// Iterate through the list of candles
	for i := 1; i < len(candles); i++ {

		if candles[i].Timestamp == 1742936400 {
			fmt.Println("")
		}

		// Get the previous and current candle
		prevCandle := candles[i-1]
		currCandle := candles[i]

		direction := ""

		if prevCandle.Open < prevCandle.Close && currCandle.Open > currCandle.Close {
			direction = models.Bear
		} else if prevCandle.Open > prevCandle.Close && currCandle.Open < currCandle.Close {
			direction = models.Bull
		}

		switch direction {
		case models.Bull:
			if prevCandle.Close-prevCandle.Low > 0 && currCandle.Open-currCandle.Low > 0 {
				low := prevCandle.Low
				if currCandle.Low < low {
					low = currCandle.Low
				}

				high := prevCandle.Close
				if currCandle.Open > high {
					high = currCandle.Open
				}

				block := models.Block{
					Type:      "RB",
					Low:       low,
					High:      high,
					Timestamp: int(currCandle.Timestamp),
					Valid:     true,
					Inverted:  false,
					Trend:     models.Bull,
				}
				blockList = append(blockList, block)
			}

		case models.Bear:
			if prevCandle.High-prevCandle.Close > 0 && currCandle.High-currCandle.Open > 0 {
				low := prevCandle.Close
				if currCandle.Open < low {
					low = currCandle.Low
				}

				high := prevCandle.High
				if currCandle.High > high {
					high = currCandle.High
				}

				block := models.Block{
					Type:      "RB",
					Low:       low,
					High:      high,
					Timestamp: int(currCandle.Timestamp),
					Valid:     true,
					Inverted:  false,
					Trend:     models.Bull,
				}
				blockList = append(blockList, block)
			}
		}

	}
	return blockList
}
