package prometheus

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/VictoriaMetrics/VictoriaMetrics/lib/logger"
)

// RegisterHandlers registers Prometheus-compatible API handlers.
func RegisterHandlers() {
	http.HandleFunc("/api/v1/series/count", seriesCountHandler)
	// Other handlers...
}

// seriesCountHandler handles /api/v1/series/count requests.
func seriesCountHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("cannot parse form: %s", err), http.StatusBadRequest)
		return
	}

	// Parse match[] parameters
	matchers, err := parseSeriesSelectors(r.Form["match[]"])
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot parse match[]: %s", err), http.StatusBadRequest)
		return
	}

	// Parse start and end timestamps
	start, end, err := parseTimeRange(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot parse time range: %s", err), http.StatusBadRequest)
		return
	}

	// Get parts from the table
	db := getDB() // Assume a function to get the DB instance
	db.incRef()
	defer db.decRef()
	tb := db.tb
	pws := tb.getParts(nil)
	defer tb.putParts(pws)

	// Count matching series
	count := countMatchingSeries(pws, matchers, start, end)

	// Return JSON response
	response := fmt.Sprintf(`{"status":"success","data":%d}`, count)
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write([]byte(response)); err != nil {
		logger.Errorf("cannot write response: %s", err)
	}
}

// parseSeriesSelectors parses match[] parameters into tag filters.
func parseSeriesSelectors(match []string) ([]*TagFilter, error) {
	if len(match) == 0 {
		// No match[] provided; return nil to count all series
		return nil, nil
	}

	var filters []*TagFilter
	for _, m := range match {
		// Parse series selector (e.g., {__name__="metric1",job=~"app.*"})
		tf, err := parseSeriesSelector(m)
		if err != nil {
			return nil, fmt.Errorf("invalid series selector %q: %w", m, err)
		}
		filters = append(filters, tf)
	}
	return filters, nil
}

// parseSeriesSelector parses a single series selector (mock implementation).
func parseSeriesSelector(s string) (*TagFilter, error) {
	// Basic parsing for {key="value"}, {key=~"regex"}
	s = strings.Trim(s, "{}")
	pairs := strings.Split(s, ",")
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid selector %q", pair)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), "\"")
		isRegexp := strings.HasPrefix(parts[1], "~")
		if isRegexp {
			value = strings.TrimPrefix(value, "~")
		}
		if key == "__name__" { // Handle __name__ for simplicity
			return &TagFilter{
				Key:        []byte(key),
				Value:      []byte(value),
				IsRegexp:   isRegexp,
				IsNegative: false,
			}, nil
		}
	}
	return nil, fmt.Errorf("no __name__ selector found in %q", s)
}

// parseTimeRange parses start and end query parameters.
func parseTimeRange(r *http.Request) (start, end time.Time, err error) {
	// Default to last 5 minutes if not specified (consistent with /api/v1/series)
	end = time.Now()
	start = end.Add(-5 * time.Minute)

	if startStr := r.FormValue("start"); startStr != "" {
		start, err = parseTimestamp(startStr)
		if err != nil {
			return start, end, fmt.Errorf("invalid start timestamp: %w", err)
		}
	}

	if endStr := r.FormValue("end"); endStr != "" {
		end, err = parseTimestamp(endStr)
		if err != nil {
			return start, end, fmt.Errorf("invalid end timestamp: %w", err)
		}
	}

	if start.After(end) {
		return start, end, fmt.Errorf("start time %s is after end time %s", start, end)
	}

	return start, end, nil
}

// parseTimestamp parses RFC3339 or Unix timestamp.
func parseTimestamp(s string) (time.Time, error) {
	// Try RFC3339
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	// Try Unix timestamp
	if ts, err := strconv.ParseInt(s, 10, 64); err == nil {
		return time.Unix(ts, 0), nil
	}
	return time.Time{}, fmt.Errorf("cannot parse timestamp %q", s)
}

// Ensure correct usage of TagFilter and PartWrapper
func countMatchingSeries(pws []*PartWrapper, filters []*TagFilter, start, end time.Time) int {
	uniqueSeries := make(map[string]struct{})

	for _, pw := range pws {
		if pw.PH.MinTimestamp > end.Unix()*1000 || pw.PH.MaxTimestamp < start.Unix()*1000 {
			continue
		}

		if len(filters) == 0 {
			for i := uint64(0); i < pw.PH.ItemsCount; i++ {
				uniqueSeries[fmt.Sprintf("%d-%d", pw.PH.MinTimestamp, i)] = struct{}{}
			}
			continue
		}

		for _, tf := range filters {
			items := SearchItemsInPart(nil, pw, tf)
			for _, item := range items {
				uniqueSeries[string(item)] = struct{}{}
			}
		}
	}

	return len(uniqueSeries)
}

// Replace TagFilter with a mock implementation or correct import
type TagFilter struct {
	Key        []byte
	Value      []byte
	IsRegexp   bool
	IsNegative bool
}

// Replace PartWrapper with a mock implementation or correct import
type PartWrapper struct {
	PH struct {
		MinTimestamp int64
		MaxTimestamp int64
		ItemsCount   uint64
	}
}

// Replace DB with a mock implementation or correct import
type DB struct {
	tb *Table
}

func (db *DB) incRef() {}
func (db *DB) decRef() {}

type Table struct{}

func (tb *Table) getParts(pws []*PartWrapper) []*PartWrapper {
	return []*PartWrapper{}
}

func (tb *Table) putParts(pws []*PartWrapper) {}

// Replace SearchItemsInPart with a mock implementation
func SearchItemsInPart(dst []string, pw *PartWrapper, tf *TagFilter) []string {
	return []string{"mockItem"}
}

// Replace getDB with a mock implementation
func getDB() *DB {
	return &DB{tb: &Table{}}
}
