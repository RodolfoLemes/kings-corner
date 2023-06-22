package pubsub

import (
	"errors"
	"fmt"
)

type PubSub[T any] interface {
	Publish(topic string, msg T) error
	Subscribe(topic string) <-chan T
	Close()
}

func New[T any]() PubSub[T] {
	return newInMemory[T]()
}

var Closed error = errors.New("closed")

type PubSubError struct {
	err error
}

func (pse *PubSubError) Error() string {
	return fmt.Sprintf("pubsub err: %s", pse.err.Error())
}
