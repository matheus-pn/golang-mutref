// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"sort"
	"testing"
)

func TestReverseSortIntSliceInsert(t *testing.T) {
	data := ints
	data1 := ints
	a := IntSlice(data[0:])
	InsertionSort(a)
	r := IntSlice(data1[0:])
	InsertionSort(sort.Reverse(r))
	for i := 0; i < len(data); i++ {
		if a[i] != r[len(data)-1-i] {
			t.Errorf("reverse sort didn't sort")
		}
		if i > len(data)/2 {
			break
		}
	}
}

func TestSortFloat64SliceInsert(t *testing.T) {
	data := float64s
	a := Float64Slice(data[0:])
	InsertionSort(a)
	if !IsSorted(a) {
		t.Errorf("sorted %v", float64s)
		t.Errorf("   got %v", data)
	}
}
