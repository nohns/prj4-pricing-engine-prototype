package engine

import (
	"math/rand"
	"time"
)

const (
	updateInterval = 30 * time.Second
    // maxStartDelay defines the upper bound for a random start delay, so prices seem to be updating in an arbitrary order
    maxStartDelay = 10 * time.Second
)

type actor struct {
	id  string
	itm item
	// Orders chan receives quantities ordered
	orders chan int
    // params chan receives updates to item parameters
    params chan ItemParams
	// Out chan is for writing price updates back to the engine
	out chan<- PriceUpdate
	// t ticker schedules price updates to be outputted
	t *time.Ticker
}

func newActor(id string, params ItemParams, out chan<- PriceUpdate) *actor {
	return &actor{
		id: id,
		itm: item{
			params: params,
		},
        orders: make(chan int),
		out: out,
	}
}

// Order notifies the actor that an order as been placed
func (a *actor) order(qty int) {
	a.orders <- qty
}


func (a *actor) start() {
	go a.listen()
    go a.primeUpdateScheduler()
}

// prime update scheduler sleeps for a random delay between 0 and maxStartDelay, to make price updates
// seem more random, but still with a set interval between them, so graphs look nicer.
func (a *actor) primeUpdateScheduler() {
    delay := time.Duration(rand.Intn(int(maxStartDelay)))
	time.Sleep(delay)
    a.t = time.NewTicker(updateInterval)
}

func (a *actor) listen() {
	for {
		select {
		case qty := <-a.orders:
			a.handleOrderPlaced(qty)
		case <-a.t.C:
			a.handleTick()
        case params := <-a.params:
            a.handleParamsUpdated(params)
		}
	}
}

// handleOrderPlaced handles incoming order quantities
func (a *actor) handleOrderPlaced(qty int) {
	a.itm.order(qty)
}

// handleTick outputs an update of the current price of the item tracked by the actor
func (a *actor) handleTick() {
	a.out <- PriceUpdate{
		id:    a.id,
		price: a.itm.price(),
		at:    time.Now(),
	}
}

// handleParamsUpdated tweaks the pricing parameters when actor is running
func (a *actor) handleParamsUpdated(params ItemParams) {
    a.itm.tweakParams(params)
}
