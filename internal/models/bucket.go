package models

import (
	"time"
)

type Bucket struct {
	CurrentCount int
	LastUpdate   time.Time
}
