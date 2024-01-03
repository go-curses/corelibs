// Copyright (c) 2024  The Go-Curses Authors
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

package run

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRunWith(t *testing.T) {
	Convey("Run", t, func() {
		stdout, stderr, status, err := Run(
			".", "perl", "-e", ""+
				"print STDOUT \"this is stdout\n\";"+
				"print STDERR \"this is stderr\n\";",
		)
		So(err, ShouldEqual, nil)
		So(status, ShouldEqual, 0)
		So(stdout, ShouldEqual, "this is stdout\n")
		So(stderr, ShouldEqual, "this is stderr\n")
	})

	Convey("With", t, func() {
		Convey("Stderr last line", func() {
			stdout, stderr, status, err := Run(
				".", "perl", "-e", ""+
					"print STDERR \"this is an error\n\";"+
					"exit(255);",
			)
			errMsg := "this is an error"
			So(err, ShouldEqual, errors.New(errMsg))
			So(status, ShouldEqual, 255)
			So(stdout, ShouldEqual, "")
			So(stderr, ShouldEqual, errMsg+"\n")
		})

		Convey("Stdout last line", func() {
			stdout, stderr, status, err := Run(
				".", "perl", "-e", ""+
					"print STDOUT \"this is an error\n\";"+
					"exit(1);",
			)
			errMsg := "this is an error"
			So(err, ShouldEqual, errors.New(errMsg))
			So(status, ShouldEqual, 1)
			So(stdout, ShouldEqual, errMsg+"\n")
			So(stderr, ShouldEqual, "")
		})

		Convey("No last line", func() {
			stdout, stderr, status, err := Run(
				".", "perl", "-e",
				"exit(133);",
			)
			So(err, ShouldEqual, errors.New("exit status 133"))
			So(status, ShouldEqual, 133)
			So(stdout, ShouldEqual, "")
			So(stderr, ShouldEqual, "")
		})
	})
}