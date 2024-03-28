package main

import (
	"exchange-algo-prototype"
	"fmt"
	"time"
)

const tickerBlueWater = "BLÅ VAND"
const tickerNnguaq = "NNGUAQ"

func main() {

	// Init exchange
	ex := exchange.New()
	ex.Add(tickerBlueWater, exchange.NewItem(25, 1.1, 20*time.Second))
	ex.Add(tickerNnguaq, exchange.NewItem(35, 1.05, 30*time.Second))

	// print index before buying anything
	fmt.Println(ex)

	// Buy a single blue water + print index
	ex.Buy(tickerBlueWater)

	fmt.Printf("PLACED BUY on %s\n\n", tickerBlueWater)

	fmt.Println(ex)

	// Wait three second and print prices
	fmt.Printf("WAITING 3 secs...\n\n")
	<-time.After(3 * time.Second)
	fmt.Println(ex)

	// Buy two blue water in an instant + print index
	ex.Buy(tickerBlueWater)
	fmt.Printf("PLACED BUY on %s\n", tickerBlueWater)
	ex.Buy(tickerBlueWater)
	fmt.Printf("PLACED BUY on %s\n\n", tickerBlueWater)
	ex.Buy(tickerNnguaq)
	fmt.Printf("PLACED BUY on %s\n\n", tickerNnguaq)
	fmt.Println(ex)

	// Wait three second and print prices
	fmt.Printf("WAITING 3 secs...\n\n")
	<-time.After(3 * time.Second)
	fmt.Println(ex)

	fmt.Printf("WAITING 10 secs...\n\n")
	<-time.After(10 * time.Second)
	fmt.Println(ex)
}
