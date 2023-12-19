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
	"os"
	"os/exec"
	"strings"
)

func Run(path, name string, argv ...string) (stdout, stderr string, status int, err error) {
	return With(Options{
		Path:    path,
		Name:    name,
		Argv:    argv,
		Environ: os.Environ(),
	})
}

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
			if status = exitError.ExitCode(); status != 0 {
				if text := strings.TrimSpace(stderr); text != "" {
					lines := strings.Split(text, "\n")
					if count := len(lines); count > 0 {
						err = errors.New(lines[count-1])
					} else if text = strings.TrimSpace(stdout); text != "" {
						lines = strings.Split(text, "\n")
						count = len(lines)
						err = errors.New(lines[count-1])
					} else {
						err = errors.New("(unspecified)")
					}
				}
				return
			}
		}
	}
	return
}