package shanks

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"shanks-algorithm/group"
	elliptic_curve "shanks-algorithm/group/elliptic-curve"
	prime_modulus "shanks-algorithm/group/prime-modulus"
	"testing"
)

type testCase struct {
	g         group.Group
	a, b      group.GroupElement
	expectedX []int64
	modulus   int64
}

func makePrimeModulusTestCase(modulus, aVal, bVal int64, expectedX []int64) *testCase {
	g := prime_modulus.NewPrimeModulusGroup(modulus)
	return &testCase{
		g:         g,
		a:         g.GetElementWithValue(aVal),
		b:         g.GetElementWithValue(bVal),
		expectedX: expectedX,
		modulus:   modulus,
	}
}

func makeEllipticCurveTestCase(modulus, b, order, xBase, yBase, xA, yA, xB, yB int64, expectedX []int64) *testCase {
	g := elliptic_curve.NewEllipticCurve(modulus, b, xBase, yBase, order)
	return &testCase{
		g:         g,
		a:         g.GetElementWithValues(xA, yA),
		b:         g.GetElementWithValues(xB, yB),
		expectedX: expectedX,
		modulus:   modulus,
	}
}

var primeModulusCases = []*testCase{
	makePrimeModulusTestCase(29, 2, 4, []int64{2}),
	makePrimeModulusTestCase(29, 6, 16, []int64{10, 24}),
	makePrimeModulusTestCase(29, 21, 8, []int64{15}),
	makePrimeModulusTestCase(29, 5, 1, []int64{0}),
	makePrimeModulusTestCase(29, 6, 1, []int64{0}),
	makePrimeModulusTestCase(29, 7, 1, []int64{0}),
	makePrimeModulusTestCase(6779, 2, 3744, []int64{1500}),
	makePrimeModulusTestCase(90127, 2, 61800, []int64{1500, 46563}),
	makePrimeModulusTestCase(397217, 2, 123093, []int64{34971, 233579}),
	makePrimeModulusTestCase(8325089, 2, 6158204, []int64{894575, 5057119}),
	makePrimeModulusTestCase(46274771, 2, 31207761, []int64{8357533}),
}

var ellipticCurveCases = []*testCase{
	makeEllipticCurveTestCase(29, 25, 31, 6, 7, 22, 14, 21, 28, []int64{8}),
	makeEllipticCurveTestCase(6779, 5528, 6938, 3941, 524, 1229, 3764, 999, 1394, []int64{5680}),
	makeEllipticCurveTestCase(90127, 12724, 90122, 43868, 42525, 73712, 53230, 67169, 9950, []int64{17405, 62466}),
	makeEllipticCurveTestCase(397217, 205386, 396481, 204918, 23673, 55853, 109517, 133119, 77136, []int64{74948}),
	makeEllipticCurveTestCase(8325089, 3607923, 8322546, 2680198, 8252287, 7835166, 1153924, 3265095, 3272592, []int64{1076127}),
	makeEllipticCurveTestCase(46274771, 27159214, 46268375, 28476117, 23486067, 22041538, 37886315, 36855233, 30436882, []int64{27295534, 42722158}),
}

func TestShanksAlgorithm_Execute(t *testing.T) {
	cases := make([]*testCase, 0, len(primeModulusCases)+len(ellipticCurveCases))
	cases = append(cases, primeModulusCases...)
	cases = append(cases, ellipticCurveCases...)

	for _, tc := range cases {
		t.Logf("modulus=%v", tc.modulus)
		algo := NewShanksAlgorithm(tc.a, tc.b, tc.g)
		x, err := algo.Execute()
		assert.NoError(t, err, "shanks algorithm ended with error")
		assert.Contains(t, tc.expectedX, x, "solved wrong")
	}
}

func TestShanksAlgorithm_ExecuteParallel(t *testing.T) {
	cases := make([]*testCase, 0, len(primeModulusCases)+len(ellipticCurveCases))
	cases = append(cases, primeModulusCases...)
	cases = append(cases, ellipticCurveCases...)

	for _, tc := range cases {
		algo := NewShanksAlgorithm(tc.a, tc.b, tc.g)
		x, err := algo.ExecuteParallel(8)
		//t.Logf("x = %d, err= %v", x, err)
		assert.NoError(t, err, "shanks algorithm ended with error")
		assert.Contains(t, tc.expectedX, x, "solved wrong")
	}
}

type numOfThreads struct {
	num       int
	numString string
}

func BenchmarkShanksPrimeModulus(b *testing.B) {
	for _, modulus := range []int64{29, 809, 6779, 90127, 397217, 8325089, 46274771} {
		benchmarkShanksPrimeModulus(b, modulus)
	}
}

func benchmarkShanksPrimeModulus(b *testing.B, modulus int64) {
	b.ResetTimer()
	for _, numOfThreads := range []numOfThreads{
		{1, "One"},
		{2, "Two"},
		{4, "Four"},
		{8, "Eight"},
	} {
		b.Run(fmt.Sprintf("Shanks%sThreads_%d", numOfThreads.numString, modulus), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tc := makePrimeModulusTestCase(modulus, 2, 1500, []int64{1500})
				algo := NewShanksAlgorithm(tc.a, tc.b, tc.g)
				algo.ExecuteParallel(numOfThreads.num)
			}
		})
	}
}

func BenchmarkShanksEllipticCurve(b *testing.B) {
	for _, tc := range ellipticCurveCases {
		benchmarkShanksEllipticCurve(b, tc)
	}
}

func benchmarkShanksEllipticCurve(b *testing.B, tc *testCase) {
	b.ResetTimer()
	for _, numOfThreads := range []numOfThreads{
		{1, "One"},
		{2, "Two"},
		{4, "Four"},
		{8, "Eight"},
	} {
		b.Run(fmt.Sprintf("Shanks%sThreads_%d", numOfThreads.numString, tc.modulus), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				algo := NewShanksAlgorithm(tc.a, tc.b, tc.g)
				algo.ExecuteParallel(numOfThreads.num)
			}
		})
	}
}
