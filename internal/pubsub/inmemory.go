package pubsub

import "sync"

type inMemory[T any] struct {
	mu       sync.Mutex
	topics   map[string][]chan T
	isClosed bool
}

func newInMemory[T any]() *inMemory[T] {
	return &inMemory[T]{
		mu:     sync.Mutex{},
		topics: make(map[string][]chan T),
	}
}

func (im *inMemory[T]) Publish(topic string, msg T) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if im.isClosed {
		return &PubSubError{Closed}
	}

	for _, ch := range im.topics[topic] {
		ch <- msg
	}

	return nil
}

func (im *inMemory[T]) Subscribe(topic string) <-chan T {
	im.mu.Lock()
	defer im.mu.Unlock()

	if im.isClosed {
		return nil
	}

	ch := make(chan T)
	im.topics[topic] = append(im.topics[topic], ch)

	return ch
}

func (im *inMemory[T]) Close() {
	im.mu.Lock()
	defer im.mu.Unlock()

	if im.isClosed {
		return
	}

	im.isClosed = true

	for _, chs := range im.topics {
		for _, ch := range chs {
			close(ch)
		}
	}
}
