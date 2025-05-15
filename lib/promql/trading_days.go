package promql

import (
	"time"
)

// TradingDayConfig defines market-specific trading days and holidays.
type TradingDayConfig struct {
	MarketName string
	Holidays   map[time.Time]bool // Date (YYYY-MM-DD) -> true if holiday
}

// DefaultTradingDayConfig returns NYSE holidays for 2021.
func DefaultTradingDayConfig() *TradingDayConfig {
	holidays := map[time.Time]bool{
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC):   true, // New Year's Day
		time.Date(2021, 1, 18, 0, 0, 0, 0, time.UTC):  true, // MLK Day
		time.Date(2021, 2, 15, 0, 0, 0, 0, time.UTC):  true, // Presidents' Day
		time.Date(2021, 4, 2, 0, 0, 0, 0, time.UTC):   true, // Good Friday
		time.Date(2021, 5, 31, 0, 0, 0, 0, time.UTC):  true, // Memorial Day
		time.Date(2021, 7, 5, 0, 0, 0, 0, time.UTC):   true, // Independence Day (observed)
		time.Date(2021, 9, 6, 0, 0, 0, 0, time.UTC):   true, // Labor Day
		time.Date(2021, 11, 25, 0, 0, 0, 0, time.UTC): true, // Thanksgiving
		time.Date(2021, 12, 24, 0, 0, 0, 0, time.UTC): true, // Christmas (observed)
	}
	return &TradingDayConfig{
		MarketName: "NYSE",
		Holidays:   holidays,
	}
}

// Global configuration for trading days.
var globalTradingDayConfig = DefaultTradingDayConfig()

// IsTradingDay checks if the given date is a trading day (Mon-Fri, not a holiday).
func IsTradingDay(t time.Time, config *TradingDayConfig) bool {
	weekday := t.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}
	date := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	return !config.Holidays[date]
}

// GetTradingDayOffset calculates the timestamp N trading days before the given time.
func GetTradingDayOffset(t time.Time, n int64, config *TradingDayConfig) time.Time {
	if n <= 0 {
		return t
	}
	current := t
	tradingDaysCounted := int64(0)
	for tradingDaysCounted < n {
		current = current.AddDate(0, 0, -1)
		if IsTradingDay(current, config) {
			tradingDaysCounted++
		}
	}
	return current
}
