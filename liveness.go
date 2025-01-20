package navm

import (
	"strconv"
)

// Liveness interval, half open [start, end)
// Can be in one of 3 states:
// 1. virtual register (not yet allocated)
// 2. physical register (allocated)
// 3. Stack position (spilled)
type Interval struct {
	physicalRegister int
	stackPosition    int
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
	if !q.active {
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
	return strconv.Itoa(q.register.value) + "->" + strconv.Itoa(q.physicalRegister) +
		"(" + strconv.Itoa(q.start) + ", " + strconv.Itoa(q.end) + ")"
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

func makeIntervals(ir *IR) []Interval {
	intervals := make([]Interval, ir.registersLength)
	// always skip first, because 0th register is unused
	// range from 1 to len(intervals)-1
	for i := 1; i < len(intervals); i++ {
		intervals[i] = Interval{register: Register{registerType: virtualRegister, value: i}}
		// Set start to max
		intervals[i].start = len(ir.instructions)
	}

	for i, instr := range ir.instructions {
		// Get all virtual registers used in this instruction
		// and update their intervals
		if instr.arg1.registerType == virtualRegister {
			intervals[instr.arg1.value].start = min(intervals[instr.arg1.value].start, i)
			intervals[instr.arg1.value].end = i + 1
		}
		if instr.ret.registerType == virtualRegister {
			intervals[instr.ret.value].start = min(intervals[instr.ret.value].start, i)
			intervals[instr.ret.value].end = i + 1
		}
		if instr.arg2.isVirtualRegister && (instr.arg2.argType == registerArg || instr.arg2.argType == address) {
			intervals[instr.arg2.value].start = min(intervals[instr.arg2.value].start, i)
			intervals[instr.arg2.value].end = i + 1
		}
	}

	return intervals
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
