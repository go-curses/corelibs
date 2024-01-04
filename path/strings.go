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

package path

import (
	"path/filepath"
	"strings"
)

// CleanWithSlash trims all leading and trailing space from the `path` given,
// uses filepath.Clean on the trimmed results. The cleaned results are
// prefixed with a leading slash producing the `clean` value of `path`.
//
// CleanWithSlash supports the Go-Enjin convention of a special type of path
// that is prefixed with a `!` instead of a slash and will correctly return
// the `clean` value with that prefix if it was present and if not uses that
// normal path separator slash.
//
// CleanWithSlash produces absolute paths but do not mistake these for
// actual filesystem paths. CleanWithSlash is intended for use in working with
// URL paths.
func CleanWithSlash(path string) (clean string) {
	var lead, trimmed string
	trimmed = strings.TrimSpace(path)
	if lead = "/"; strings.HasPrefix(path, "!") {
		lead = "!"
	}
	if len(trimmed) > 2 {
		if trimmed[0] == lead[0] {
			trimmed = trimmed[1:]
		}
		if last := len(trimmed) - 1; last >= 0 && trimmed[last] == '/' {
			trimmed = trimmed[:last]
		}
	}
	cleaned := filepath.Clean(trimmed)
	if cleaned == "" || cleaned == "." || cleaned == lead {
		cleaned = ""
	}
	clean = lead + cleaned
	return
}

// CleanWithSlashes is a wrapper around CleanWithSlash which ensures the
// `clean` path begins and ends with a slash
func CleanWithSlashes(path string) (clean string) {
	if clean = CleanWithSlash(path); clean != "/" {
		clean += "/"
	}
	return
}

// Join joins the given URL path parts and uses filepath.Clean on the results
func Join(parts ...string) (joined string) {
	joined = strings.Join(parts, "/")
	joined = filepath.Clean(joined)
	return
}

// JoinWithSlash joins the given URL path parts and uses CleanWithSlash on the
// results
func JoinWithSlash(paths ...string) (joined string) {
	joined = strings.Join(paths, "/")
	joined = CleanWithSlash(joined)
	return
}

// JoinWithSlashes joins the given URL path parts and uses CleanWithSlashes on
// the results
func JoinWithSlashes(paths ...string) (joined string) {
	joined = strings.Join(paths, "/")
	joined = CleanWithSlashes(joined)
	return
}

// TrimSlash returns the filepath cleaned and without any trailing slash
func TrimSlash(path string) (clean string) {
	if path == "" {
		return
	}
	clean = strings.TrimSpace(path)
	clean = filepath.Clean(clean)
	clean = strings.TrimSuffix(clean, "/")
	return
}

// TrimSlashes returns the filepath cleaned and without any leading or
// trailing slashes
func TrimSlashes(path string) (clean string) {
	if path == "" {
		return
	}
	clean = strings.TrimSpace(path)
	clean = filepath.Clean(clean)
	clean = strings.TrimPrefix(clean, "/")
	clean = strings.TrimSuffix(clean, "/")
	return
}

// SafeConcatRelPath prunes all empty and current dir paths from the
// list given, using TrimSlashes, then joins all the paths together,
// ensuring the output is prefixed with the given root, and has no
// leading or trailing slash and has is filepath.Clean
func SafeConcatRelPath(root string, paths ...string) (out string) {
	var outs []string
	for _, path := range paths {
		if v := TrimSlashes(path); v != "" && v != "." {
			outs = append(outs, v)
		}
	}
	out = strings.Join(outs, "/")
	root = TrimSlashes(root)
	out = strings.TrimPrefix(out, root)
	out = root + "/" + out
	out = filepath.Clean(out)
	if size := len(out); size > 0 && out[0] == '/' {
		out = out[1:]
	}
	return
}

// SafeConcatUrlPath wraps SafeConcatRelPath with a prefixing slash making
// an absolute URL path
func SafeConcatUrlPath(root string, paths ...string) (out string) {
	out = "/" + SafeConcatRelPath(root, paths...)
	return
}

