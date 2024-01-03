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
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Run is a wrapper around With configured with the given Options and the
// default os.Environ
func Run(path, name string, argv ...string) (stdout, stderr string, status int, err error) {
	return With(Options{
		Path:    path,
		Name:    name,
		Argv:    argv,
		Environ: os.Environ(),
	})
}

// With is a blocking function which runs a command with the Options given and
// if there is an error, looks for the last non-empty line of output to STDERR
// (and if that's empty, checks STDOUT) and returns that as the function `err`
// return value. If both STDERR and STDOUT are empty, the error message is:
// "exit status %d" where the `%d` is replaced with the status code.
func With(options Options) (stdout, stderr string, status int, err error) {
	cmd := exec.Command(options.Name, options.Argv...)
	cmd.Stdin = nil
	cmd.Dir = options.Path
	cmd.Env = options.Environ

	var ob, eb bytes.Buffer
	cmd.Stdout = &ob
	cmd.Stderr = &eb

	err = cmd.Run()

	stdout = ob.String()
	stderr = eb.String()

	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			status = exitError.ExitCode() // always non-zero

			// helper func to find the last, non-empty, line of the given text
			lastLine := func(text string) (errMsg string) {
				if text = strings.TrimSpace(text); text != "" {
					lines := strings.Split(text, "\n")
					for idx := len(lines) - 1; idx >= 0; idx-- {
						if msg := strings.TrimSpace(lines[idx]); msg != "" {
							errMsg = msg
							return
						}
					}
				}
				return
			}

			// check stderr and then stdout for a last line to use as the
			// error message
			if errMsg := lastLine(stderr); errMsg != "" {
				err = errors.New(errMsg)
				return
			} else if errMsg = lastLine(stdout); errMsg != "" {
				err = errors.New(errMsg)
				return
			}

			// both stderr and stdout had no messaging, use a sane fallback
			// message instead
			err = fmt.Errorf("exit status %d", status)
		}
	}
	return
}