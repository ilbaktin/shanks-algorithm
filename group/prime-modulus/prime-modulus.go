package prime_modulus

import (
	"crypto/sha256"
	"hash"
	"shanks-algorithm/group"
)

type PrimeModulusGroup struct {
	modulus	int64
	hashFunc hash.Hash
}

func NewPrimeModulusGroup(modulus int64) *PrimeModulusGroup {
	return &PrimeModulusGroup{
		modulus: modulus,
		hashFunc: sha256.New(),
	}
}

func (g *PrimeModulusGroup) GetElementWithValue(value int64) *PrimeModulusElement {
	return &PrimeModulusElement{
		value: value,
		g: g,
	}
}

func (g *PrimeModulusGroup) GroupOrder() int64 {
	return g.modulus - 1
}

//func (g *PrimeModulusGroup) ElementOrder(element group.GroupElement) int64 {
//	return
//}

type PrimeModulusElement struct {
	value	int64
	g *PrimeModulusGroup
}


func (pme *PrimeModulusElement) Add(other group.GroupElement) group.GroupElement {
	return &PrimeModulusElement{
		value: (pme.value * other.(*PrimeModulusElement).value) % pme.g.modulus,
		g: pme.g,
	}
}

func (pme *PrimeModulusElement) Sub(other group.GroupElement) group.GroupElement {
	return &PrimeModulusElement{
		value: (pme.value - other.(*PrimeModulusElement).value) % pme.g.modulus,
		g: pme.g,
	}
}

func (pme *PrimeModulusElement) Equal(other group.GroupElement) bool {
	return pme.value == other.(*PrimeModulusElement).value
}

func (pme *PrimeModulusElement) IsNeutralElement() bool {
	return pme.value == 1
}

func (pme *PrimeModulusElement) Pow(power int64) group.GroupElement {
	result := int64(1)
	for i := int64(0); i < power; i++ {
		result *= pme.value
		result = result % pme.g.modulus
	}
	return &PrimeModulusElement{
		value: result,
		g: pme.g,
	}
}

func (pme *PrimeModulusElement) Hash() interface{} {
	return pme.value
}

