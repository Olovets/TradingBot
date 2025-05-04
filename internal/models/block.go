package models

const Bull = "bull"
const Bear = "bear"

type Block struct {
	Type          string  `json:"type"`
	Low           float64 `json:"low"`
	High          float64 `json:"high"`
	Timestamp     int     `json:"timestamp"`
	Valid         bool    `json:"valid"`
	Inverted      bool    `json:"inverted"`
	Trend         string  `json:"trend"`
	ConfirmedTime int     `json:"confirmed_time"`
	InvertTime    int     `json:"invert_time"`
}

func ReturnLast(candles []Block) Block {
	lastBlock := Block{}
	for _, candle := range candles {
		if candle.Timestamp > lastBlock.Timestamp {
			lastBlock = candle
		}
	}
	return lastBlock
}

func AggregateBlocks(blocks [][]Block) []Block {
	aggregatedBlocks := []Block{}
	for _, blockList := range blocks {
		aggregatedBlocks = append(aggregatedBlocks, blockList...)
	}
	return aggregatedBlocks
}

func ReturnNearFTA(POI Block, blocks []Block) Block {
	switch POI.Trend {
	case Bull:
		ftaBlocks := []Block{}

		for _, block := range blocks {
			if block.Trend == Bear {
				ftaBlocks = append(ftaBlocks, block)
			}
		}

		lastFta := Block{}
		lastNearPrice := 0.0

		for _, block := range ftaBlocks {
			if (lastNearPrice == 0.0 && block.Low > POI.Low) || (block.Low > POI.Low && block.Low < lastNearPrice) {
				lastNearPrice = block.Low
				lastFta = block
			}
		}

		return lastFta
	case Bear:
		ftaBlocks := []Block{}

		for _, block := range blocks {
			if block.Trend == Bear {
				ftaBlocks = append(ftaBlocks, block)
			}
		}

		lastFta := Block{}
		lastNearPrice := 0.0

		for _, block := range ftaBlocks {
			if (lastNearPrice == 0.0 && block.High < POI.High) || (block.High < POI.High && block.High > lastNearPrice) {
				lastNearPrice = block.High
				lastFta = block
			}
		}

		return lastFta
	}

	return Block{}
}
