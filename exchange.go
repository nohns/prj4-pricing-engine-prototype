package exchange

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrTickerNotFound = errors.New("ticker not found")
)

type Exchange struct {
	tickers map[string]*Item
}

func New() *Exchange {
	return &Exchange{
		tickers: make(map[string]*Item),
	}
}

func (e *Exchange) Add(ticker string, item Item) {
	e.tickers[ticker] = &item
}

func (e *Exchange) Buy(ticker string) error {
	i := e.tickers[ticker]
	if i == nil {
		return ErrTickerNotFound
	}

	i.Buy()
	return nil
}

func (e *Exchange) String() string {
	sb := strings.Builder{}

	sb.WriteString("--- Exchange prices ---\n")
	for t, i := range e.tickers {
		fmt.Fprintf(&sb, " - %s: %.2f DKK\n", t, i.CurrPrice())
	}
	sb.WriteString("-----------------------\n")

	return sb.String()
}
