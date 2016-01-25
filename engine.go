package main

import (
	"bufio"
	"io"
)

// engine is the main program state
type engine struct {
	nxtl     string        // the next line
	pat      string        // the pattern space, possibly nil
	hold     string        // the hold buffer,   possibly nil
	appl     *string       // any lines we've been asked to 'a\'ppend, usually nil
	lastl    bool          // true if it's the last line
	ins      []instruction // the instruction stream
	ip       int           // the current locaiton in the instruction stream
	input    *bufio.Reader // the input stream
	output   *bufio.Writer // the output stream
	lineno   int           // current line number
	modified bool          // have we modified the pattern space?
}

// a sed instruction is mostly a function transforming an engine
type instruction func(*engine) error

// Run executes the instructions until we hit an error.  The most
// common "error" will be io.EOF, which we will translate to nil
func run(e *engine) error {
	var err error

	// prime the engine by filling nxtl... roll back the IP and lineno
	err = cmd_fillNext(e)
	e.ip = 0
	e.lineno = 0

	for err == nil {
		err = e.ins[e.ip](e)
	}

	if err == io.EOF {
		err = nil
	}

	ferr := e.output.Flush() // attempt to flush output
	if ferr != nil && err == nil {
		err = ferr
	}

	return err
}
