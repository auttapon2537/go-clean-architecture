package entity

import (
	"time"
)

// MemoryLog represents a memory log entity for MongoDB storage
type MemoryLog struct {
	ID            string    `json:"id" bson:"_id,omitempty"`
	Timestamp     time.Time `json:"timestamp" bson:"timestamp"`
	Alloc         uint64    `json:"alloc" bson:"alloc"`
	TotalAlloc    uint64    `json:"totalAlloc" bson:"totalAlloc"`
	Sys           uint64    `json:"sys" bson:"sys"`
	NumGC         uint32    `json:"numGC" bson:"numGC"`
	GCCPUFraction float64   `json:"gcCPUFraction" bson:"gcCPUFraction"`
	NumGoroutine  int       `json:"numGoroutine" bson:"numGoroutine"`
}
