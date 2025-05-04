// internal/marketdata/alpha_vantage.go
package marketdata

import (
	"encoding/json"
	"fmt"
	"github.com/Olovets/TradingBot/internal/models"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

const binanceAPIURL = "https://api.binance.com/api/v3/klines"

// Helper function to round to 5 decimal places
func roundToFiveDecimals(value float64) float64 {
	factor := math.Pow(10, 5) // 10^5 for 5 decimal places
	return math.Round(value*factor) / factor
}

func FetchCandlesLast31Days(pair string, apiKey string) ([]models.Candle, error) {
	const limit = 500                            // Maximum number of candles per request
	const intervalDuration = 5 * 60              // 5-minute interval in seconds
	const oneDay = 24 * 60 * 60                  // One day in seconds
	const daysToFetch = 31                       // Fetch data for 31 days
	const intervalInMilliseconds = 5 * 60 * 1000 // 5-minute interval in milliseconds

	endTime := time.Now().Unix() * 1000                  // Current time in milliseconds
	startTime := endTime - (daysToFetch * oneDay * 1000) // Start time 31 days ago in milliseconds

	var candles []models.Candle

	for {
		// Build the API request URL with startTime, endTime, and limit
		params := fmt.Sprintf("?symbol=%s&interval=5m&limit=%d&startTime=%d&endTime=%d",
			pair, limit, startTime, endTime)
		url := binanceAPIURL + params

		// Create the HTTP request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating HTTP request: %v", err)
		}

		// Add the API key to the headers (if required)
		if apiKey != "" {
			req.Header.Add("X-MBX-APIKEY", apiKey)
		}

		// Make the HTTP request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error fetching data from Binance: %v", err)
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %v", err)
		}

		// Parse the JSON response
		var apiResponse [][]interface{}
		if err := json.Unmarshal(body, &apiResponse); err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON response: %v", err)
		}

		// Break the loop if no more data is returned
		if len(apiResponse) == 0 {
			break
		}

		// Convert the API response into Candle structs
		for _, entry := range apiResponse {
			if len(entry) < 6 {
				continue
			}

			timestamp := int64(entry[0].(float64)) // Binance timestamps are in milliseconds

			candle := models.Candle{
				Pair:      pair,
				Timestamp: timestamp / 1000, // Store in seconds
				Open:      roundToFiveDecimals(parseFloat(entry[1].(string))),
				High:      roundToFiveDecimals(parseFloat(entry[2].(string))),
				Low:       roundToFiveDecimals(parseFloat(entry[3].(string))),
				Close:     roundToFiveDecimals(parseFloat(entry[4].(string))),
			}
			candles = append(candles, candle)
		}

		// Update startTime to fetch the next batch of candles
		if len(apiResponse) < limit {
			break // If fewer candles are returned, no more data is available
		}
		startTime += int64(limit * intervalInMilliseconds)
	}

	return candles, nil
}

// FetchCandles fetches 5-minute candles for a given trading pair from Binance
func FetchCandles(pair string, apiKey string, since int64) ([]models.Candle, error) {
	var candles []models.Candle
	const limit = 500                      // Maximum number of candles per request
	var lastTimestamp int64 = since * 1000 // Convert `since` to milliseconds

	for {
		// Build the API request URL with startTime and limit
		params := fmt.Sprintf("?symbol=%s&interval=5m&limit=%d", pair, limit)
		if lastTimestamp > 0 {
			params += fmt.Sprintf("&startTime=%d", lastTimestamp)
		}
		url := binanceAPIURL + params

		// Create the HTTP request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating HTTP request: %v", err)
		}

		// Add the API key to the headers (if required)
		if apiKey != "" {
			req.Header.Add("X-MBX-APIKEY", apiKey)
		}

		// Make the HTTP request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error fetching data from Binance: %v", err)
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %v", err)
		}

		// Parse the JSON response
		var apiResponse [][]interface{}
		if err := json.Unmarshal(body, &apiResponse); err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON response: %v", err)
		}

		// Break the loop if no more data is returned
		if len(apiResponse) == 0 {
			break
		}

		// Convert the API response into Candle structs
		for _, entry := range apiResponse {
			if len(entry) < 6 {
				continue
			}

			timestamp := int64(entry[0].(float64)) // Binance timestamps are in milliseconds
			candleTime := time.Unix(timestamp/1000, 0)

			// Skip weekends if simulating forex behavior
			weekday := candleTime.Weekday()
			if weekday == time.Saturday || weekday == time.Sunday {
				continue
			}

			candle := models.Candle{
				Pair:      pair,
				Timestamp: timestamp / 1000, // Store in seconds
				Open:      roundToFiveDecimals(parseFloat(entry[1].(string))),
				High:      roundToFiveDecimals(parseFloat(entry[2].(string))),
				Low:       roundToFiveDecimals(parseFloat(entry[3].(string))),
				Close:     roundToFiveDecimals(parseFloat(entry[4].(string))),
			}
			candles = append(candles, candle)
		}

		// Update the `lastTimestamp` to fetch the next batch of candles
		lastTimestamp = int64(apiResponse[len(apiResponse)-1][0].(float64)) + 1

		// Exit the loop if the maximum candle count is reached (optional)
		if len(apiResponse) < limit {
			break
		}
	}

	return candles, nil
}

// parseFloat safely parses a string into a float64
func parseFloat(value string) float64 {
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Printf("Error parsing float: %v", err)
		return 0
	}
	return result
}
