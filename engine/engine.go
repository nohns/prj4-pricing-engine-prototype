package engine

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrItemAlreadyTracked = errors.New("item already tracked")
	ErrItemNotFound       = errors.New("item not found")
)

const (
	priceUpdateBufSize = 10 // Arbitrary size, which seems like a good fit :)
)

type PriceUpdate struct {
	id    string
	price float64
	at    time.Time
}

type Engine struct {
	// actors keeps track of actors registered. ONLY interact with it when mu lock is acquired
	actors  map[string]*actor
	updates chan PriceUpdate
	mu      sync.RWMutex
}

// New instantiates a pricing engine, which can track items over time.
func New() *Engine {
	return &Engine{
		actors:  make(map[string]*actor),
		updates: make(chan PriceUpdate, priceUpdateBufSize),
	}
}

type ItemParams struct {
	upperb float64
	lowerb float64
}

func (e *Engine) TrackItem(id string, params ItemParams) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.actors[id]; ok {
		return ErrItemAlreadyTracked
	}
	e.actors[id] = newActor(id, params, e.updates)

	return nil
}

func (e *Engine) OrderItem(id string, qty int) error {
	a, err := e.actor(id)
	if err != nil {
		return err
	}

	a.order(qty)
	return nil
}

func (e *Engine) actor(id string) (*actor, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	a, ok := e.actors[id]
	if !ok {
		return nil, ErrItemNotFound
	}

	return a, nil
}

func (e *Engine) ReadUpdate() PriceUpdate {
	return <-e.updates
}

func (e *Engine) Start() {
    e.mu.RLock()
    defer e.mu.RUnlock()

    for _, a := range e.actors {
        a.start()
    }
}
