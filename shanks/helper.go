package shanks

type int64Range struct {
	start, end int64
}

func splitRange(rangeSize, batchNumber int64) []int64Range {
	if batchNumber <= 1 {
		return []int64Range{
			{
				start: 0,
				end:   rangeSize,
			},
		}
	}
	if batchNumber > rangeSize {
		batchNumber = rangeSize
	}
	batchSize := rangeSize / batchNumber
	reminder := rangeSize % batchNumber
	if reminder > 0 {
		batchNumber = batchNumber - 1
		batchSize = rangeSize / batchNumber
		reminder = rangeSize % batchNumber
	}

	ranges := make([]int64Range, 0, batchNumber+1)
	for i := int64(0); i < batchNumber; i++ {
		ranges = append(ranges, int64Range{
			start: i * batchSize,
			end:   (i + 1) * batchSize,
		})
	}
	if reminder > 0 {
		ranges = append(ranges, int64Range{
			start: batchNumber * batchSize,
			end:   rangeSize,
		})
	}
	return ranges
}
