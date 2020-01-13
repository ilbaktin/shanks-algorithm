package elliptic_curve

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewElipticCurve(t *testing.T) {
	ec, err := NewRandomEllipticCurveGroup(29)
	t.Logf("%v", err)
	assert.NoError(t, err, "error while initializing curve")
	t.Logf("%#v", ec.ec)
}

//func TestElipticCurveOrder(t *testing.T) {
//	xBase := big.NewInt(43868)
//	yBase := big.NewInt(42525)
//	ecParams := &elliptic.CurveParams{
//		P:       big.NewInt(90127),
//		N:       big.NewInt(90122),
//		B:       big.NewInt(12724),
//		Gx:      xBase,
//		Gy:      yBase,
//		BitSize: 5,
//		Name:    "random curve",
//	}
//
//	//t.Logf("%v", big.NewInt(266).Bytes())
//	t.Logf("(%v, %v) is on curve? %v", xBase, yBase, ecParams.IsOnCurve(xBase, yBase))
//	for i := big.NewInt(1); i.Cmp(big.NewInt(100)) == -1; i.Add(i, big.NewInt(1)) {
//		x, y := ecParams.ScalarMult(xBase, yBase, i.Bytes())
//		t.Logf("%v * (%v, %v) = (%v, %v)", i, xBase, yBase, x, y)
//		if x.Cmp(big.NewInt(0)) == 0 && y.Cmp(big.NewInt(0)) == 0 {
//			t.Logf("order = %v", i)
//			break
//		}
//	}
//	xA := big.NewInt(73712)
//	yA := big.NewInt(53230)
//	deg := big.NewInt(62466)
//	xRes, yRes := ecParams.ScalarMult(xA, yA, deg.Bytes())
//	t.Logf("Check if algorithm works properly: %v * (%v, %v) = (%v, %v)", deg, xA, yA, xRes, yRes)
//}

//makeEllipticCurveTestCase(46274771, 27159214, 46268375, 28476117, 23486067, 22041538, 37886315, 36855233, 30436882, []int64{27295534}),
//makeEllipticCurveTestCase(90127, 12724, 90122, 43868, 42525, 24531, 38095, 9879, 16693, []int64{64999}),
//(73712,53230), b=(67169,9950)
