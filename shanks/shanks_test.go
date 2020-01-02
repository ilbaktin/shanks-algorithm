package shanks

import (
	"github.com/stretchr/testify/assert"
	"shanks-algorithm/group"
	prime_modulus "shanks-algorithm/group/prime-modulus"
	"testing"
)

type testCase struct {
	g         group.Group
	a, b      group.GroupElement
	expectedX int64
}

func makePrimeModulusTestCase(modulus, aVal, bVal, expectedX int64) *testCase {
	g := prime_modulus.NewPrimeModulusGroup(modulus)
	return &testCase{
		g:         g,
		a:         g.GetElementWithValue(aVal),
		b:         g.GetElementWithValue(bVal),
		expectedX: expectedX,
	}
}

func TestShanksAlgorithm_Execute(t *testing.T) {
	cases := []*testCase{
		makePrimeModulusTestCase(29, 2, 4, 2),
		makePrimeModulusTestCase(29, 6, 16, 10),
		makePrimeModulusTestCase(29, 21, 8, 15),
		makePrimeModulusTestCase(29, 5, 1, 0),
		makePrimeModulusTestCase(29, 6, 1, 0),
		makePrimeModulusTestCase(29, 7, 1, 0),
	}
	for _, tc := range cases {
		algo := NewShanksAlgorithm(tc.a, tc.b, tc.g)
		x, err := algo.Execute()
		assert.NoError(t, err, "shanks algorithm ended with error")
		assert.Equal(t, tc.expectedX, x, "solved wrong: got x=%d, want x=%d", x, tc.expectedX)
	}
}
