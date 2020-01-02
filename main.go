package main

import (
	"log"
	prime_modulus "shanks-algorithm/group/prime-modulus"
	"shanks-algorithm/shanks"
)

func main() {
	g := prime_modulus.NewPrimeModulusGroup(29)
	log.Println("Group initialized")
	algo := shanks.NewShanksAlgorithm(g.GetElementWithValue(6), g.GetElementWithValue(16), g)
	log.Println("Algorithm initialized")
	x, err := algo.Execute()
	if err != nil {
		log.Fatalf("shanks algorithm ended with error: err=%v", err)
	}
	log.Printf("shanks algorithm succeeded: x=%d", x)
}
