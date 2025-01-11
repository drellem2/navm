package navm

import (
	"testing"
)

func init() {
}

func TestLivenessQueue(t *testing.T) {
	// Push a few values and check that they are sorted
	q := LivenessQueue{}
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
