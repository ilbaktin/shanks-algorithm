package main

import (
	"fmt"
	"log"
	elliptic_curve "shanks-algorithm/group/elliptic-curve"
	"shanks-algorithm/shanks"
)

func main() {
	//g := prime_modulus.NewPrimeModulusGroup(29)
	//a := g.GetElementWithValue(2)
	//b := g.GetElementWithValue(4)

	//g, _ := elliptic_curve.NewRandomEllipticCurveGroup(90127)
	g := elliptic_curve.NewEllipticCurve(90127, 12724, 43868, 42525, 90122)
	a, _ := g.RandomElement()
	b, _ := g.RandomElement()

	fmt.Printf("Group: %s\n", g)
	fmt.Printf("a=%s, b=%s\n", a, b)

	algo := shanks.NewShanksAlgorithm(a, b, g)
	log.Println("Algorithm initialized")

	x, err := algo.ExecuteParallel(8)
	//x, err := algo.Execute()
	if err != nil {
		log.Fatalf("shanks algorithm ended with error: err=%v", err)
	}
	log.Printf("shanks algorithm succeeded: x=%d", x)
}
