package common

import (
	"fmt"
	"runtime"
)

const (
	UniqueCode = "unique_code"
)

func MemUsage() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	result := fmt.Sprintf("memoryusage::Alloc = %v MB::TotalAlloc = %v MB::Sys = %v MB::tNumGC = %v", bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC)
	return result
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
