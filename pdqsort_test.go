package main

import (
	"math/rand"
	"sort"
	"testing"
)

func TestSortFloat64SlicePdq(t *testing.T) {
	data := float64s
	a := Float64Slice(data[0:])
	PdqSort(a)
	if !IsSorted(a) {
		t.Errorf("sorted %v", float64s)
		t.Errorf("   got %v", data)
	}
}

func TestHeapsortBMPdq(t *testing.T) {
	testBentleyMcIlroy(t, PdqSort, func(n int) int { return n * lg(n) * 12 / 10 })
}

func TestReverseSortIntSlicePdq(t *testing.T) {
	data := ints
	data1 := ints
	a := IntSlice(data[0:])
	PdqSort(a)
	r := IntSlice(data1[0:])
	PdqSort(sort.Reverse(r))
	for i := 0; i < len(data); i++ {
		if a[i] != r[len(data)-1-i] {
			t.Errorf("reverse sort didn't sort")
		}
		if i > len(data)/2 {
			break
		}
	}
}

func TestBreakPatterns(t *testing.T) {
	// Special slice used to trigger breakPatterns.
	data := make([]int, 30)
	for i := range data {
		data[i] = 10
	}
	data[(len(data)/4)*1] = 0
	data[(len(data)/4)*2] = 1
	data[(len(data)/4)*3] = 2
	PdqSort(IntSlice(data))
}

type nonDeterministicTestingData struct {
	r *rand.Rand
}

func (t *nonDeterministicTestingData) Len() int {
	return 500
}
func (t *nonDeterministicTestingData) Less(i, j int) bool {
	if i < 0 || j < 0 || i >= t.Len() || j >= t.Len() {
		panic("nondeterministic comparison out of bounds")
	}
	return t.r.Float32() < 0.5
}
func (t *nonDeterministicTestingData) Swap(i, j int) {
	if i < 0 || j < 0 || i >= t.Len() || j >= t.Len() {
		panic("nondeterministic comparison out of bounds")
	}
}

func TestNonDeterministicComparison(t *testing.T) {
	// Ensure that sort.Sort does not panic when Less returns inconsistent results.
	// See https://golang.org/issue/14377.
	defer func() {
		if r := recover(); r != nil {
			t.Error(r)
		}
	}()

	td := &nonDeterministicTestingData{
		r: rand.New(rand.NewSource(0)),
	}

	for i := 0; i < 10; i++ {
		PdqSort(td)
	}
}
