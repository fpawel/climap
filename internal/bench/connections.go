package bench

import "sync/atomic"

type connections struct {
	*atomic.Int64
}

func (x connections) Inc() int64 {
	return x.Int64.Add(1)
}

func (x connections) Get() int64 {
	return x.Int64.Load()
}

func (x connections) Dec() int64 {
	return x.Int64.Add(-1)
}
