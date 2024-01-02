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

package maths

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAs(t *testing.T) {

	Convey("As Int, Uint, Int64, Uint64", t, func() {
		So(AsInt(1.0), ShouldEqual, int(1))
		So(AsUint(1.0), ShouldEqual, uint(1))
		So(AsInt64(1.0), ShouldEqual, int64(1))
		So(AsUint64(1.0), ShouldEqual, uint64(1))
	})

	Convey("Atoi", t, func() {
		So(Atoi("10"), ShouldEqual, 10)
	})

	Convey("ToInt", t, func() {
		So(ToInt(int(1), 0), ShouldEqual, 1)
		So(ToInt(int8(1), 0), ShouldEqual, 1)
		So(ToInt(int16(1), 0), ShouldEqual, 1)
		So(ToInt(int32(1), 0), ShouldEqual, 1)
		So(ToInt(int64(1), 0), ShouldEqual, 1)
		So(ToInt(uint(1), 0), ShouldEqual, 1)
		So(ToInt(uint8(1), 0), ShouldEqual, 1)
		So(ToInt(uint16(1), 0), ShouldEqual, 1)
		So(ToInt(uint32(1), 0), ShouldEqual, 1)
		So(ToInt(uint64(1), 0), ShouldEqual, 1)
		So(ToInt(float32(1), 0), ShouldEqual, 1)
		So(ToInt(float64(1), 0), ShouldEqual, 1)
		So(ToInt("1", 0), ShouldEqual, 0)
	})

	Convey("ToInt64", t, func() {
		So(ToInt64(int(1), 0), ShouldEqual, int64(1))
		So(ToInt64(int8(1), 0), ShouldEqual, int64(1))
		So(ToInt64(int16(1), 0), ShouldEqual, int64(1))
		So(ToInt64(int32(1), 0), ShouldEqual, int64(1))
		So(ToInt64(int64(1), 0), ShouldEqual, int64(1))
		So(ToInt64(uint(1), 0), ShouldEqual, int64(1))
		So(ToInt64(uint8(1), 0), ShouldEqual, int64(1))
		So(ToInt64(uint16(1), 0), ShouldEqual, int64(1))
		So(ToInt64(uint32(1), 0), ShouldEqual, int64(1))
		So(ToInt64(uint64(1), 0), ShouldEqual, int64(1))
		So(ToInt64(float32(1), 0), ShouldEqual, int64(1))
		So(ToInt64(float64(1), 0), ShouldEqual, int64(1))
		So(ToInt64("1", 0), ShouldEqual, int64(0))
	})

	Convey("ToUint", t, func() {
		So(ToUint(int(1), 0), ShouldEqual, uint(1))
		So(ToUint(int8(1), 0), ShouldEqual, uint(1))
		So(ToUint(int16(1), 0), ShouldEqual, uint(1))
		So(ToUint(int32(1), 0), ShouldEqual, uint(1))
		So(ToUint(int64(1), 0), ShouldEqual, uint(1))
		So(ToUint(uint(1), 0), ShouldEqual, uint(1))
		So(ToUint(uint8(1), 0), ShouldEqual, uint(1))
		So(ToUint(uint16(1), 0), ShouldEqual, uint(1))
		So(ToUint(uint32(1), 0), ShouldEqual, uint(1))
		So(ToUint(uint64(1), 0), ShouldEqual, uint(1))
		So(ToUint(float32(1), 0), ShouldEqual, uint(1))
		So(ToUint(float64(1), 0), ShouldEqual, uint(1))
		So(ToUint("1", 0), ShouldEqual, uint(0))
	})

	Convey("ToUint64", t, func() {
		So(ToUint64(int(1), 0), ShouldEqual, uint64(1))
		So(ToUint64(int8(1), 0), ShouldEqual, uint64(1))
		So(ToUint64(int16(1), 0), ShouldEqual, uint64(1))
		So(ToUint64(int32(1), 0), ShouldEqual, uint64(1))
		So(ToUint64(int64(1), 0), ShouldEqual, uint64(1))
		So(ToUint64(uint(1), 0), ShouldEqual, uint64(1))
		So(ToUint64(uint8(1), 0), ShouldEqual, uint64(1))
		So(ToUint64(uint16(1), 0), ShouldEqual, uint64(1))
		So(ToUint64(uint32(1), 0), ShouldEqual, uint64(1))
		So(ToUint64(uint64(1), 0), ShouldEqual, uint64(1))
		So(ToUint64(float32(1), 0), ShouldEqual, uint64(1))
		So(ToUint64(float64(1), 0), ShouldEqual, uint64(1))
		So(ToUint64("1", 0), ShouldEqual, uint64(0))
	})

	Convey("ToFloat64", t, func() {
		So(ToFloat64(int(1), 0), ShouldEqual, float64(1))
		So(ToFloat64(int8(1), 0), ShouldEqual, float64(1))
		So(ToFloat64(int16(1), 0), ShouldEqual, float64(1))
		So(ToFloat64(int32(1), 0), ShouldEqual, float64(1))
		So(ToFloat64(int64(1), 0), ShouldEqual, float64(1))
		So(ToFloat64(uint(1), 0), ShouldEqual, float64(1))
		So(ToFloat64(uint8(1), 0), ShouldEqual, float64(1))
		So(ToFloat64(uint16(1), 0), ShouldEqual, float64(1))
		So(ToFloat64(uint32(1), 0), ShouldEqual, float64(1))
		So(ToFloat64(uint64(1), 0), ShouldEqual, float64(1))
		So(ToFloat64(float32(1), 0), ShouldEqual, float64(1))
		So(ToFloat64(float64(1), 0), ShouldEqual, float64(1))
		So(ToFloat64("1", 0), ShouldEqual, float64(0))
	})

}
