package engine

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrItemAlreadyTracked = errors.New("item already tracked")
	ErrItemNotFound       = errors.New("item not found")
	ErrQuantityBelowOne   = errors.New("quantity below one")
	ErrItemParamsInvalid  = errors.New("item params invalid")
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
	MaxPrice float64
	MinPrice float64

	// StartPrice is the price of the item at the LastOrdered time.
	StartPrice float64

	// BuyMultiplier is the multiplier that decides how much the price
	// of an item increases when exactly one is ordered.
	BuyMultiplier float64

	// HalfTime specifies the amount of time before a price reaches half
	// its orignal price, assuming no orders placed.
	HalfTime int

	LastOrdered time.Time
}

func (ip *ItemParams) validate() error {
    if ip.MaxPrice < ip.MinPrice {
        return fmt.Errorf("%w: max price must be larger than min price", ErrItemParamsInvalid)
    }
    if ip.BuyMultiplier < 1 {
        return fmt.Errorf("%w: buy multiplier must be larger than 1", ErrItemParamsInvalid)
    }
    if ip.HalfTime < 1 {
        return fmt.Errorf("%w: half time must be greater than 1 second", ErrItemParamsInvalid)
    }
    if ip.MaxPrice < ip.StartPrice || ip.MinPrice > ip.StartPrice {
        return fmt.Errorf("%w: start price must be between min and max price", ErrItemParamsInvalid)
    }
    return nil
}

// Track item registers an item and proceeds to track its price, when the engine starts.
func (e *Engine) TrackItem(id string, params ItemParams) error {
    if err := params.validate(); err != nil {
        return err
    }
    
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.actors[id]; ok {
		return ErrItemAlreadyTracked
	}
	e.actors[id] = newActor(id, params, e.updates)

	return nil
}

// Order item takes the id of a given item and the quantity ordered and proceeds to
// increase the price.
func (e *Engine) OrderItem(id string, qty int) error {
	if qty < 1 {
		return ErrQuantityBelowOne
	}
	a, err := e.actor(id)
	if err != nil {
		return err
	}

	a.order(qty)
	return nil
}

// Tweak an items param while engine is running.
func (e *Engine) TweakItem(id string, newparams ItemParams) error {
    if err := newparams.validate(); err != nil {
        return err
    }
    
	a, err := e.actor(id)
	if err != nil {
		return err
	}

	a.updateParams(newparams)
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
