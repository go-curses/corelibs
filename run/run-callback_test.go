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
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCallback(t *testing.T) {
	Convey("Success", t, func() {
		Convey("Nil", func() {
			pid, _, err := Callback(
				Options{
					Path: ".",
					Name: "perl",
					Argv: []string{
						"-e",
						"exit(0);",
					},
				},
				nil, nil,
			)
			So(err, ShouldEqual, nil)
			So(pid, ShouldNotEqual, 0)
		})

		Convey("Outputs", func() {
			var stdout, stderr []string
			pid, done, err := Callback(
				Options{
					Path: ".",
					Name: "perl",
					Argv: []string{
						"-e", "" +
							"print STDOUT \"this is stdout\n\";" +
							"print STDERR \"this is stderr\n\";" +
							"exit(0);",
					},
				},
				func(line string) {
					stdout = append(stdout, line)
				},
				func(line string) {
					stderr = append(stderr, line)
				},
			)
			So(err, ShouldEqual, nil)
			So(pid, ShouldNotEqual, 0)
			<-done // wait for exit
			So(stdout, ShouldEqual, []string{"this is stdout"})
			So(stderr, ShouldEqual, []string{"this is stderr"})
		})

		Convey("Exec Errors", func() {
			var stdout, stderr []string
			pid, _, err := Callback(
				Options{
					Path: ".",
					Name: "",
				},
				func(line string) {
					stdout = append(stdout, line)
				},
				func(line string) {
					stderr = append(stderr, line)
				},
			)
			So(err.Error(), ShouldEqual, "exec: no command")
			So(pid, ShouldEqual, 0)
			So(stdout, ShouldEqual, []string(nil))
			So(stderr, ShouldEqual, []string(nil))
		})

		Convey("Command Errors", func() {
			var stdout, stderr []string
			pid, done, err := Callback(
				Options{
					Path: ".",
					Name: "perl",
					Argv: []string{
						"-e",
						"exit(255);",
					},
				},
				func(line string) {
					stdout = append(stdout, line)
				},
				func(line string) {
					stderr = append(stderr, line)
				},
			)
			So(err, ShouldEqual, nil)
			So(pid, ShouldNotEqual, 0)
			<-done
			So(stdout, ShouldEqual, []string(nil))
			So(stderr, ShouldEqual, []string{"exit status 255"})
		})
	})
}