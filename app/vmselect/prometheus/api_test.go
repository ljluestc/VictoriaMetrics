package prometheus

import (
	"testing"
)

func TestValidatePrimaryVindexValues(t *testing.T) {
	tests := []struct {
		values    []interface{}
		expectErr bool
	}{
		{[]interface{}{1, "test", 3.14}, false},  // Valid values
		{[]interface{}{nil, "test", 3.14}, true}, // Contains NULL
		{[]interface{}{1, nil, 3.14}, true},      // Contains NULL
	}

	for _, test := range tests {
		err := validatePrimaryVindexValues(test.values)
		if (err != nil) != test.expectErr {
			t.Errorf("validatePrimaryVindexValues(%v) = %v, expectErr %v", test.values, err, test.expectErr)
		}
	}
}
