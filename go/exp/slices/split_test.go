// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices_test

import (
	"slices"
	"testing"

	slices_ "github.com/searKing/golang/go/exp/slices"
)

var splitTests = []struct {
	s    []int
	sep  int
	want [][]int
}{
	{nil, -1, nil},
	{nil, 0, nil},
	{nil, 1, nil},
	{[]int{}, -1, [][]int{}},
	{[]int{}, 0, [][]int{}},
	{[]int{}, 1, [][]int{}},
	{[]int{1}, -1, [][]int{{1}}},
	{[]int{1}, 0, [][]int{{1}}},
	{[]int{1}, 1, [][]int{{1}}},
	{[]int{1}, 2, [][]int{{1}}},
	{[]int{1, 2, 3}, -1, [][]int{{1}, {2}, {3}}},
	{[]int{1, 2, 3}, 0, [][]int{{1}, {2}, {3}}},
	{[]int{1, 2, 3}, 1, [][]int{{1}, {2}, {3}}},
	{[]int{1, 2, 3}, 2, [][]int{{1, 2}, {3}}},
	{[]int{1, 2, 3}, 3, [][]int{{1, 2, 3}}},
	{[]int{1, 2, 3}, 4, [][]int{{1, 2, 3}}},
	{[]int{1, 2, 3, 4}, -1, [][]int{{1}, {2}, {3}, {4}}},
	{[]int{1, 2, 3, 4}, 0, [][]int{{1}, {2}, {3}, {4}}},
	{[]int{1, 2, 3, 4}, 1, [][]int{{1}, {2}, {3}, {4}}},
	{[]int{1, 2, 3, 4}, 2, [][]int{{1, 2}, {3, 4}}},
	{[]int{1, 2, 3, 4}, 3, [][]int{{1, 2, 3}, {4}}},
	{[]int{1, 2, 3, 4}, 4, [][]int{{1, 2, 3, 4}}},
	{[]int{1, 2, 3, 4}, 5, [][]int{{1, 2, 3, 4}}},
}

func TestSplit(t *testing.T) {
	for i, test := range splitTests {
		got := slices_.Split(test.s, test.sep)
		if len(test.want) != len(got) {
			t.Errorf("#%d: Split(%v, %v) = %v, want %v", i, test.s, test.sep, got, test.want)
			continue
		}
		for i := range got {
			if !slices.Equal(got[i], test.want[i]) {
				t.Errorf("#%d: Split(%v, %v) = %v, want %v", i, test.s, test.sep, got, test.want)
				break
			}
		}
	}
}

