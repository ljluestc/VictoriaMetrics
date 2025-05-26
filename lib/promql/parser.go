package promql

import (
	"fmt"
	"strconv"
)

// Define or import T_NUMBER
const T_NUMBER = 1 // Replace with the correct value if needed

type parser struct {
	// Add the missing lexer field
	lexer *lexer
}

func parseDuration(s string) (duration int64, tradingDays int64, err error) {
	var n float64
	var unit string
	for i, r := range s {
		if r >= '0' && r <= '9' || r == '.' {
			continue
		}
		n, err = strconv.ParseFloat(s[:i], 64)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid duration: %s", s)
		}
		unit = s[i:]
		break
	}
	switch unit {
	case T_DURATION_SECONDS:
		return int64(n * 1000), 0, nil
	case T_DURATION_MINUTES:
		return int64(n * 60 * 1000), 0, nil
	case T_DURATION_HOURS:
		return int64(n * 3600 * 1000), 0, nil
	case T_DURATION_DAYS:
		return int64(n * 24 * 3600 * 1000), 0, nil
	case T_DURATION_WEEKS:
		return int64(n * 7 * 24 * 3600 * 1000), 0, nil
	case T_DURATION_YEARS:
		return int64(n * 365 * 24 * 3600 * 1000), 0, nil
	case T_DURATION_TRADING_DAYS:
		return 0, int64(n), nil // Return number of trading days
	default:
		return 0, 0, fmt.Errorf("unknown duration unit: %s", unit)
	}
}

// Update parseOffset to handle trading days
func (p *parser) parseOffset() (offset int64, tradingDaysOffset int64, err error) {
	p.lexer.Next()
	if p.lexer.Token() != T_NUMBER {
		return 0, 0, fmt.Errorf("expected number after offset")
	}
	s := p.lexer.Text()
	p.lexer.Next()
	offset, tradingDaysOffset, err = parseDuration(s + p.lexer.Text())
	if err != nil {
		return 0, 0, err
	}
	return offset, tradingDaysOffset, nil
}
