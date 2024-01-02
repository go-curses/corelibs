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

package notify

import (
	"bytes"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNotify(t *testing.T) {

	Convey("Builder", t, func() {
		b := New(Debug)
		So(b, ShouldNotEqual, nil)
		n := b.Make()
		So(n, ShouldNotEqual, nil)
		So(n.Level(), ShouldEqual, Debug)
		n = New(Quiet).
			SetOut(nil).
			SetErr(nil).
			Make()
		So(n, ShouldNotEqual, nil)
		So(n.Level(), ShouldEqual, Quiet)
		So(n.Stdout(), ShouldEqual, nil)
		So(n.Stderr(), ShouldEqual, nil)
		n = New(Quiet).
			SetLevel(Info).
			Make()
		So(n, ShouldNotEqual, nil)
		So(n.Level(), ShouldEqual, Info)
	})

	Convey("Notifier", t, func() {
		Convey("Defaults", func() {
			n := New(Debug).Make()
			So(n, ShouldNotEqual, nil)
			So(n.Stdout(), ShouldEqual, os.Stdout)
			So(n.Stderr(), ShouldEqual, os.Stderr)
		})

		Convey("Modified", func() {
			n := New(Info).Make()
			So(n.Level(), ShouldEqual, Info)
			n.ModifyLevel(Debug)
			So(n.Level(), ShouldEqual, Debug)
			o, e := bytes.Buffer{}, bytes.Buffer{}
			n.ModifyOut(&o).ModifyErr(&e)
			So(n.Stdout(), ShouldNotEqual, os.Stdout)
			So(n.Stderr(), ShouldNotEqual, os.Stderr)
		})

		Convey("Info Level", func() {
			o, e := bytes.Buffer{}, bytes.Buffer{}
			n := New(Info).SetOut(&o).SetErr(&e).Make()
			n.Info("info: %v\n", "test")
			So(o.String(), ShouldEqual, "info: test\n")
			n.Debug("debug: test\n")
			So(o.String(), ShouldEqual, "info: test\n")
			So(e.String(), ShouldEqual, "")
			n.Error("error: test\n")
			So(e.String(), ShouldEqual, "error: test\n")
		})

		Convey("Debug Level", func() {
			o, e := bytes.Buffer{}, bytes.Buffer{}
			n := New(Debug).SetOut(&o).SetErr(&e).Make()
			n.Debug("debug: %v\n")
			So(o.String(), ShouldEqual, "debug: %v\n")
			So(e.String(), ShouldEqual, "")
			n.Info("info: test\n")
			So(o.String(), ShouldEqual, "debug: %v\ninfo: test\n")
			n.Error("error: test\n")
			So(e.String(), ShouldEqual, "error: test\n")
		})

		Convey("Error Level", func() {
			o, e := bytes.Buffer{}, bytes.Buffer{}
			n := New(Error).SetOut(&o).SetErr(&e).Make()
			n.Error("error: %v\n", "test")
			So(o.String(), ShouldEqual, "")
			So(e.String(), ShouldEqual, "error: test\n")
		})

	})

}
