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
	"os"
)

var (
	DefaultDirMode  = os.FileMode(0770)
	DefaultFileMode = os.FileMode(0660)
)

// Permissions is a wrapper around os.Stat, returning just the
// stat.Mode().Perm() value and any error from the stat call
func Permissions(path string) (perms fs.FileMode, err error) {
	var st os.FileInfo
	if st, err = os.Stat(path); err == nil {
		perms = st.Mode().Perm()
	}
	return
}

// Overwrite overwrites the given file, preserving existing
// permissions
func Overwrite(path, content string) (err error) {
	var perms os.FileMode
	if perm, ee := Permissions(path); ee != nil {
		perms = DefaultFileMode
	} else {
		perms = perm
	}
	err = OverwriteWithPerms(path, content, perms)
	return
}

// OverwriteWithPerms overwrites the given file and ensures the permissions
// are specifically the perms given. Normal unix umask may prevent correct
// permissions, use: `old := syscall.Umask(0); defer syscall.Umask(old)` to
// guarantee the specified perms
func OverwriteWithPerms(path, content string, perms fs.FileMode) (err error) {
	var stat os.FileInfo
	if err = os.WriteFile(path, []byte(content), perms); err != nil {
		return
	} else if stat, err = os.Stat(path); err == nil && stat.Mode().Perm() != perms {
		err = os.Chmod(path, perms)
	}
	return
}

// BackupAndOverwrite uses BackupName (passing argv along) to derive a
// non-existing file name, uses CopyFile to back up the current file
// contents and then uses Overwrite to update the original file contents.
func BackupAndOverwrite(path, content string, argv ...interface{}) (backup string, err error) {
	for backup = BackupName(path, argv...); Exists(backup); {
		backup = BackupName(backup, argv...)
	}

	if _, err = CopyFile(path, backup); err == nil {
		err = Overwrite(path, content)
	}
	return
}

// MoveFile tries to rename `src` to `dst` and if that works, nothing else is
// done. If renaming did not work, MoveFile copies `src` to `dst` and if
// successful, removes the original file. MoveFile can only move regular files
// and not pipes, char devices, etc.
func MoveFile(src, dst string) (err error) {
	src, _ = Abs(src)
	dst, _ = Abs(dst)
	if src == dst {
		// nop, same file
	} else if !Exists(src) {
		err = fmt.Errorf(`file not found`)
	} else if !IsRegularFile(src) {
		err = fmt.Errorf(`not a regular file`)
	} else if err = os.Rename(src, dst); err == nil {
		// rename worked, no need to copy+remove
	} else if _, err = CopyFile(src, dst); err == nil {
		// TODO: figure out how to test CopyFile actually working when
		//       os.Rename fails!
		// copy worked, remove file
		err = os.Remove(src)
	}
	return
}

// PruneEmptyDirs finds all directories starting at the given path, checks if
func PruneEmptyDirs(path string) (err error) {
	var all []string
	if all, err = ListDirs(path, true); err != nil {
		return
	}
	for _, dir := range all {
		var files []string
		if files, err = ListFiles(dir, true); err != nil {
			return
		}
		if len(files) == 0 {
			if err = os.Remove(dir); err != nil {
				return
			}
		}
	}
	return
}

func CopyFile(src, dst string) (copied int64, err error) {
	// see: https://opensource.com/article/18/6/copying-files-go
	var stat os.FileInfo
	if stat, err = os.Stat(src); err != nil {
		return
	} else if !stat.Mode().IsRegular() {
		err = fmt.Errorf("not a regular file")
		return
	}

	var source *os.File
	if source, err = os.Open(src); err != nil {
		return
	}
	defer source.Close()

	var destination *os.File
	if destination, err = os.Create(dst); err != nil {
		return
	}
	defer destination.Close()

	copied, err = io.Copy(destination, source)
	return
}