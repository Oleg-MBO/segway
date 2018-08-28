package client

import (
	"fmt"
)

// ErrParceData is user when
type ErrParceData struct {
	command string
	err     error
}

func (err *ErrParceData) Error() string {
	return fmt.Sprintf(`Error during parse data of command "%s" err:%#v`, err.command, err.err)
}

// ErrUndefinnedCommand error is used when command from esp
type ErrUndefinnedCommand struct {
	command string
}

func (err *ErrUndefinnedCommand) Error() string {
	return fmt.Sprintf(`undefinded command: %s`, err.command)
}

// ErrEspIsDone is used when driver finish read commands
type ErrEspIsDone struct {
}

func (err *ErrEspIsDone) Error() string {
	return "Esp is done"
}
