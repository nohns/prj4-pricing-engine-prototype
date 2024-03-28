package exchange

import (
	"math"
	"time"
)

type Item struct {
	baseprice float64
	a, b      float64
	bmult     float64
	cooldown  time.Duration
	lastbuy   time.Time
}

// Creates a new exchange item, with a base price and the decay cooldown period
func NewItem(price float64, buyMultiplier float64, cooldown time.Duration) Item {
	b := price * buyMultiplier
	y1, y2 := b, price
	x1, x2 := 0, int(cooldown/time.Second) // Convert cooldown to seconds
	a := math.Pow(y2/y1, 1/float64(x2-x1)) // find growth coefficient for decay function

	return Item{
		baseprice: price,
		a:         a,
		b:         b,
		bmult:     buyMultiplier,
		cooldown:  cooldown,
	}
}

// Triggers a buy on the item, raising its price
func (i *Item) Buy() {
	// Reset exponentiel growth when cooldown is reached
	if i.PriceDecayed() {
		i.b = i.baseprice
	}

	i.lastbuy = time.Now()
	i.b *= i.bmult
}

// Calculate the current price, at the time when called
func (i Item) CurrPrice() float64 {
	x := time.Now().Sub(i.lastbuy) / time.Second
	if i.PriceDecayed() {

		return i.baseprice
	}

	// Return decay value as we are still in cooldown
	return i.b * math.Pow(i.a, float64(x))
}

func (i Item) PriceDecayed() bool {
	sincebuy := time.Now().Sub(i.lastbuy)
	return sincebuy > i.cooldown
}
