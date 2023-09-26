package pubsub

import (
	"sync"
)

type Subscriber[T any] chan T

type PubSub[T any] struct {
	subMap map[Subscriber[T]]Subscriber[T]
	mux    sync.Mutex
}

func New[T any]() *PubSub[T] {
	return &PubSub[T]{mux: sync.Mutex{}, subMap: map[Subscriber[T]]Subscriber[T]{}}
}

func (p *PubSub[T]) Subscribe() *Subscriber[T] {
	p.mux.Lock()
	defer p.mux.Unlock()
	sub := make(Subscriber[T], 1)
	p.subMap[sub] = sub

	return &sub
}

func (p *PubSub[T]) Publish(msg T) {
	p.mux.Lock()
	defer p.mux.Unlock()

	for _, sub := range p.subMap {
		sub <- msg
	}
}

func (p *PubSub[T]) Close(sub *Subscriber[T]) {
	p.mux.Lock()
	defer p.mux.Unlock()

	target := p.subMap[*sub]
	if target != nil {
		close(target)
		delete(p.subMap, target)
	}
}
