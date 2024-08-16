package domain

import "time"

// TimeProvider is an interface to provide time
type TimeProvider interface {
	UtcNow() time.Time
}
