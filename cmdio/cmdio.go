// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cmdio

import (
	"fmt"
	"io"
)

type CmdIO interface {
	Out(s string)
	Outf(format string, args ...interface{})

	Stdout() io.Writer
}

func NewCmdIO(stdout io.Writer) CmdIO {
	return &cmdIO{
		stdout: stdout,
	}
}

type cmdIO struct {
	stdout io.Writer
}

func (cmdIO *cmdIO) Out(s string) {
	_, _ = io.WriteString(cmdIO.stdout, s)
}

func (cmdIO *cmdIO) Outf(format string, args ...interface{}) {
	_, _ = io.WriteString(cmdIO.stdout, fmt.Sprintf(format, args...))
}

func (cmdIO *cmdIO) Stdout() io.Writer {
	return cmdIO.stdout
}
