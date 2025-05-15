package storage

import (
	"sync"
	"time"

	"github.com/VictoriaMetrics/VictoriaMetrics/lib/logger"
)

var dataFlushInterval = 5 * time.Second

type partition struct {
	name                string
	s                   *Storage
	inmemoryParts       []*partWrapper
	smallParts          []*partWrapper
	bigParts            []*partWrapper
	partsLock           sync.Mutex
	inmemoryRowsDeleted counter
	smallRowsDeleted    counter
	bigRowsDeleted      counter
}

type partWrapper struct {
	isInMerge bool
	p         part
}

type counter struct {
	Value int
}

type partitionMetrics struct {
	// Define fields here
}

func (pt *partition) removeStaleParts() {
	if pt.s.keepOldData {
		logger.Infof("Skipping stale parts removal in partition %q due to -keepOldData=true", pt.name)
		return
	}
	startTime := time.Now()
	retentionDeadline := startTime.UnixMilli() - pt.s.retentionMsecs

	var pws []*partWrapper
	pt.partsLock.Lock()
	for _, pw := range pt.inmemoryParts {
		if !pw.isInMerge && pw.p.ph.MaxTimestamp < retentionDeadline {
			pt.inmemoryRowsDeleted.Add(pw.p.ph.RowsCount)
			pw.isInMerge = true
			pws = append(pws, pw)
		}
	}
	for _, pw := range pt.smallParts {
		if !pw.isInMerge && pw.p.ph.MaxTimestamp < retentionDeadline {
			pt.smallRowsDeleted.Add(pw.p.ph.RowsCount)
			pw.isInMerge = true
			pws = append(pws, pw)
		}
	}
	for _, pw := range pt.bigParts {
		if !pw.isInMerge && pw.p.ph.MaxTimestamp < retentionDeadline {
			pt.bigRowsDeleted.Add(pw.p.ph.RowsCount)
			pw.isInMerge = true
			pws = append(pws, pw)
		}
	}
	pt.partsLock.Unlock()

	pt.swapSrcWithDstParts(pws, nil, partSmall)
}

func (pt *partition) swapSrcWithDstParts(src, dst []*partWrapper, partType string) {
	// Placeholder implementation
}

const partSmall = "small"

func somePartitionFunction() []*partWrapper {
	parts := []*partWrapper{}
	c := counter{Value: 10}
	_ = c
	return parts
}
