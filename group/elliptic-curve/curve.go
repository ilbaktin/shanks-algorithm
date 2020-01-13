package elliptic_curve

import (
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/pkg/errors"
	"math/big"
	mrand "math/rand"
	"shanks-algorithm/group"
	"time"
)

type EllipticCurveGroup struct {
	ec elliptic.Curve
}

func NewEllipticCurve(modulus, b, xBase, yBase, orderBase int64) *EllipticCurveGroup {
	ec := &elliptic.CurveParams{Name: "specific curve"}
	ec.P = big.NewInt(modulus)
	ec.B = big.NewInt(b)
	ec.Gx = big.NewInt(xBase)
	ec.Gy = big.NewInt(yBase)
	ec.N = big.NewInt(orderBase)
	ec.BitSize = big.NewInt(modulus).BitLen()

	return &EllipticCurveGroup{ec: ec}
}

func NewRandomEllipticCurveGroup(modulusInt64 int64) (*EllipticCurveGroup, error) {
	mrand.Seed(time.Now().UTC().UnixNano())
	modulus := big.NewInt(modulusInt64)

	ec := &elliptic.CurveParams{Name: "random curve"}
	ec.P = modulus
	ec.BitSize = modulus.BitLen()
	b, err := rand.Int(rand.Reader, modulus)
	if err != nil {
		return nil, err
	}
	ec.B = b

	fmt.Printf("P=%v\n", modulus)
	fmt.Printf("B=%v\n", b)
	fmt.Println(isCurveSuitable(big.NewInt(-3), b, modulus))

	var x, y, basePointOrder *big.Int

TOP:
	for {
		//fmt.Println()
		x, err = rand.Int(rand.Reader, modulus)
		if err != nil {
			return nil, err
		}
		//fmt.Printf("x=%v\n", x)

		ySqr := new(big.Int).Mul(x, x)
		//fmt.Printf("x^2=%v\n", ySqr)

		ySqr.Mul(ySqr, x)
		//fmt.Printf("x^3=%v\n", ySqr)

		ySqr.Sub(ySqr, new(big.Int).Mul(big.NewInt(3), x))
		//fmt.Printf("x^3 - 3x=%v\n", ySqr)

		ySqr.Add(ySqr, b)
		//fmt.Printf("x^3 - 3x + B=%v\n", ySqr)

		ySqr.Mod(ySqr, modulus)
		//fmt.Printf("x^3 - 3x + B mod %v=%v\n", modulus, ySqr)

		if big.Jacobi(ySqr, modulus) != 1 {
			//fmt.Printf("y^2=%v is not jacobi\n", ySqr)
			continue
		}

		yPositive := ySqr.ModSqrt(ySqr, modulus)
		//fmt.Printf("checking if x=%v,y=%v on curve\n", x, yPositive)
		yNegative := new(big.Int).Sub(modulus, yPositive)
		//fmt.Printf("checking if x=%v,y=%v on curve\n", x, yNegative)
		if ec.IsOnCurve(x, yPositive) {
			y = yPositive
		} else if ec.IsOnCurve(x, yPositive) {
			y = yNegative
		} else {
			continue
		}
		fmt.Printf("(%v, %v) on curve\n", x, y)

		order := big.NewInt(1)
		xTmp := new(big.Int).Set(x)
		yTmp := new(big.Int).Set(y)
		for {
			xTmp, yTmp = ec.Add(xTmp, yTmp, x, y)
			if xTmp.Cmp(x) == 0 && yTmp.Cmp(y) == 0 {
				//fmt.Printf("order of (%v,%v) is %v\n", x, y, order)
				if true { // order.ProbablyPrime(64) && order.Cmp(modulus) == 1
					basePointOrder = order
					//fmt.Printf("order %v is probably prime\n", order)
					break TOP
				} else {
					//fmt.Printf("order %v is not prime, continue seaching\n", order)
					continue TOP
				}
			}
			order.Add(order, big.NewInt(1))
		}
	}
	ec.Gx = x
	ec.Gy = y
	ec.N = basePointOrder

	return &EllipticCurveGroup{
		ec: ec,
	}, nil
}

