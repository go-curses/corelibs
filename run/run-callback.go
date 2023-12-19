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

package run

import (
	"bufio"
	"io"
	"os/exec"

	"github.com/go-curses/cdk"
)

// Callback runs the command defined by the Options and pipes the standard and error output streams to the stdout and
// stderr functions given (line by line)
func Callback(options Options, stdout, stderr func(line string)) (pid int, done chan bool, err error) {
	cmd := exec.Command(options.Name, options.Argv...)
	cmd.Dir = options.Path
	cmd.Stdin = nil
	cmd.Env = options.Environ

	var o, e io.ReadCloser
	if stdout != nil {
		if o, err = cmd.StdoutPipe(); err != nil {
			return
		}
	}
	if stderr != nil {
		if e, err = cmd.StderrPipe(); err != nil {
			return
		}
	}

	if err = cmd.Start(); err != nil {
		return
	}

	pid = cmd.Process.Pid

	if stdout != nil {
		cdk.Go(func() {
			so := bufio.NewScanner(o)
			for so.Scan() {
				stdout(so.Text())
			}
		})
	}

	if stderr != nil {
		cdk.Go(func() {
			se := bufio.NewScanner(e)
			for se.Scan() {
				line := se.Text()
				stderr(line)
			}
		})
	}

	done = make(chan bool)
	cdk.Go(func() {
		if err = cmd.Wait(); err != nil {
			stderr(err.Error())
			done <- false
		} else {
			done <- true
		}
	})
	return
}