var splitNTests = []struct {
	s    []int
	sep  int
	n    int
	want [][]int
}{
	{nil, -1, -1, nil},
	{nil, -1, 0, nil},
	{nil, -1, 1, nil},
	{nil, 0, -1, nil},
	{nil, 0, 0, nil},
	{nil, 0, 1, nil},
	{[]int{}, -1, -1, nil},
	{[]int{}, -1, 0, nil},
	{[]int{}, -1, 1, [][]int{{}}},
	{[]int{}, 0, -1, nil},
	{[]int{}, 0, 0, nil},
	{[]int{}, 0, 1, [][]int{}},
	{[]int{}, 1, -1, nil},
	{[]int{}, 1, 0, nil},
	{[]int{}, 1, 1, [][]int{{}}},
	{[]int{1}, -1, -1, [][]int{{1}}},
	{[]int{1}, -1, 0, nil},
	{[]int{1}, -1, 1, [][]int{{1}}},
	{[]int{1}, -1, 2, [][]int{{1}}},
	{[]int{1}, 0, -1, [][]int{{1}}},
	{[]int{1}, 0, 0, nil},
	{[]int{1}, 0, 1, [][]int{{1}}},
	{[]int{1}, 0, 2, [][]int{{1}}},
	{[]int{1}, 1, -1, [][]int{{1}}},
	{[]int{1}, 1, 0, nil},
	{[]int{1}, 1, 1, [][]int{{1}}},
	{[]int{1}, 1, 2, [][]int{{1}}},
	{[]int{1}, 2, -1, [][]int{{1}}},
	{[]int{1}, 2, 0, nil},
	{[]int{1}, 2, 1, [][]int{{1}}},
	{[]int{1}, 2, 2, [][]int{{1}}},
	{[]int{1, 2, 3}, -1, -1, [][]int{{1}, {2}, {3}}},
	{[]int{1, 2, 3}, -1, 0, nil},
	{[]int{1, 2, 3}, -1, 1, [][]int{{1, 2, 3}}},
	{[]int{1, 2, 3}, -1, 2, [][]int{{1}, {2, 3}}},
	{[]int{1, 2, 3}, -1, 3, [][]int{{1}, {2}, {3}}},
	{[]int{1, 2, 3}, -1, 4, [][]int{{1}, {2}, {3}}},
	{[]int{1, 2, 3}, 0, -1, [][]int{{1}, {2}, {3}}},
	{[]int{1, 2, 3}, 0, 0, nil},
	{[]int{1, 2, 3}, 0, 1, [][]int{{1, 2, 3}}},
	{[]int{1, 2, 3}, 0, 2, [][]int{{1}, {2, 3}}},
	{[]int{1, 2, 3}, 0, 3, [][]int{{1}, {2}, {3}}},
	{[]int{1, 2, 3}, 0, 4, [][]int{{1}, {2}, {3}}},
	{[]int{1, 2, 3}, 1, -1, [][]int{{1}, {2}, {3}}},
	{[]int{1, 2, 3}, 1, 0, nil},
	{[]int{1, 2, 3}, 1, 1, [][]int{{1, 2, 3}}},
	{[]int{1, 2, 3}, 1, 2, [][]int{{1}, {2, 3}}},
	{[]int{1, 2, 3}, 1, 3, [][]int{{1}, {2}, {3}}},
	{[]int{1, 2, 3}, 1, 4, [][]int{{1}, {2}, {3}}},
	{[]int{1, 2, 3}, 2, -1, [][]int{{1, 2}, {3}}},
	{[]int{1, 2, 3}, 2, 0, nil},
	{[]int{1, 2, 3}, 2, 1, [][]int{{1, 2, 3}}},
	{[]int{1, 2, 3}, 2, 2, [][]int{{1, 2}, {3}}},
	{[]int{1, 2, 3}, 2, 3, [][]int{{1, 2}, {3}}},
	{[]int{1, 2, 3}, 2, 4, [][]int{{1, 2}, {3}}},
	{[]int{1, 2, 3}, 3, -1, [][]int{{1, 2, 3}}},
	{[]int{1, 2, 3}, 3, 0, nil},
	{[]int{1, 2, 3}, 3, 1, [][]int{{1, 2, 3}}},
	{[]int{1, 2, 3}, 3, 2, [][]int{{1, 2, 3}}},
	{[]int{1, 2, 3}, 3, 3, [][]int{{1, 2, 3}}},
	{[]int{1, 2, 3}, 3, 4, [][]int{{1, 2, 3}}},
	{[]int{1, 2, 3}, 4, -1, [][]int{{1, 2, 3}}},
	{[]int{1, 2, 3}, 4, 0, nil},
	{[]int{1, 2, 3}, 4, 1, [][]int{{1, 2, 3}}},
	{[]int{1, 2, 3}, 4, 2, [][]int{{1, 2, 3}}},
	{[]int{1, 2, 3}, 4, 3, [][]int{{1, 2, 3}}},
	{[]int{1, 2, 3}, 4, 4, [][]int{{1, 2, 3}}},
	{[]int{1, 2, 3, 4}, -1, -1, [][]int{{1}, {2}, {3}, {4}}},
	{[]int{1, 2, 3, 4}, -1, 0, nil},
	{[]int{1, 2, 3, 4}, -1, 1, [][]int{{1, 2, 3, 4}}},
	{[]int{1, 2, 3, 4}, -1, 2, [][]int{{1, 2}, {3, 4}}},
	{[]int{1, 2, 3, 4}, -1, 3, [][]int{{1}, {2}, {3, 4}}},
	{[]int{1, 2, 3, 4}, -1, 4, [][]int{{1}, {2}, {3}, {4}}},
	{[]int{1, 2, 3, 4}, 0, -1, [][]int{{1}, {2}, {3}, {4}}},
	{[]int{1, 2, 3, 4}, 0, 0, nil},
	{[]int{1, 2, 3, 4}, 0, 1, [][]int{{1, 2, 3, 4}}},
	{[]int{1, 2, 3, 4}, 0, 2, [][]int{{1, 2}, {3, 4}}},
	{[]int{1, 2, 3, 4}, 0, 3, [][]int{{1}, {2}, {3, 4}}},
	{[]int{1, 2, 3, 4}, 0, 4, [][]int{{1}, {2}, {3}, {4}}},
	{[]int{1, 2, 3, 4}, 1, -1, [][]int{{1}, {2}, {3}, {4}}},
	{[]int{1, 2, 3, 4}, 1, 0, nil},
	{[]int{1, 2, 3, 4}, 1, 1, [][]int{{1, 2, 3, 4}}},
	{[]int{1, 2, 3, 4}, 1, 2, [][]int{{1}, {2, 3, 4}}},
	{[]int{1, 2, 3, 4}, 1, 3, [][]int{{1}, {2}, {3, 4}}},
	{[]int{1, 2, 3, 4}, 1, 4, [][]int{{1}, {2}, {3}, {4}}},
	{[]int{1, 2, 3, 4}, 2, -1, [][]int{{1, 2}, {3, 4}}},
	{[]int{1, 2, 3, 4}, 2, 0, nil},
	{[]int{1, 2, 3, 4}, 2, 1, [][]int{{1, 2, 3, 4}}},
	{[]int{1, 2, 3, 4}, 2, 2, [][]int{{1, 2}, {3, 4}}},
	{[]int{1, 2, 3, 4}, 2, 3, [][]int{{1, 2}, {3, 4}}},
	{[]int{1, 2, 3, 4}, 2, 4, [][]int{{1, 2}, {3, 4}}},
	{[]int{1, 2, 3, 4}, 3, -1, [][]int{{1, 2, 3}, {4}}},
	{[]int{1, 2, 3, 4}, 3, 0, nil},
	{[]int{1, 2, 3, 4}, 3, 1, [][]int{{1, 2, 3, 4}}},
	{[]int{1, 2, 3, 4}, 3, 2, [][]int{{1, 2, 3}, {4}}},
	{[]int{1, 2, 3, 4}, 3, 3, [][]int{{1, 2, 3}, {4}}},
	{[]int{1, 2, 3, 4}, 3, 4, [][]int{{1, 2, 3}, {4}}},
	{[]int{1, 2, 3, 4}, 3, 5, [][]int{{1, 2, 3}, {4}}},
	{[]int{1, 2, 3, 4}, 4, -1, [][]int{{1, 2, 3, 4}}},
	{[]int{1, 2, 3, 4}, 4, 0, nil},
	{[]int{1, 2, 3, 4}, 4, 1, [][]int{{1, 2, 3, 4}}},
	{[]int{1, 2, 3, 4}, 4, 2, [][]int{{1, 2, 3, 4}}},
	{[]int{1, 2, 3, 4}, 4, 3, [][]int{{1, 2, 3, 4}}},
	{[]int{1, 2, 3, 4}, 4, 4, [][]int{{1, 2, 3, 4}}},
	{[]int{1, 2, 3, 4}, 4, 5, [][]int{{1, 2, 3, 4}}},
	{[]int{1, 2, 3, 4}, 5, -1, [][]int{{1, 2, 3, 4}}},
	{[]int{1, 2, 3, 4}, 5, 0, nil},
	{[]int{1, 2, 3, 4}, 5, 1, [][]int{{1, 2, 3, 4}}},
	{[]int{1, 2, 3, 4}, 5, 2, [][]int{{1, 2, 3, 4}}},
	{[]int{1, 2, 3, 4}, 5, 3, [][]int{{1, 2, 3, 4}}},
	{[]int{1, 2, 3, 4}, 5, 4, [][]int{{1, 2, 3, 4}}},
	{[]int{1, 2, 3, 4}, 5, 5, [][]int{{1, 2, 3, 4}}},
}

func TestSplitN(t *testing.T) {
	for i, test := range splitNTests {
		got := slices_.SplitN(test.s, test.sep, test.n)
		if len(test.want) != len(got) {
			t.Errorf("#%d: SplitN(%v, %v, %v) = %v, want %v", i, test.s, test.sep, test.n, got, test.want)
			continue
		}
		for i := range got {
			if !slices.Equal(got[i], test.want[i]) {
				t.Errorf("#%d: SplitN(%v, %v, %v) = %v, want %v", i, test.s, test.sep, test.n, got, test.want)
				break
			}
		}
	}
}
