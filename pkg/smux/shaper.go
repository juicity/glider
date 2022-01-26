package smux

func _itimediff(later, earlier uint32) int32 {
	return (int32)(later - earlier)
}

type shaperHeap []writeRequest

func (h shaperHeap) Len() int           { return len(h) }
func (h shaperHeap) Less(i, j int) bool { return _itimediff(h[j].prio, h[i].prio) > 0 }
func (h shaperHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *shaperHeap) Push(x any)        { *h = append(*h, x.(writeRequest)) }

func (h *shaperHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
