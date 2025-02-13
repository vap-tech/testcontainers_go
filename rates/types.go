package rates

import (
	"fmt"
	"time"
)

type Rates struct {
	day   time.Time
	value float64
}

func (r Rates) String() string {
	return fmt.Sprintf("date: %s value: %.4f", r.day.Format("02-01-2006"), r.value)
}
