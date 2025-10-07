package data

import "github.com/Clayal10/enders_game/lib/cross"

type Queue[T any] struct {
	q       []T
	r, w, l int
}

func NewQueue[T any](length int) *Queue[T] {
	return &Queue[T]{
		q: make([]T, length),
		l: length,
	}
}

func (q *Queue[T]) Enqueue(obj ...T) {
	for _, o := range obj {
		q.q[q.w%q.l] = o
		q.w++
	}
}

func (q *Queue[T]) Dequeue() (result T, err error) {
	if q.r == q.w {
		return result, cross.ErrQueueEmpty
	}
	result = q.q[q.r%q.l]
	q.r++
	return
}

func (q *Queue[T]) IsEmpty() bool {
	return q.r == q.w
}
