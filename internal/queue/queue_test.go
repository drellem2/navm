package queue

import (
	"testing"
)

func init() {
}

func TestQueue(t *testing.T) {
	q := Queue{}
	if !q.Empty() {
		t.Errorf("Expected true, got false")
	}
	q.Push(1)
	q.Push(2)
	q.Push(3)
	q.Push(4)
	if q.Len() != 4 {
		t.Errorf("Expected 4, got %d", q.
			Len())
	}
	if q.Pop() != 1 {
		t.Errorf("Expected 1, got %d", q.Pop())
	}
	if q.Pop() != 2 {
		t.Errorf("Expected 2, got %d", q.Pop())
	}
	if q.Empty() {
		t.Errorf("Expected false, got true")
	}
	if q.Peek() != 3 {
		t.Errorf("Expected 3, got %d", q.Peek())
	}
}
