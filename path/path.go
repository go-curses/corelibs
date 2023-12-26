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
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"

	"github.com/go-curses/corelibs/strings"
)

var (
	DefaultDirMode  os.FileMode = 0770
	DefaultFileMode os.FileMode = 0660
)

func Permissions(path string) (perms fs.FileMode, err error) {
	var st os.FileInfo
	if st, err = os.Stat(path); err == nil {
		perms = st.Mode()
	}
	return
}

func IsHidden(path string) bool {
	name := filepath.Base(path)
	return len(name) > 0 && name[0] == '.'
}

func IsPlainText(src string) (isPlain bool) {
	if IsFile(src) {
		if kind, err := mimetype.DetectFile(src); err == nil {
			if isPlain = kind.Is("text/plain"); isPlain {
				return
			}
			for parent := kind.Parent(); parent != nil; parent = parent.Parent() {
				if isPlain = parent.Is("text/plain"); isPlain {
					return
				}
			}
		}
	}
	return
}

func Overwrite(path, content string) (err error) {
	var perms os.FileMode
	if perm, ee := Permissions(path); ee != nil {
		perm = DefaultFileMode
	} else {
		perms = perm
	}
	err = OverwriteWithPerms(path, content, perms)
	return
}

func OverwriteWithPerms(path, content string, perm fs.FileMode) (err error) {
	err = os.WriteFile(path, []byte(content), perm)
	return
}

func BackupAndOverwrite(path, backup, content string) (err error) {
	var perms os.FileMode
	if perms, err = Permissions(path); err != nil {
		return
	}

	for Exists(backup) {
		backup = strings.IncrementFilePath(backup)
	}

	if _, err = CopyFile(path, backup); err != nil {
		return
	}
	err = OverwriteWithPerms(path, content, perms)
	return
}

// Diff returns a unified diff comparing two files
func Diff(src, dst string) (unified string, err error) {
	var edits []gotextdiff.TextEdit
	if !Exists(src) || IsDir(src) {
		err = fmt.Errorf(`"%v" not found or not a file`, src)
		return
	} else if !Exists(dst) || IsDir(dst) {
		err = fmt.Errorf(`"%v" not found or not a file`, dst)
		return
	}
	var source, modified []byte
	if source, err = ReadFile(src); err != nil {
		return
	} else if modified, err = ReadFile(dst); err != nil {
		return
	}
	edits = myers.ComputeEdits(span.URIFromPath(src), string(source), string(modified))
	unified = fmt.Sprint(gotextdiff.ToUnified(src, dst, string(source), edits))
	return
}

func MoveFile(src, dst string) (err error) {
	if !Exists(src) || IsDir(src) {
		err = fmt.Errorf(`file not found or is not a regular file`)
		return
	} else if err = os.Rename(src, dst); err == nil {
		// rename worked, no need to copy+remove
		return
	} else if _, err = CopyFile(src, dst); err != nil {
		return
	} else if err = os.Remove(src); err != nil {
		err = fmt.Errorf("error removing old file: %w", err)
		return
	}
	return
}