func (ecg *EllipticCurveGroup) GetElementWithValues(x, y int64) group.GroupElement {
	return &EllipticCurveGroupElement{
		x:   big.NewInt(x),
		y:   big.NewInt(y),
		ecg: ecg,
	}
}

func isCurveSuitable(a, b, p *big.Int) bool {
	left := new(big.Int).Mul(big.NewInt(4), a)
	left.Mul(left, a)
	left.Mul(left, a)
	right := new(big.Int).Mul(big.NewInt(27), b)
	right.Mul(right, b)

	left.Add(left, right)
	res := left.Mod(left, p)
	fmt.Println(res)
	return res.Cmp(big.NewInt(0)) != 0

}

func (ecg *EllipticCurveGroup) GroupOrder() int64 {
	return ecg.ec.Params().N.Int64()
}

func (ecg *EllipticCurveGroup) String() string {
	return fmt.Sprintf("%#v", ecg.ec.Params())
}

func (ecg *EllipticCurveGroup) RandomElement() (*EllipticCurveGroupElement, error) {
	mrand.Seed(time.Now().UTC().UnixNano())
	order, err := rand.Int(rand.Reader, ecg.ec.Params().N)
	if err != nil {
		return nil, errors.Wrapf(err, "gen rand big int")
	}
	x, y := ecg.ec.ScalarBaseMult(order.Bytes())
	return &EllipticCurveGroupElement{
		x:   x,
		y:   y,
		ecg: ecg,
	}, nil
}

type EllipticCurveGroupElement struct {
	x, y *big.Int
	ecg  *EllipticCurveGroup
}

func (ecge *EllipticCurveGroupElement) Add(other group.GroupElement) group.GroupElement {
	otherElem := other.(*EllipticCurveGroupElement)
	x, y := ecge.ecg.ec.Add(ecge.x, ecge.y, otherElem.x, otherElem.y)
	return &EllipticCurveGroupElement{
		x:   x,
		y:   y,
		ecg: ecge.ecg,
	}
}

func (ecge *EllipticCurveGroupElement) Sub(other group.GroupElement) group.GroupElement {
	otherElem := other.(*EllipticCurveGroupElement)
	x, y := ecge.ecg.ec.Add(ecge.x, ecge.y, otherElem.x, new(big.Int).Sub(ecge.ecg.ec.Params().P, ecge.y))
	return &EllipticCurveGroupElement{
		x:   x,
		y:   y,
		ecg: ecge.ecg,
	}
}

func (ecge *EllipticCurveGroupElement) Pow(k int64) group.GroupElement {
	bigK := big.NewInt(k)
	x, y := ecge.ecg.ec.ScalarMult(ecge.x, ecge.y, bigK.Bytes())
	return &EllipticCurveGroupElement{
		x:   x,
		y:   y,
		ecg: ecge.ecg,
	}
}

func (ecge *EllipticCurveGroupElement) IsNeutralElement() bool {
	return ecge.x == ecge.ecg.ec.Params().Gx && ecge.y == ecge.ecg.ec.Params().Gy
}

func (ecge *EllipticCurveGroupElement) Hash() interface{} {
	return fmt.Sprintf("%s_%s", ecge.x.String(), ecge.y.String())
}

func (ecge *EllipticCurveGroupElement) Equal(other group.GroupElement) bool {
	otherElem := other.(*EllipticCurveGroupElement)
	return ecge.x.Cmp(otherElem.x) == 0 && ecge.y.Cmp(otherElem.y) == 0
}

func (ecge *EllipticCurveGroupElement) String() string {
	return fmt.Sprintf("(%v,%v)", ecge.x, ecge.y)
}
