package navm

import (
	"strconv"
)

// Liveness interval, half open [start, end)
type Interval struct {
	physicalRegister int
	register         Register
	start            int
	end              int
}

// Two modes for different uses in the register allocation linear scan
// If active, we sort by start time ascending as is currently the case
// If inactive, we sort by end time ascending

// Priority queue with smallest start values first
type LivenessQueue struct {
	intervals []Interval
	active    bool
}

func (q *LivenessQueue) Push(i Interval) {
	// Insert in sorted order
	if len(q.intervals) == 0 {
		q.intervals = append(q.intervals, i)
		return
	}
	if q.active {
		for j, v := range q.intervals {
			if v.start > i.start {
				q.intervals = append(q.intervals[:j], append([]Interval{i}, q.intervals[j:]...)...)
				return
			}
		}
		q.intervals = append(q.intervals, i)

	} else {
		for j, v := range q.intervals {
			if v.end > i.end {
				q.intervals = append(q.intervals[:j], append([]Interval{i}, q.intervals[j:]...)...)
				return
			}
		}
		q.intervals = append(q.intervals, i)
	}
}

func (q *LivenessQueue) Pop() Interval {
	if len(q.intervals) == 0 {
		panic("Empty queue")
	}
	i := q.intervals[0]
	q.intervals = q.intervals[1:]
	return i
}

func (q *LivenessQueue) PopLast() Interval {
	if len(q.intervals) == 0 {
		panic("Empty queue")
	}
	i := q.intervals[len(q.intervals)-1]
	q.intervals = q.intervals[:len(q.intervals)-1]
	return i
}

func (q *LivenessQueue) Len() int {
	return len(q.intervals)
}

func (q *LivenessQueue) Empty() bool {
	return len(q.intervals) == 0
}

func (q *LivenessQueue) Peek() Interval {
	if len(q.intervals) == 0 {
		panic("Empty queue")
	}
	return q.intervals[0]
}

func (q *Interval) Print() string {
	return "(" + strconv.Itoa(q.start) + ", " + strconv.Itoa(q.end) + ")"
}

func (q *LivenessQueue) Print() string {
	// Return string instead of printing
	str := "["
	for _, i := range q.intervals {
		str += i.Print()
	}
	return str + "]"
}

func (q *LivenessQueue) Remove(i int) {
	q.intervals = append(q.intervals[:i], q.intervals[i+1:]...)
}
