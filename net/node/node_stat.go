package node

import "sync/atomic"

type NodeStat struct {
	Inbound  StatHandler
	Outbound StatHandler
}

type StatHandler struct {
	Count uint64
}

func (h *StatHandler) Add(delta uint64) {
	atomic.AddUint64(&h.Count, delta)
}

func (h *StatHandler) AddOne() {
	atomic.AddUint64(&h.Count, 1)
}
