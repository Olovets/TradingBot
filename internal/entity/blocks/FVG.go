package blocks

import (
	"fmt"
	"github.com/Olovets/TradingBot/internal/models"
)

type FvgBlock struct {
}

func NewFvgBlock() *FvgBlock {
	return &FvgBlock{}
}

func (b *FvgBlock) IdentifyAll(candles []models.Candle) []models.Block {
	var blockList []models.Block

	// Iterate through the list of candles
	for i := 1; i < len(candles)-1; i++ {
		// Get the previous and current candle
		prevCandle := candles[i-1]
		currCandle := candles[i]
		nextCandle := candles[i+1]

		direction := ""

		if currCandle.ID == 33 {
			fmt.Println("")
		}

		if prevCandle.High < nextCandle.High && currCandle.Open < currCandle.Close {
			direction = models.Bull
		} else if prevCandle.Low > nextCandle.Low && currCandle.Open > currCandle.Close {
			direction = models.Bear
		}

		switch direction {
		case models.Bull:
			if prevCandle.High < nextCandle.Low {
				block := models.Block{
					Type:      "FVG",
					Low:       prevCandle.High,
					High:      nextCandle.Low,
					Timestamp: int(nextCandle.Timestamp),
					Valid:     true,
					Inverted:  false,
					Trend:     models.Bull,
				}
				blockList = append(blockList, block)
			}
		case models.Bear:
			if prevCandle.Low > nextCandle.High {
				block := models.Block{
					Type:      "FVG",
					Low:       nextCandle.High,
					High:      prevCandle.Low,
					Timestamp: int(nextCandle.Timestamp),
					Valid:     true,
					Inverted:  false,
					Trend:     models.Bear,
				}
				blockList = append(blockList, block)

			}
		}

	}
	return blockList
}
