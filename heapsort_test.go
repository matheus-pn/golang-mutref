// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math/rand"
	"testing"
)

const (
	_Sawtooth = iota
	_Rand
	_Stagger
	_Plateau
	_Shuffle
	_NDist
)

const (
	_Copy = iota
	_Reverse
	_ReverseFirstHalf
	_ReverseSecondHalf
	_Sorted
	_Dither
	_NMode
)

// IsSorted reports whether data is sorted.
func IsSorted(data Interface) bool {
	n := data.Len()
	for i := n - 1; i > 0; i-- {
		if data.Less(i, i-1) {
			return false
		}
	}
	return true
}

// IntsAreSorted reports whether the slice x is sorted in increasing order.
func IntsAreSorted(x []int) bool { return IsSorted(IntSlice(x)) }

type testingData struct {
	desc        string
	t           *testing.T
	data        []int
	maxswap     int // number of swaps allowed
	ncmp, nswap int
}

func (d *testingData) Len() int { return len(d.data) }
func (d *testingData) Less(i, j int) bool {
	d.ncmp++
	return d.data[i] < d.data[j]
}
func (d *testingData) Swap(i, j int) {
	if d.nswap >= d.maxswap {
		d.t.Fatalf("%s: used %d swaps sorting slice of %d", d.desc, d.nswap, len(d.data))
	}
	d.nswap++
	d.data[i], d.data[j] = d.data[j], d.data[i]
}

func lg(n int) int {
	i := 0
	for 1<<uint(i) < n {
		i++
	}
	return i
}

type IntSlice []int

func (x IntSlice) Len() int           { return len(x) }
func (x IntSlice) Less(i, j int) bool { return x[i] < x[j] }
func (x IntSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// Sort is a convenience method: x.Sort() calls Sort(x).
func (x IntSlice) Sort() { Sort(x) }

func Ints(x []int) { Sort(IntSlice(x)) }

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestHeapsortBM(t *testing.T) {
	testBentleyMcIlroy(t, Heap, func(n int) int { return n * lg(n) * 12 / 10 })
}

func testBentleyMcIlroy(t *testing.T, sort func(Interface), maxswap func(int) int) {
	sizes := []int{100, 1023, 1024, 1025}
	if testing.Short() {
		sizes = []int{100, 127, 128, 129}
	}
	dists := []string{"sawtooth", "rand", "stagger", "plateau", "shuffle"}
	modes := []string{"copy", "reverse", "reverse1", "reverse2", "sort", "dither"}
	var tmp1, tmp2 [1025]int
	for _, n := range sizes {
		for m := 1; m < 2*n; m *= 2 {
			for dist := 0; dist < _NDist; dist++ {
				j := 0
				k := 1
				data := tmp1[0:n]
				for i := 0; i < n; i++ {
					switch dist {
					case _Sawtooth:
						data[i] = i % m
					case _Rand:
						data[i] = rand.Intn(m)
					case _Stagger:
						data[i] = (i*m + i) % n
					case _Plateau:
						data[i] = min(i, m)
					case _Shuffle:
						if rand.Intn(m) != 0 {
							j += 2
							data[i] = j
						} else {
							k += 2
							data[i] = k
						}
					}
				}

				mdata := tmp2[0:n]
				for mode := 0; mode < _NMode; mode++ {
					switch mode {
					case _Copy:
						for i := 0; i < n; i++ {
							mdata[i] = data[i]
						}
					case _Reverse:
						for i := 0; i < n; i++ {
							mdata[i] = data[n-i-1]
						}
					case _ReverseFirstHalf:
						for i := 0; i < n/2; i++ {
							mdata[i] = data[n/2-i-1]
						}
						for i := n / 2; i < n; i++ {
							mdata[i] = data[i]
						}
					case _ReverseSecondHalf:
						for i := 0; i < n/2; i++ {
							mdata[i] = data[i]
						}
						for i := n / 2; i < n; i++ {
							mdata[i] = data[n-(i-n/2)-1]
						}
					case _Sorted:
						for i := 0; i < n; i++ {
							mdata[i] = data[i]
						}
						// Ints is known to be correct
						// because mode Sort runs after mode _Copy.
						Ints(mdata)
					case _Dither:
						for i := 0; i < n; i++ {
							mdata[i] = data[i] + i%5
						}
					}

					desc := fmt.Sprintf("n=%d m=%d dist=%s mode=%s", n, m, dists[dist], modes[mode])
					d := &testingData{desc: desc, t: t, data: mdata[0:n], maxswap: maxswap(n)}
					sort(d)
					// Uncomment if you are trying to improve the number of compares/swaps.
					//t.Logf("%s: ncmp=%d, nswp=%d", desc, d.ncmp, d.nswap)

					// If we were testing C qsort, we'd have to make a copy
					// of the slice and sort it ourselves and then compare
					// x against it, to ensure that qsort was only permuting
					// the data, not (for example) overwriting it with zeros.
					//
					// In go, we don't have to be so paranoid: since the only
					// mutating method Sort can call is TestingData.swap,
					// it suffices here just to check that the final slice is sorted.
					if !IntsAreSorted(mdata) {
						t.Fatalf("%s: ints not sorted\n\t%v", desc, mdata)
					}
				}
			}
		}
	}
}
