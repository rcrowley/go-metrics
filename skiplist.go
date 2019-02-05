package metrics

import (
	"fmt"
	"math"
	"sort"
)

type llNode struct {
	val  int64
	next *llNode
	prev *llNode
}

type expLL struct {
	head *llNode
	tail *llNode
	size int
	skip []*llNode
}

func (e *expLL) Clear() {
	e.head = nil
	e.tail = nil
	e.size = 0
}

func (e *expLL) Max() int64 {
	if e.tail != nil {
		return e.tail.val
	}
	return 0
}

func (e *expLL) Min() int64 {
	if e.head != nil {
		return e.head.val
	}
	return 0
}

func (e *expLL) Size() int {
	return e.size
}

func (e *expLL) Sum() int64 {
	sum := int64(0)
	for i := e.head; i != nil; i = i.next {
		sum += i.val
	}
	return sum
}

func (e *expLL) getAtPosition(pos int) int64 {
	div := pos / SkipInterval
	j := pos % SkipInterval
	if div > 0 {
		j++ // because first skipIterval is off by one
	}
	iter := e.skip[div]
	for i := 0; i < j; i++ {
		iter = iter.next
	}
	return iter.val
}

func (e *expLL) Percentiles(ps []float64) []float64 {
	scores := make([]float64, len(ps))
	if e.size > 0 {
		for i, p := range ps {
			pos := p * float64(e.size+1)
			if pos < 1.0 {
				scores[i] = float64(e.head.val)
			} else if pos >= float64(e.size) {
				scores[i] = float64(e.tail.val)
			} else {
				lower := float64(e.getAtPosition(int(pos) - 1))
				upper := float64(e.getAtPosition(int(pos)))
				scores[i] = lower + (pos-math.Floor(pos))*(upper-lower)
			}
		}
	}
	return scores
}

func (e *expLL) Variance() float64 {
	if 0 == e.size {
		return 0.0
	}
	m := float64(e.Sum()) / float64(e.size)
	var sum float64
	for i := e.head; i != nil; i = i.next {
		d := float64(i.val) - m
		sum += d * d
	}
	return sum / float64(e.size)
}

var dd = false

func logl(str ...interface{}) {
	if dd {
		fmt.Print(str...)
	}
}

func logln(str ...interface{}) {
	if dd {
		fmt.Println(str...)
	}
}

func (e *expLL) Values() []int64 {
	ret := make([]int64, 0, e.size)
	for i := e.head; i != nil; i = i.next {
		ret = append(ret, i.val)
	}
	return ret
}

func (e *expLL) Push(s int64) int {
	logln("pushing", s)
	tmp := &llNode{val: s}
	e.size++
	if e.head == nil {
		e.head, e.tail = tmp, tmp
		e.skip = append(e.skip, e.head)
		return 0 // no rank yet
	}

	idx, start := e.getIndexAndStart(s)

	j := idx * SkipInterval

	var i *llNode
	prev := start.prev
	for i = start; i != nil && i.val < s; i = i.next {
		prev = i
		j++
	}

	logln("exited loop with j", j, "and value", i, "and prev", prev)

	if prev == nil {
		logln("insert first element in list")
		tmp.next = e.head
		tmp.prev = nil
		e.head.prev = tmp
		e.head = tmp
	} else {
		// insert element before this one unless it's the end of the list
		logln("inserting after", prev.val)
		prev.next = tmp
		tmp.prev = prev
		if i == nil {
			e.tail = tmp
		} else {
			tmp.next = i
			i.prev = tmp
		}
	}
	if i == start {
		logln("need to move this skip too")
		idx-- // need to move this skip too
	}
	logln("moving skip indicies over, starting with idx ", idx+1)
	for ii := idx + 1; ii < len(e.skip); ii++ {
		logln(e.skip[ii].val, "->", e.skip[ii].prev)
		e.skip[ii] = e.skip[ii].prev
		if e.skip[ii] == nil {
			panic("we've just moved a skip index wrong to a nil")
		}
	}

	if e.size/SkipInterval >= len(e.skip) {
		if e.tail == nil {
			panic("e.tail is nil!")
		}
		logln("adding a skip size", e.size, len(e.skip))
		e.skip = append(e.skip, e.tail)
	}
	e.output()
	logln("returning", j, "/", e.size)
	return j
}

func (e *expLL) output() {
	if dd {
		logl("Current list [")
		for t := e.head; t != nil; t = t.next {
			logl(t.val, ",")
		}
		logln("]")
		logln("Current skip [")
		for ii, sk := range e.skip {
			logln("\tskip", ii, sk.val, sk.prev, sk.next)
		}
		logln("]")
		logln("head", e.head)
		logln("tail", e.tail)
	}
}

func (e *expLL) getIndexAndStart(v int64) (int, *llNode) {
	// find the node using a skiplist if exists
	idx := sort.Search(len(e.skip), func(i int) bool {
		return v < e.skip[i].val
	})
	logln("idx", idx, len(e.skip))
	var start *llNode
	if idx == 0 {
		start = e.head
	} else {
		idx--
		start = e.skip[idx]
	}
	logln("starting at value ", start.val)
	return idx, start
}

func (e *expLL) Remove(v int64) int64 {
	if e.head == nil {
		return 0
	}
	logln("removing", v)
	idx, start := e.getIndexAndStart(v)
	var i *llNode
	prev := start.prev
	for i = start; i != nil && i.val < v; i = i.next {
		prev = i
	}
	logln("removing", idx, i, start, prev)
	if i != nil && i.val == v {
		// update skip map
		if i == start {
			logln("removing the starting value from the skip list, move that too")
			idx-- // need to move this skip too
		}
		trim := false
		logln("moving skip list over")
		for ii := idx + 1; ii < len(e.skip); ii++ {
			logln(e.skip[ii], "->", e.skip[ii].next)
			if e.skip[ii].next == nil {
				if ii != len(e.skip)-1 {
					panic("nil next not on end")
				}
				logln("trimming off the end in remove of ", v, len(e.skip))
				trim = true
			} else {
				e.skip[ii] = e.skip[ii].next
			}
		}
		if trim {
			e.skip = e.skip[:len(e.skip)-1] // trim off the end
			logln("trimmed off the end in remove of ", v, len(e.skip))
		}

		if e.head == i {
			e.head = i.next
			if e.head != nil {
				e.head.prev = nil
			}
		}
		if e.tail == i {
			e.tail = i.prev
			if e.tail != nil {
				e.tail.next = nil
			}
		}
		if prev != nil {
			prev.next = i.next
			if i.next != nil {
				i.next.prev = prev
			}
		}
		i.next, i.prev = nil, nil
		e.size--
		e.output()
		return i.val
	}
	return 0
}

var SkipInterval = 100

func newExpLL(size int) *expLL {
	skipSize := size / SkipInterval
	return &expLL{
		skip: make([]*llNode, 0, skipSize),
	}
}
