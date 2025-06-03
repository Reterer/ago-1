package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	tasks    []Task
	heapSize int
}

func NewScheduler() Scheduler {
	return Scheduler{}
}

func (s *Scheduler) grow() {
	if len(s.tasks) == 0 {
		s.tasks = make([]Task, 1)
		return
	}

	s.tasks = append(s.tasks, make([]Task, len(s.tasks))...)
}

func (s *Scheduler) AddTask(task Task) {
	if len(s.tasks) == s.heapSize {
		s.grow()
	}

	s.tasks[s.heapSize] = task
	pos := s.heapSize
	s.heapSize++

	s.fixHeapUp(pos)
}

func (s *Scheduler) fixHeapUp(i int) {
	parent := (i - 1) / 2
	for i > 0 && s.tasks[parent].Priority < s.tasks[i].Priority {
		s.tasks[i], s.tasks[parent] = s.tasks[parent], s.tasks[i]

		i = parent
		parent = (i - 1) / 2
	}
}

func (s *Scheduler) fixHeapDown(i int) {
	for {
		left := i*2 + 1
		right := i*2 + 2
		largest := i

		if left < s.heapSize && s.tasks[left].Priority > s.tasks[largest].Priority {
			largest = left
		}

		if right < s.heapSize && s.tasks[right].Priority > s.tasks[largest].Priority {
			largest = right
		}

		if largest == i {
			break
		}

		s.tasks[i], s.tasks[largest] = s.tasks[largest], s.tasks[i]
		i = largest
	}
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	var pos int
	for i := range s.tasks {
		if s.tasks[i].Identifier == taskID {
			pos = i
			break
		}
	}

	oldPriority := s.tasks[pos].Priority
	s.tasks[pos].Priority = newPriority

	if oldPriority > newPriority {
		s.fixHeapDown(pos)
	} else {
		s.fixHeapUp(pos)
	}
}

func (s *Scheduler) GetTask() Task {
	res := s.tasks[0]

	s.heapSize--
	s.tasks[0] = s.tasks[s.heapSize]
	s.fixHeapDown(0)

	return res
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)
	task1.Priority = 100 // Так как повысили приоритет задаче

	task = scheduler.GetTask()
	assert.Equal(t, task1, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}
