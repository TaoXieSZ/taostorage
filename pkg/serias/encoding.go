package serias

import (
	"sync"
)

type Serias struct {
	sync.Mutex
	T0  uint64
	t   uint64
	val float64
}
