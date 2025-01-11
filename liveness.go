package navm

import (
	"strconv"
)

// Liveness interval, half open [start, end)
type Interval struct {
	start int
	end   int
}

// Priority queue with smallest start values first
type LivenessQueue struct {
	intervals []Interval
}

func (q *LivenessQueue) Push(i Interval) {

	// Insert in sorted order
	if len(q.intervals) == 0 {
		q.intervals = append(q.intervals, i)
		return
	}
	for j, v := range q.intervals {
		if v.start > i.start {
			q.intervals = append(q.intervals[:j], append([]Interval{i}, q.intervals[j:]...)...)
			return
		}
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

func (q *LivenessQueue) Print() string {
	// Return string instead of printing
	str := "["
	for _, i := range q.intervals {
		str += "("
		str += strconv.Itoa(i.start)
		str += ", "
		str += strconv.Itoa(i.end)
		str += ")"
	}
	return str + "]"
}

func (q *LivenessQueue) Remove(i int) {
	q.intervals = append(q.intervals[:i], q.intervals[i+1:]...)
}
