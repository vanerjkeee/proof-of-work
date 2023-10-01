package main

import "math/rand"

var Quotes = []string{"Qoote 1", "Quote 2", "Quote 3", "Quote 4", "Quote 5"}

func GetQuote() string {
	return Quotes[rand.Intn(len(Quotes))]
}
