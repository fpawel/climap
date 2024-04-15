package main

import "sync/atomic"

type Connections struct {
	*atomic.Int64
}

func (x Connections) Inc() int64 {
	return x.Int64.Add(1)
}

func (x Connections) Get() int64 {
	return x.Int64.Load()
}

func (x Connections) Dec() int64 {
	return x.Int64.Add(-1)
}
