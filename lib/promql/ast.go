package promql

import "github.com/prometheus/prometheus/model/labels"

type EvalConfig struct {
	StartTimestamp int64
	EndTimestamp   int64
	Storage        interface{}
}

type Vector []struct {
	Metric map[string]string
	Value  float64
}

const (
	T_DURATION_SECONDS = "s"
	T_DURATION_MINUTES = "m"
	T_DURATION_HOURS   = "h"
	T_DURATION_DAYS    = "d"
	T_DURATION_WEEKS   = "w"
	T_DURATION_YEARS   = "y"
)

type parser struct{}
