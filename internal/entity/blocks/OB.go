package blocks

import (
	"github.com/Olovets/TradingBot/common"
	"github.com/Olovets/TradingBot/internal/models"
)

type OrderBlock struct {
}

func NewOrderBlock() *OrderBlock {
	return &OrderBlock{}
}

func (b *OrderBlock) IdentifyAll(candles []models.Candle) []models.Block {
	ranges := IdentifyOBRanges(candles)
	confirmedRanges := ConfirmOBs(candles, ranges)

	OBs := CheckInvertingOBs(candles, confirmedRanges)

	OBs = InvalidOBs(candles, OBs)

	return OBs
}

func InvalidOBs(candles []models.Candle, blocks []models.Block) []models.Block {
	validatedBlocks := make([]models.Block, 0)

	for i := 0; i < len(blocks); i++ {
		if blocks[i].InvertTime == 0 {
			validatedBlocks = append(validatedBlocks, blocks[i])
		}

		candles = models.GetAllCandlesAfterTime(candles, int64(blocks[i].InvertTime))

		if blocks[i].Trend == models.Bull {
			invalidValue := blocks[i].Low

			for j := 0; j < len(candles); j++ {
				if candles[j].Close < invalidValue {
					break
				}
			}

		} else if blocks[i].Trend == models.Bear {
			invalidValue := blocks[i].High

			for j := 0; j < len(candles); j++ {
				if candles[j].Close > invalidValue {
					break
				}
			}
		}
	}
	return validatedBlocks
}

func CheckInvertingOBs(candles []models.Candle, blocks []models.Block) []models.Block {
	confirmedBlocks := make([]models.Block, 0)

	for i := 0; i < len(blocks); i++ {
		candles = models.GetAllCandlesAfterTime(candles, int64(blocks[i].Timestamp))

		if blocks[i].Trend == models.Bull {
			confirmValue := blocks[i].Low
			invalidValue := blocks[i].High

			for j := 0; j < len(candles); j++ {
				if candles[j].Close < confirmValue {
					blocks[i].ConfirmedTime = int(candles[j].Timestamp)
					blocks[i].Trend = models.Bear
					confirmedBlocks = append(confirmedBlocks, blocks[i])
					break
				}

				if candles[j].Close > invalidValue {
					break
				}
			}

		} else if blocks[i].Trend == models.Bear {
			confirmValue := blocks[i].High
			invalidValue := blocks[i].Low

			for j := 0; j < len(candles); j++ {
				if candles[j].Close > confirmValue {
					blocks[i].ConfirmedTime = int(candles[j].Timestamp)
					blocks[i].Trend = models.Bull

					confirmedBlocks = append(confirmedBlocks, blocks[i])
					break
				}

				if candles[j].Close < invalidValue {
					break
				}
			}

		}

		confirmedBlocks = append(confirmedBlocks, blocks[i])
	}
	return confirmedBlocks
}

func ConfirmOBs(candles []models.Candle, blocks []models.Block) []models.Block {
	confirmedBlocks := make([]models.Block, 0)

	for i := 0; i < len(blocks); i++ {
		candles = models.GetAllCandlesAfterTime(candles, int64(blocks[i].Timestamp))

		if blocks[i].Trend == models.Bull {
			confirmValue := blocks[i].Low
			invalidValue := blocks[i].High

			for j := 0; j < len(candles); j++ {
				if candles[j].Close < confirmValue {
					blocks[i].ConfirmedTime = int(candles[j].Timestamp)
					blocks[i].Trend = models.Bear
					confirmedBlocks = append(confirmedBlocks, blocks[i])
					break
				}

				if candles[j].Close > invalidValue {
					break
				}
			}

		} else if blocks[i].Trend == models.Bear {
			confirmValue := blocks[i].High
			invalidValue := blocks[i].Low

			for j := 0; j < len(candles); j++ {
				if candles[j].Close > confirmValue {
					blocks[i].ConfirmedTime = int(candles[j].Timestamp)
					blocks[i].Trend = models.Bull

					confirmedBlocks = append(confirmedBlocks, blocks[i])
					break
				}

				if candles[j].Close < invalidValue {
					break
				}
			}

		}
	}
	return confirmedBlocks
}

// Function to identify Fair Value Gaps (FVGs) and return them as a list of Blocks
func IdentifyOBRanges(candles []models.Candle) []models.Block {
	var alreadyUsed []uint
	blockList := make([]models.Block, 0)

	// Iterate through the list of candles
	for i := 0; i < len(candles); i++ {
		if common.InUintSlice(alreadyUsed, uint(candles[i].Timestamp)) {
			continue
		}

		trend := candles[i].BullOrBearCandle()

		if trend == "bull" {
			low := candles[i].Open
			high := candles[i].Close
			timestamp := candles[i].Timestamp

			for j := i + 1; j < len(candles); j++ {
				if candles[j].IsBearish() {
					high = candles[j].Open
					break
				} else {
					alreadyUsed = append(alreadyUsed, uint(candles[j].Timestamp))
					timestamp = candles[j].Timestamp
				}
			}

			block := models.Block{
				Type:      "OB",
				Low:       low,
				High:      high,
				Timestamp: int(timestamp),
				Valid:     true,
				Inverted:  false,
				Trend:     models.Bull,
			}
			blockList = append(blockList, block)
		} else if trend == "bear" {
			high := candles[i].Open
			low := candles[i].Close
			timestamp := candles[i].Timestamp

			for j := i + 1; j < len(candles); j++ {
				if candles[j].IsBullish() {
					low = candles[j].Open
					break
				} else {
					alreadyUsed = append(alreadyUsed, uint(candles[j].Timestamp))
					timestamp = candles[j].Timestamp
				}
			}

			block := models.Block{
				Type:      "OB",
				Low:       low,
				High:      high,
				Timestamp: int(timestamp),
				Valid:     true,
				Inverted:  false,
				Trend:     models.Bear,
			}
			blockList = append(blockList, block)
		} else {
			continue
		}

	}

	return blockList
}
