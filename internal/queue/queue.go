package queue

import (
	"strconv"
)

// Convenience struct to create queues of integers

type Queue struct {
	data []int
}

func (q *Queue) Push(i int) {
	q.data = append(q.data, i)
}

func (q *Queue) Pop() int {
	if len(q.data) == 0 {
		panic("Empty queue")
	}
	i := q.data[0]
	q.data = q.data[1:]
	return i
}

func (q *Queue) Len() int {
	return len(q.data)
}

func (q *Queue) Empty() bool {
	return len(q.data) == 0
}

func (q *Queue) Peek() int {
	if len(q.data) == 0 {
		panic("Empty queue")
	}
	return q.data[0]
}

func (q *Queue) Print() string {
	var s string
	for _, i := range q.data {
		s += " " + strconv.Itoa(i)
	}
	return s
}
