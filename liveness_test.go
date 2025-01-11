package navm

import (
	"testing"
)

func init() {
}

func TestLivenessQueueActiveStartsWithSmall(t *testing.T) {
	q := LivenessQueue{}
	q.active = true
	q.Push(Interval{start: 0, end: 1})
	q.Push(Interval{start: 3, end: 4})
	q.Push(Interval{start: 1, end: 2})
	q.Push(Interval{start: 2, end: 3})

	if q.Len() != 4 {
		t.Errorf("Expected 4, got %d", q.Len())
	}
	if q.Pop().start != 0 {
		t.Errorf("Expected 0, got %d", q.Pop().start)
	}
	if q.Pop().start != 1 {
		t.Errorf("Expected 1, got %d", q.Pop().start)
	}
	if q.Pop().start != 2 {
		t.Errorf("Expected 2, got %d", q.Pop().start)
	}
	if q.Pop().start != 3 {
		t.Errorf("Expected 3, got %d", q.Pop().start)
	}
	if q.Empty() != true {
		t.Errorf("Expected true, got false")
	}
}

func TestLivenessQueueActive(t *testing.T) {
	// Push a few values and check that they are sorted
	q := LivenessQueue{}
	q.active = true
	q.Push(Interval{start: 3, end: 4})
	q.Push(Interval{start: 1, end: 2})
	q.Push(Interval{start: 2, end: 3})
	q.Push(Interval{start: 0, end: 1})
	if q.intervals[0].start != 0 {
		t.Errorf("Expected 0, got %d", q.intervals[0].start)
	}
	if q.intervals[1].start != 1 {
		t.Errorf("Expected 1, got %d", q.intervals[1].start)
	}
	if q.intervals[2].start != 2 {
		t.Errorf("Expected 2, got %d", q.intervals[2].start)
	}
	if q.intervals[3].start != 3 {
		t.Errorf("Expected 3, got %d", q.intervals[3].start)
	}

	// Test pop
	i := q.Pop()
	if i.start != 0 {
		t.Errorf("Expected 0, got %d", i.start)
	}
	if q.Len() != 3 {
		t.Errorf("Expected 3, got %d", q.Len())
	}
	if q.Peek().start != 1 {
		t.Errorf("Expected 1, got %d", q.Peek().start)
	}

	// Test remove
	q.Remove(1)

	if q.Len() != 2 {
		t.Errorf("Expected 2, got %d", q.Len())
	}
	if q.Peek().start != 1 {
		t.Errorf("Expected 1, got %d", q.Peek().start)
	}

	j := q.Pop()
	if j.start != 1 {
		t.Errorf("Expected 1, got %d", j.start)
	}

	if q.Peek().start != 3 {
		t.Errorf("Expected 3, got %d", q.Peek().start)
	}
}

func TestLivenessQueueInactive(t *testing.T) {
	// Push a few values and check that they are sorted
	q := LivenessQueue{}
	q.active = false
	q.Push(Interval{start: 3, end: 4})
	q.Push(Interval{start: 1, end: 2})
	q.Push(Interval{start: 2, end: 3})
	q.Push(Interval{start: 0, end: 1})
	if q.intervals[0].end != 4 {
		t.Errorf("Expected 4, got %d", q.intervals[0].end)
	}
	if q.intervals[1].end != 3 {
		t.Errorf("Expected 3, got %d", q.intervals[1].end)
	}
	if q.intervals[2].end != 2 {
		t.Errorf("Expected 2, got %d", q.intervals[2].end)
	}
	if q.intervals[3].end != 1 {
		t.Errorf("Expected 1, got %d", q.intervals[3].end)
	}

	// Test pop
	i := q.Pop()
	if i.end != 4 {
		t.Errorf("Expected 4, got %d", i.end)
	}
	if q.Len() != 3 {
		t.Errorf("Expected 3, got %d", q.Len())
	}
	if q.Peek().end != 3 {
		t.Errorf("Expected 3, got %d", q.Peek().end)
	}

	// Test remove
	q.Remove(1)

	if q.Len() != 2 {
		t.Errorf("Expected 2, got %d", q.Len())
	}
	if q.Peek().end != 3 {
		t.Errorf("Expected 3, got %d", q.Peek().end)
	}

	j := q.Pop()
	if j.end != 3 {
		t.Errorf("Expected 3, got %d", j.end)
	}

	if q.Peek().end != 1 {
		t.Errorf("Expected 1, got %d", q.Peek().end)
	}

}
