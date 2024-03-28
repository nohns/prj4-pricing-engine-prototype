package app

type Test interface {
    PlaceOrder(bid int, qty int) error
}
