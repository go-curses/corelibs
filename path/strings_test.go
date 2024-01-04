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

package path

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStrings(t *testing.T) {
	Convey("CleanWithSlash, CleanWithSlashes", t, func() {
		So(CleanWithSlash("thing"), ShouldEqual, "/thing")
		So(CleanWithSlash("./thing"), ShouldEqual, "/thing")
		So(CleanWithSlash("/thing/"), ShouldEqual, "/thing")
		So(CleanWithSlash("!thing"), ShouldEqual, "!thing")
		So(CleanWithSlash("/"), ShouldEqual, "/")
		So(CleanWithSlashes("/"), ShouldEqual, "/")
		So(CleanWithSlashes("/thing"), ShouldEqual, "/thing/")
	})

	Convey("Join, JoinWithSlash, JoinWithSlashes", t, func() {
		So(Join("one"), ShouldEqual, "one")
		So(Join("one", "two"), ShouldEqual, "one/two")
		So(Join("/one", "/two"), ShouldEqual, "/one/two")
		So(JoinWithSlash("one"), ShouldEqual, "/one")
		So(JoinWithSlash("one", "two"), ShouldEqual, "/one/two")
		So(JoinWithSlash("/one", "/two"), ShouldEqual, "/one/two")
		So(JoinWithSlashes("one"), ShouldEqual, "/one/")
		So(JoinWithSlashes("one", "two"), ShouldEqual, "/one/two/")
		So(JoinWithSlashes("/one", "/two"), ShouldEqual, "/one/two/")
	})

	Convey("TrimSlash, TrimSlashes", t, func() {
		So(TrimSlash(""), ShouldEqual, "")
		So(TrimSlash("one/"), ShouldEqual, "one")
		So(TrimSlash("one/two/"), ShouldEqual, "one/two")
		So(TrimSlash("/one/two/"), ShouldEqual, "/one/two")
		So(TrimSlashes("one/"), ShouldEqual, "one")
		So(TrimSlashes("one/two/"), ShouldEqual, "one/two")
		So(TrimSlashes("/one/two/"), ShouldEqual, "one/two")
	})

	Convey("SafeConcatRelPath, SafeConcatUrlPath", t, func() {
		So(SafeConcatRelPath(""), ShouldEqual, "")
		So(SafeConcatRelPath("", "", "."), ShouldEqual, "")
		So(SafeConcatRelPath("", "/one", ".", "two"), ShouldEqual, "one/two")
		So(SafeConcatRelPath("top/", "/one", ".", "two"), ShouldEqual, "top/one/two")
		So(SafeConcatRelPath("/top/", "/one", ".", "two"), ShouldEqual, "top/one/two")
		So(SafeConcatUrlPath(""), ShouldEqual, "/")
		So(SafeConcatUrlPath("/top/", "/one", ".", "two"), ShouldEqual, "/top/one/two")
	})

	Convey("TrimPrefix, TrimDotSlash", t, func() {
		So(TrimPrefix("/one", "one"), ShouldEqual, "")
		So(TrimPrefix("/one/two/many", "one"), ShouldEqual, "two/many")
		So(TrimPrefix("one//two/many", "one"), ShouldEqual, "two/many")
		So(TrimPrefix("/one/two/many", "two/many"), ShouldEqual, "one/two/many")
		So(TrimDotSlash("/nope"), ShouldEqual, "/nope")
		So(TrimDotSlash("./yup"), ShouldEqual, "yup")
	})

	Convey("TopDirectory", t, func() {
		So(TopDirectory("/one/two"), ShouldEqual, "one")
		So(TopDirectory("one/two"), ShouldEqual, "one")
		So(TopDirectory("/one"), ShouldEqual, "one")
	})

	Convey("MatchExact", t, func() {
		So(MatchExact("one/two", "one/two"), ShouldEqual, true)
		So(MatchExact("/one/two", "one/two/"), ShouldEqual, true)
		So(MatchExact("/one/two", "two/one"), ShouldEqual, false)
	})

	Convey("MatchCut", t, func() {
		check := func(path, prefix, expectedSuffix string, expectedMatched bool) {
			suffix, matched := MatchCut(path, prefix)
			So(matched, ShouldEqual, expectedMatched)
			So(suffix, ShouldEqual, expectedSuffix)
		}
		check("/one/two/", "one/two", "", true)
		check("/one/two/", "one/", "two", true)
		check("/one/two/", "many/", "", false)
	})

	Convey("Base, BasePath", t, func() {
		So(Base(""), ShouldEqual, "")
		So(Base("/one"), ShouldEqual, "one")
		So(Base("/one/file"), ShouldEqual, "file")
		So(Base("/one/file.txt"), ShouldEqual, "file")
		So(Base("/one/file.txt.tmpl"), ShouldEqual, "file")
		So(Base("/one/file.txt.tmpl.bak"), ShouldEqual, "file")
		So(BasePath(""), ShouldEqual, "")
		So(BasePath("/one"), ShouldEqual, "/one")
		So(BasePath("/one/file"), ShouldEqual, "/one/file")
		So(BasePath("/one/file.txt"), ShouldEqual, "/one/file")
		So(BasePath("/one/file.txt.tmpl"), ShouldEqual, "/one/file")
		So(BasePath("/one/file.txt.tmpl.bak"), ShouldEqual, "/one/file.txt")
	})

	Convey("TrimRelativeToRoot", t, func() {
		So(TrimRelativeToRoot("/one/two", "root"), ShouldEqual, "")
		So(TrimRelativeToRoot("/one/two", "one"), ShouldEqual, "two")
	})

	Convey("ParseParentPaths", t, func() {
		So(ParseParentPaths(""), ShouldEqual, []string(nil))
		So(ParseParentPaths("/"), ShouldEqual, []string(nil))
		So(ParseParentPaths("/one/two/many"), ShouldEqual, []string{
			"one",
			"one/two",
			"one/two/many",
		})
	})
}