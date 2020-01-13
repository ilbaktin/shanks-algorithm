package shanks

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitRange(t *testing.T) {
	testCases := []struct {
		rangeSize, batchNumber int64
		expected               []int64Range
	}{
		{
			10,
			15,
			[]int64Range{
				{0, 1},
				{1, 2},
				{2, 3},
				{3, 4},
				{4, 5},
				{5, 6},
				{6, 7},
				{7, 8},
				{8, 9},
				{9, 10},
			},
		},
		{
			50,
			5,
			[]int64Range{
				{0, 10},
				{10, 20},
				{20, 30},
				{30, 40},
				{40, 50},
			},
		},
		{
			50,
			4,
			[]int64Range{
				{0, 16},
				{16, 32},
				{32, 48},
				{48, 50},
			},
		},
		{
			1000,
			0,
			[]int64Range{
				{0, 1000},
			},
		},
		{
			1000,
			1,
			[]int64Range{
				{0, 1000},
			},
		},
		{
			1000,
			-100,
			[]int64Range{
				{0, 1000},
			},
		},
	}

	for _, tc := range testCases {
		result := splitRange(tc.rangeSize, tc.batchNumber)
		assert.ElementsMatch(t, tc.expected, result, "wrong ranges, want %v, got %v", tc.expected, result)
	}
}
