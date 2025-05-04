package blocks

import "github.com/Olovets/TradingBot/internal/models"

type BlockI interface {
	IdentifyAll(candles []models.Candle) []models.Block
}
