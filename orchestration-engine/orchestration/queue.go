package orchestration

import (
	"errors"
	"sync"
)

// * thread-safe task queue
type TaskQueue struct {
	queue chan *Task
	mu    sync.Mutex
	size  int
}

func NewTaskQueue(capacity int) *TaskQueue {
	return &TaskQueue{
		queue: make(chan *Task, capacity),
		size:  0,
	}
}

func (q *TaskQueue) Enqueue(task *Task) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	select {
	case q.queue <- task:
		q.size++
		return nil
	default:
		return errors.New("task queue is full")
	}
}

// Dequeue removes and returns the next task from the queue
func (q *TaskQueue) Dequeue() (*Task, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.size == 0 {
		return nil, errors.New("queue is empty")
	}

	task := <-q.queue
	q.size--
	return task, nil
}

// Size returns the current size of the task queue
func (q *TaskQueue) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.size
}

func (q *TaskQueue) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.size == 0
}

func (q *TaskQueue) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Drain the channel by reading all items
	for len(q.queue) > 0 {
		<-q.queue
	}
	q.size = 0
}
