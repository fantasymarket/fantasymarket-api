package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type IntTestDataRange struct {
	seed        string
	min         int64
	max         int64
	expectation int64
}

var rangeData = []IntTestDataRange{
	{"test", -10, 10, -7},
	{"test2", -100, -10, -68},
	{"test3", -5, 5, -12},
	{"test4", -99999, 99999, -115385},
}

func TestInt64HashRangePanics(t *testing.T) {
	assert.PanicsWithValue(t, "invalid argument to Int64HashRange: min cannot be greater than max", func() {
		Int64HashRange(1, -1, "test")
	})
}

func TestInt64HashnPanics(t *testing.T) {
	assert.PanicsWithValue(t, "invalid argument to Int64Hashn: n can't be less than 0", func() {
		Int64Hashn(-1, "test")
	})
}

func TestInt64HashRange(t *testing.T) {
	for _, test := range rangeData {
		result := Int64HashRange(test.min, test.max, test.seed)
		assert.Equal(t, test.expectation, result)
	}
}
