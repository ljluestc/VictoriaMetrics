package prometheus

import (
	"fmt"
	"net/http"
)

// Removed duplicate declarations of queryDuration, QueryRangeHandler, and queryRangeHandler

var queryDuration = 0

func QueryRangeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "QueryRangeHandler")
}

func queryRangeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "queryRangeHandler")
}

func validatePrimaryVindexValues(values []interface{}) error {
	for _, value := range values {
		if value == nil {
			return fmt.Errorf("NULL values are not allowed for primary vindex columns")
		}
	}
	return nil
}

// Example usage in the insert logic
func processInsertQuery(query string, values []interface{}) error {
	// Validate primary vindex values
	if err := validatePrimaryVindexValues(values); err != nil {
		return err
	}

	// ...existing logic for processing the insert query...
	return nil
}
