package utils_test

import (
	"fantasymarket/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

type includesTestData struct {
	slice      []string
	someData   []string
	inclData   string
	someExpect bool
	inclExpect bool
}

var includesData = []includesTestData{
	{
		slice:      []string{"one", "two", "three", "four", "five"},
		someData:   []string{"one", "five"},
		inclData:   "one",
		someExpect: true,
		inclExpect: true,
	}, {
		slice:      []string{"one", "two", "three", "four", "five"},
		someData:   []string{"two", "six"},
		inclData:   "six",
		someExpect: true,
		inclExpect: false,
	},
	{
		slice:      []string{"one", "two", "three", "four", "five"},
		someData:   []string{"six", "seven"},
		inclData:   "six",
		someExpect: false,
		inclExpect: false,
	},
	{
		slice:      []string{"one", "two", "three", "four", "five"},
		someData:   []string{"", " "},
		inclData:   "",
		someExpect: false,
		inclExpect: false,
	},
	{
		slice:      []string{},
		someData:   []string{"", "six"},
		inclData:   "",
		someExpect: false,
		inclExpect: false,
	},
}

func testSome(t *testing.T) {
	for _, test := range includesData {
		assert.Equal(t, test.someExpect, utils.Some(test.slice, test.someData))
		assert.Equal(t, test.inclExpect, utils.Includes(test.slice, test.inclData))
	}
}
