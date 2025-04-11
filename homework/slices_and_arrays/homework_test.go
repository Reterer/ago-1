package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Number interface {
	int | int8 | int16 | int32 | int64
}

type CircularQueue[T Number] struct {
	values []T
	first  int // первый элемент в очереди
	end    int // следующий свободный элемент
}

func NewCircularQueue[T Number](size int) CircularQueue[T] {
	return CircularQueue[T]{
		values: make([]T, size),
		first:  -1, // Если в циклической очереди нет элементов, то first принимает значение -1
	}
}

func (q *CircularQueue[T]) Push(value T) bool {
	if q.Full() {
		return false
	}

	// Если в очереди нет элементов, то начинаем заполнять её сначала
	if q.first == -1 {
		q.first = 0
		q.end = 0
	}

	q.values[q.end] = value

	q.end = q.IncIdx(q.end)

	return true
}

func (q *CircularQueue[T]) Pop() bool {
	if q.Empty() {
		return false
	}

	q.first = q.IncIdx(q.first)

	// Очередь стала пустой
	if q.first == q.end {
		q.first = -1
	}

	return true
}

func (q *CircularQueue[T]) Front() T {
	if q.Empty() {
		return -1
	}

	return q.values[q.first]
}

func (q *CircularQueue[T]) Back() T {
	if q.Empty() {
		return -1
	}

	back := q.DecIdx(q.end)
	return q.values[back]
}

func (q *CircularQueue[T]) Empty() bool {
	return q.first == -1
}

func (q *CircularQueue[T]) Full() bool {
	return q.first == q.end
}

func (q *CircularQueue[T]) IncIdx(idx int) int {
	shifted := idx + 1
	if shifted == len(q.values) {
		return 0
	}
	return shifted
}

func (q *CircularQueue[T]) DecIdx(idx int) int {
	shifted := idx - 1
	if shifted < 0 {
		return len(q.values) - 1
	}
	return shifted
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[int](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}
