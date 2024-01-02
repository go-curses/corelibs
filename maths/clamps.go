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
	"math"
)

// Clamp returns the value if it is greater-than-or-equal-to the minimum and
// less-than-or-equal-to the maximum arguments; if the value is less than the
// minimum, the minimum is returned; if the value is greater than the maximum,
// the maximum is returned
func Clamp[V Number](value, min, max V) V {
	if value >= min && value <= max {
		return value
	}
	if value > max {
		return max
	}
	return min
}

// Floor returns the value if it is greater-than-or-equal-to the minimum and
// returns the minimum otherwise
func Floor[V Number](value, min V) V {
	if value < min {
		return min
	}
	return value
}

// Ceil returns the value if it is less-than-or-equal-to the maximum and
// returns the maximum otherwise
func Ceil[V Number](value, max V) V {
	if value > max {
		return max
	}
	return value
}

// Round returns the value rounded to the nearest whole value
func Round[V Decimal](x V) (rounded int) {
	return int(math.Floor(float64(x) + 0.5))
}

// RoundUp returns the value rounded up
func RoundUp[V Decimal](value V) (rounded int) {
	return int(math.Ceil(float64(value)))
}

// RoundDown returns the value rounded down
func RoundDown[V Decimal](value V) (rounded int) {
	return int(math.Floor(float64(value)))
}
