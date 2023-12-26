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
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
)

func Permissions(path string) (perms fs.FileMode, err error) {
	var st os.FileInfo
	if st, err = os.Stat(path); err == nil {
		perms = st.Mode()
	}
	return
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func ReadFile(path string) (contents string, err error) {
	var raw []byte
	if raw, err = ioutil.ReadFile(path); err == nil {
		contents = string(raw)
	}
	return
}

func WriteFile(path, contents string, perms fs.FileMode) (err error) {
	err = ioutil.WriteFile(path, []byte(contents), perms)
	return
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil && info.IsDir()
}

func IsHidden(path string) bool {
	name := filepath.Base(path)
	return len(name) > 0 && name[0] == '.'
}

func Ls(path string, all bool, recursive bool) (paths []string) {
	if !IsDir(path) {
		paths = append(paths, path)
		return
	}
	if recursive {
		_ = filepath.Walk(path, func(p string, info os.FileInfo, e error) error {
			if e == nil && !info.IsDir() {
				if all || !IsHidden(info.Name()) {
					paths = append(paths, p)
				}
			}
			return nil
		})
		return
	}
	if entries, err := os.ReadDir(path); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				name := entry.Name()
				if all || !IsHidden(name) {
					paths = append(paths, path+string(os.PathSeparator)+name)
				}
			}
		}
	}
	return
}

func Overwrite(path, content string) (err error) {
	err = OverwriteWithPerms(path, content, 0644)
	return
}

func OverwriteWithPerms(path, content string, perm fs.FileMode) (err error) {
	err = ioutil.WriteFile(path, []byte(content), perm)
	return
}

func BackupAndOverwrite(path, backup, content string) (err error) {
	var perms os.FileMode
	if perms, err = Permissions(path); err != nil {
		return
	}
	if err = CopyFile(path, backup); err != nil {
		return
	}
	err = OverwriteWithPerms(path, content, perms)
	return
}

// Diff returns a unified diff comparing two files
func Diff(src, dst string) (unified string, err error) {
	var edits []gotextdiff.TextEdit
	var source string
	if !Exists(src) || IsDir(src) {
		err = fmt.Errorf(`"%v" not found or not a file`, src)
		return
	}
	if !Exists(dst) || IsDir(dst) {
		err = fmt.Errorf(`"%v" not found or not a file`, dst)
		return
	}
	var modified string
	if source, err = ReadFile(src); err != nil {
		return
	}
	if modified, err = ReadFile(dst); err != nil {
		return
	}
	edits = myers.ComputeEdits(span.URIFromPath(src), source, modified)
	unified = fmt.Sprint(gotextdiff.ToUnified(src, dst, source, edits))
	return
}

func CopyFile(src, dst string) (err error) {
	if !Exists(src) || IsDir(src) {
		return fmt.Errorf(`"%v" not found or not a file`, src)
	}
	var srcFile, dstFile *os.File
	if srcFile, err = os.Open(src); err != nil {
		return fmt.Errorf("error opening file: %s", err)
	}
	if dstFile, err = os.Create(dst); err != nil {
		srcFile.Close()
		return fmt.Errorf("error creating file: %s", err)
	}
	defer dstFile.Close()
	defer srcFile.Close()
	if _, e := io.Copy(dstFile, srcFile); e != nil {
		return fmt.Errorf("error copying file: %s", e)
	}
	return nil
}

func MoveFile(src, dst string) (err error) {
	if !Exists(src) || IsDir(src) {
		return fmt.Errorf(`"%v" not found or not a file`, src)
	}
	if err = os.Rename(src, dst); err == nil {
		// rename worked, no need to copy
		return
	}
	if err = CopyFile(src, dst); err != nil {
		return
	}
	err = os.Remove(src)
	if err != nil {
		return fmt.Errorf("error removing file: %s", err)
	}
	return nil
}