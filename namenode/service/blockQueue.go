package NNService

import (
	"container/list"
	"sync"
)

type DNBlockQueue struct {
	q  *list.List
	mu sync.Mutex
}

func InitQueue() *DNBlockQueue {
	return &DNBlockQueue{
		q: list.New(),
	}
}

func (queue *DNBlockQueue) EnQueue(data string) {
	if data == "" {
		return
	}

	queue.mu.Lock()
	defer queue.mu.Unlock()
	queue.q.PushBack(data)
}

func (queue *DNBlockQueue) DeQueue() string {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	if element := queue.q.Front(); element != nil {
		queue.q.Remove(element)
		return element.Value.(string)
	}
	return ""
}

func (queue *DNBlockQueue) Clear() {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	for element := queue.q.Front(); element != nil; {
		elementNext := element.Next()
		queue.q.Remove(element)
		element = elementNext
	}
}

func (queue *DNBlockQueue) Size() int64 {
	return int64(queue.q.Len())
}
