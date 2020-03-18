package utils

import (
	"fmt"
	"testing"
)

type IntTestData struct {
	min         int64
	max         int64
	expectation int64
}

var Intdata = []IntTestData{
	{-10, 10, -10},
	{-100, -10, -70},
	{-5, 5, -5},
	{-99999, 99999, -53655},
}

func TestRandInt64Panics(t *testing.T) {
	assertPanic(t, func() {
		RandInt64(1, -1, 1)
	})
}

func TestRandInt64WithInts(t *testing.T) {
	fmt.Println("Testing RandInt64 - Int")

	seed := int64(5)

	for _, test := range Intdata {
		if result := RandInt64(test.min, test.max, seed); result != test.expectation {
			t.Fatal("Expected ", test.expectation, ", got ", result)
		}
	}
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}
