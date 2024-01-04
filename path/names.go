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
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	RxIncrementNameSuffix = regexp.MustCompile(`\((\d+)\)\s*$`)
	RxIncrementPathSuffix = regexp.MustCompile(`.(\d+)\s*$`)
)

// IncrementLabel examines the given label for the specific sequence of leader,
// prefix, one or more digits and suffix characters. If the sequence is found,
// the digits are parsed into an integer and incremented by the given increment
// value and reassembled into the modified label output. If the sequence is not
// found, the sequence is appended to the whole label given with a number of 1.
//
// If the increment given is zero, an increment of 1 is assumed as an increment
// of zero would not actually increment anything. Negative increments are
// valid, resulting in a decrement rather than actual increment.
//
// Examples:
//
//	IncrementLabel("Title", " ", "(", ")", 1) == "Title (1)"
//	IncrementLabel("Title (10)", " ", "(", ")", 1) == "Title (11)"
//	IncrementLabel("filename", ".", "", "", 1) == "filename.1"
//	IncrementLabel("filename.10", ".", "", "", 1) == "filename.11"
//	IncrementLabel("filename.10", ".", "", "", -1) == "filename.9"
func IncrementLabel(label, leader, prefix, suffix string, increment int) (modified string, integer int) {
	if len(label) == 0 {
		return // nop, early out
	} else if increment == 0 {
		increment = 1
	}

	suffixLen, prefixLen, leaderLen := len(suffix), len(prefix), len(leader)

	var number string
	labelLen := len(label)
	cursor := labelLen - suffixLen

	if label[cursor:] == suffix {
		// found suffix, look for number
		for idx := cursor - 1; idx >= -1; idx-- {
			if idx > -1 && unicode.IsDigit(rune(label[idx])) {
				// prepend to maintain order of digits in the number
				number = string(label[idx]) + number
				continue // keep checking for more digits
			}
			break
		}
		if numberLen := len(number); numberLen > 0 {
			// found number, look for prefix
			cursor -= numberLen
			if label[cursor-prefixLen:cursor] == prefix {
				// found prefix look for leader
				cursor -= prefixLen
				if label[cursor-leaderLen:cursor] == leader {
					// found leader, modify label
					cursor -= leaderLen
					// Atoi err is nil unless it failed to convert a valid
					// case, which given that there are unit tests for
					// standard packages, is not plausible or requiring
					// formal error handling.
					integer, _ = strconv.Atoi(number)
					if integer+increment <= 0 {
						// decremented to no increment in the modified label
						modified = label[:cursor]
						return
					}
					incremented := strconv.Itoa(integer + increment)
					modified = label[:cursor] + leader + prefix + incremented + suffix
					return
				}
			}
		}
	}

	// this is the first increment
	if increment <= 0 {
		integer = 1
		modified = label + leader + prefix + "1" + suffix
		return
	}

	integer = increment
	modified = label + leader + prefix + strconv.Itoa(integer) + suffix
	return
}

// IncrementFileName is a wrapper around IncrementLabel, equivalent to:
// `IncrementLabel(name, " ", "(", ")", 1)`
func IncrementFileName(name string) (modified string) {
	modified, _ = IncrementLabel(name, " ", "(", ")", 1)
	return
}

// IncrementFilePath is a wrapper around IncrementLabel, equivalent to:
// `IncrementLabel(name, ".", "", "", 1)`
func IncrementFilePath(name string) (modified string) {
	modified, _ = IncrementLabel(name, ".", "", "", 1)
	return
}

// IncrementFileBackup removes the given extension if present, uses
// IncrementFilePath to modify the name with the extension appended.
// If the extension is empty, a "~" is used.
//
// Examples:
//
//	IncrementFileBackup("test.txt", ".bak") == "test.txt.bak"
//	IncrementFileBackup("test.txt.bak", ".bak") == "test.txt.1.bak"
//	IncrementFileBackup("test.txt.1.bak", ".bak") == "test.txt.2.bak"
func IncrementFileBackup(name, extension string) (modified string) {
	var trimmed string
	if extension == "" {
		extension = "~"
	}
	trimmed = strings.TrimSuffix(name, extension)
	if trimmed == name {
		modified = name + extension
		return
	}
	modified = IncrementFilePath(trimmed) + extension
	return
}

// BackupName is a more flexible version of IncrementFileBackup,
// with a variable length list of values configuring the backup
// file naming process. The following is a breakdown of the `argv`
// patterns based on the number of arguments in the `argv`:
//
//	defaults: ext="~", leader=".", prefix="", suffix="", inc=1
//
//	0: all defaults used
//	1: ext=argv[0].(string)
//	2: 1 + leader=argv[1].(string)
//	3: 2 + inc=argv[2].(int)
//	4: 2 + prefix=argv[2].(string) + suffix=argv[3].(string)
//	5: 4 + inc=argv[4].(int)
//
// Examples:
//
//   - zero argv
//     BackupName("file.txt") == "file.txt~"
//     BackupName("file.txt~") == "file.txt.1~"
//   - one argv
//     BackupName("file.txt", ".bak") == "file.txt.bak"
//     BackupName("file.txt.bak", ".bak") == "file.txt.1.bak"
//   - two argv
//     BackupName("file.txt", ".bak", "~") == "file.txt~1.bak"
//     BackupName("file.txt~1.bak", ".bak", "~") == "file.txt~2.bak"
//   - three argv
//     BackupName("file.txt", ".bak", "~", 2) == "file.txt~1.bak"
//     BackupName("file.txt~1.bak", ".bak", "~", 2) == "file.txt~3.bak"
//   - four argv
//     BackupName("file.txt", ".bak", "~", "[", "]") == "file.txt.bak"
//     BackupName("file.txt.bak", ".bak", "~", "[", "]") == "file.txt~[1].bak"
//   - five argv
//     BackupName("file.txt", ".bak", "~", "[", "]", 2) == "file.txt.bak"
//     BackupName("file.txt~[1].bak", ".bak", "~", "[", "]", 2) == "file.txt~[3].bak"
func BackupName(name string, argv ...interface{}) (modified string) {
	extension, leader, prefix, suffix := "~", ".", "", ""
	increment := 1

	switch len(argv) {
	case 1:
		extension, _ = argv[0].(string)
	case 2:
		extension, _ = argv[0].(string)
		leader, _ = argv[1].(string)
	case 3:
		extension, _ = argv[0].(string)
		leader, _ = argv[1].(string)
		increment, _ = argv[2].(int)
	case 4:
		extension, _ = argv[0].(string)
		leader, _ = argv[1].(string)
		prefix, _ = argv[2].(string)
		suffix, _ = argv[3].(string)
	case 5:
		extension, _ = argv[0].(string)
		leader, _ = argv[1].(string)
		prefix, _ = argv[2].(string)
		suffix, _ = argv[3].(string)
		increment, _ = argv[4].(int)
	}

	label := strings.TrimSuffix(name, extension)
	var initial int
	if increment <= 0 {
		// first increment is actually a decrement, and we don't do
		// negative numbered backups, so the initial increment is 1
		initial = 1
	} else {
		// first increment is the same as the increment setting
		initial = increment
	}
	if v, i := IncrementLabel(label, leader, prefix, suffix, increment); i == initial && label == name {
		// backup extension was not present and this is the first increment
		modified = label + extension
	} else {
		modified = v + extension
	}
	return
}