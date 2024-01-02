// Copyright (c) 2023  The Go-Curses Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package maps

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSortedKeys(t *testing.T) {

	Convey("Standard", t, func() {
		m := map[string]struct{}{
			"one":  {},
			"two":  {},
			"many": {},
		}
		So(SortedKeys(m), ShouldEqual, []string{
			"many", "one", "two",
		})
	})

}

func TestSortedNumbers(t *testing.T) {

	Convey("SortedNumbers", t, func() {
		mInt := map[int]struct{}{
			2: {},
			1: {},
			0: {},
		}
		So(SortedNumbers(mInt), ShouldEqual, []int{
			0, 1, 2,
		})
		mFloat := map[float32]struct{}{
			2.5: {},
			1.0: {},
			3.4: {},
		}
		So(SortedNumbers(mFloat), ShouldEqual, []float32{
			1.0, 2.5, 3.4,
		})
	})

	Convey("ReverseSortedNumbers", t, func() {
		mInt := map[int]struct{}{
			2: {},
			1: {},
			0: {},
		}
		So(ReverseSortedNumbers(mInt), ShouldEqual, []int{
			2, 1, 0,
		})
		mFloat := map[float32]struct{}{
			2.5: {},
			1.0: {},
			3.4: {},
		}
		So(ReverseSortedNumbers(mFloat), ShouldEqual, []float32{
			3.4, 2.5, 1.0,
		})
	})

}