// TrimPrefix is like strings.TrimPrefix but for URL paths so that the prefix
// looked for is a prefixing URL path regardless of any slashes present in the
// prefix or path argument values.
//
// TrimPrefix is primarily used in Go-Enjin filesystem driver implementations
// to coalesce given path into an underlying filesystem's actual path.
//
// Example:
//
//	TrimPrefix("/one", "one") == ""
//	TrimPrefix("/one/two/many", "one") == "two/many"
//	TrimPrefix("/one/two/many", "two/many") == "one/two/many"
func TrimPrefix(path, prefix string) (modified string) {
	prefix = TrimSlashes(prefix)
	modified = TrimSlashes(path)
	if pl := len(prefix); pl > 0 {
		if ml := len(modified); ml > pl {
			if modified[0:pl] == prefix {
				modified = modified[pl+1:]
			}
		} else if modified == prefix {
			return ""
		}
	}
	modified = TrimSlashes(modified)
	return
}

// TrimDotSlash trims any leading `./` from the given path
func TrimDotSlash(path string) (out string) {
	if out = path; len(out) >= 2 && out[0:2] == "./" {
		out = out[2:]
	}
	return
}

// TopDirectory returns the top directory in the path given
func TopDirectory(path string) (name string) {
	name = TrimSlashes(path)
	name, _, _ = strings.Cut(name, "/")
	return
}

// MatchExact compares cleaned versions of path and prefix
func MatchExact(path, prefix string) (match bool) {
	match = CleanWithSlash(path) == CleanWithSlash(prefix)
	return
}

// MatchCut returns the suffix of path if path is prefixed with the given
// prefix path. If the path and prefix match exactly, suffix will be empty
// and matched will be true
//
// Examples:
//
//	MatchCut("/one/two/", "one/two") == "", true
//	MatchCut("/one/two/many", "one/two") == "many", true
//	MatchCut("/one/two/many", "one") == "two/many", true
func MatchCut(path, prefix string) (suffix string, matched bool) {
	path = CleanWithSlash(path)
	prefix = CleanWithSlash(prefix)
	if matched = path == prefix; matched {
		return
	} else if matched = strings.HasPrefix(path, prefix+"/"); matched {
		suffix = path[len(prefix)+1:]
	}
	return
}

// Base returns the base name of the path without any file extensions
func Base(path string) (name string) {
	name = filepath.Base(path)
	for extn := filepath.Ext(name); extn != ""; extn = filepath.Ext(name) {
		name = name[:len(name)-len(extn)]
	}
	return
}

// BasePath returns the path without any primary or secondary file extensions,
// tertiary and other extensions present will remain as-is
//
// Example:
//
//	BasePath("one/file.txt.tmpl.bak") == "one/file.txt"
func BasePath(path string) (basePath string) {
	basePath = path
	extn, extra := ExtExt(path)
	if extn != "" {
		basePath = strings.TrimSuffix(basePath, "."+extn)
		if extra != "" {
			basePath = strings.TrimSuffix(basePath, "."+extra)
		}
	}
	return
}

// Ext returns the extension of the file (without the dot)
func Ext(path string) (extn string) {
	if extn = filepath.Ext(path); extn != "" {
		extn = extn[1:]
	}
	return
}

// ExtExt returns the extension of the file (without the dot) and any secondary
// extension found in the path
//
// Example:
//
//	ExtExt("page.html.tmpl") => "tmpl", "html"
func ExtExt(path string) (primary, secondary string) {
	primary = Ext(path)
	trimmed := TrimExt(path)
	secondary = Ext(trimmed)
	return
}

// TrimExt returns the path without the last extension, if there are any
func TrimExt(path string) (out string) {
	if extn := filepath.Ext(path); extn != "" {
		out = path[0 : len(path)-len(extn)]
	}
	return
}

// TrimRelativeToRoot truncates the root from path and removes any leading
// or trailing slashes. If the root is not present in the path, rel is empty
func TrimRelativeToRoot(path, root string) (rel string) {
	path = TrimSlashes(path)
	root = TrimSlashes(root)
	if rootLen := len(root); len(path) >= rootLen {
		if path[:rootLen] == root {
			rel = path[rootLen+1:]
			rel = TrimSlashes(rel)
		}
	}
	return
}

// ParseParentPaths is expecting a directory path and returns a list that
// walks up the total path with each item in the list
//
// Example:
//
//	ParseParentPaths("one/two/many") == []string{
//	    "one",
//	    "one/two",
//	    "one/two/many",
//	}
func ParseParentPaths(path string) (parents []string) {
	if size := len(path); size == 0 {
		return
	} else if size > 0 && path[0] == '/' {
		path = path[1:]
	}
	parts := strings.Split(path, "/")
	if len(parts) == 1 && parts[0] == "" {
		return
	}
	for i := 0; i < len(parts); i++ {
		var parent string
		parent = strings.Join(parts[0:i+1], "/")
		parents = append(parents, parent)
	}
	return
}