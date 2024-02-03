package martinez_go

type SweepEventComparator func(e1, e2 *SweepEvent) int

// TinyQueueDefault represents a priority queue with an underlying slice.
type TinyQueueDefault struct {
	data    []*SweepEvent
	compare SweepEventComparator
}

// NewTinyQueueDefault creates a new instance of TinyQueueDefault
// It takes an optional slice of data and a compare function.
func NewTinyQueueDefault(data []*SweepEvent, compare SweepEventComparator) *TinyQueueDefault {
	if compare == nil {
		compare = defaultCompare
	}
	tq := &TinyQueueDefault{
		data:    data,
		compare: compare,
	}
	if len(data) > 0 {
		for i := len(data)>>1 - 1; i >= 0; i-- {
			tq.down(i)
		}
	}
	return tq
}

// defaultCompare provides a basic comparison function for integers.
func defaultCompare(e1, e2 *SweepEvent) int {
	return CompareEvents(e1, e2)
}

func (tq *TinyQueueDefault) Len() int {
	return len(tq.data)
}

// Push adds an item to the queue.
func (tq *TinyQueueDefault) Push(item *SweepEvent) {
	tq.data = append(tq.data, item)
	tq.up(len(tq.data) - 1)
}

// Pop removes and returns the top item from the queue.
func (tq *TinyQueueDefault) Pop() *SweepEvent {
	if len(tq.data) == 0 {
		return nil
	}
	top := tq.data[0]
	lastIndex := len(tq.data) - 1
	tq.data[0] = tq.data[lastIndex]
	tq.data = tq.data[:lastIndex]
	if len(tq.data) > 0 {
		tq.down(0)
	}
	return top
}

// Peek returns the top item from the queue without removing it.
func (tq *TinyQueueDefault) Peek() *SweepEvent {
	if len(tq.data) == 0 {
		return nil
	}
	return tq.data[0]
}

// up moves the item at the given position up to its proper place in the queue.
func (tq *TinyQueueDefault) up(pos int) {
	item := tq.data[pos]
	for pos > 0 {
		parent := (pos - 1) >> 1
		if tq.compare(item, tq.data[parent]) >= 0 {
			break
		}
		tq.data[pos] = tq.data[parent]
		pos = parent
	}
	tq.data[pos] = item
}

// down moves the item at the given position down to its proper place in the queue.
func (tq *TinyQueueDefault) down(pos int) {
	item := tq.data[pos]
	halfLength := len(tq.data) >> 1
	for pos < halfLength {
		left := (pos << 1) + 1
		right := left + 1
		best := left
		if right < len(tq.data) && tq.compare(tq.data[right], tq.data[left]) < 0 {
			best = right
		}
		if tq.compare(tq.data[best], item) >= 0 {
			break
		}
		tq.data[pos] = tq.data[best]
		pos = best
	}
	tq.data[pos] = item
}

// TinyQueueDefault represents a priority queue with an underlying slice.
type TinyQueue interface {
	Len() int
	Push(*SweepEvent)
	Pop() *SweepEvent
	Peek() *SweepEvent
}
