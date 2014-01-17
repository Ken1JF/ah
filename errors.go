/*
 *  File:		src/gitHub.com/Ken1JF/ahgo/ah/errors.go
 *  Project:	abst-hier
 *
 *  Created by Ken Friedenbach on 12/08/09.
 *	Copyright 2009-2014, all rights reserved.
 *  Much of this logic is based on the scanner for Go,
 *  whiich may be found in:
 *  		${GOROOT}/src/pkg/go/scanner/
 *
 *  Copyright 2009 The Go Authors. All rights reserved.
 *  Use of this source code is governed by a BSD-style
 *  license that can be found in the LICENSE file.
 */

package ah

import (
	"fmt"
	"io"
	"sort"
)

// Token source positions are represented by a Position value.
// A Position is valid if the line number is > 0.
//
type Position struct {
	Filename string // filename, if any
	Offset   int    // byte offset, starting at 0
	Line     int    // line number, starting at 1
	Column   int    // column number, starting at 1 (character count)
}

// Pos is an accessor method for anonymous Position fields.
// It returns its receiver.
//
func (pos *Position) Pos() Position { return *pos }

// IsValid returns true if the position is valid.
//
func (pos *Position) IsValid() bool { return pos.Line > 0 }

func (pos *Position) String() string {
	s := pos.Filename
	if pos.IsValid() {
		if s != "" {
			s += ":"
		}
		s += fmt.Sprintf("%d:%d", pos.Line, pos.Column)
	}
	if s == "" {
		s = "???"
	}
	return s
}

// NoPos is used when there is no corresponding source position for a token.
var NoPos Position

// In an ErrorList, an error is represented by a *ErrorWithPosition.
// The position Pos, if valid, points to the beginning
// of the offending token, and the error condition is described
// by Msg.
//
type ErrorWithPosition struct {
	Pos Position
	Msg string
}

// ErrorWithPosition implements the error interface
//
func (e ErrorWithPosition) Error() string {
	if e.Pos.Filename != "" || e.Pos.IsValid() {
		// don't print "<unknown position>"
		return e.Pos.String() + ": " + e.Msg
	}
	return e.Msg
}

// ErrorList is a list of *Errors.
// The zero value for an ErrorList is an empty ErrorList ready to use.
//
type ErrorList []*ErrorWithPosition

// Add adds an ErrorWithPosition with given position and error message to an ErrorList.
func (p *ErrorList) Add(pos Position, msg string) {
	*p = append(*p, &ErrorWithPosition{pos, msg})
}

// Reset resets an ErrorList to no errors.
func (p *ErrorList) Reset() { *p = (*p)[0:0] }

// ErrorList implements the sort Interface.
func (p ErrorList) Len() int      { return len(p) }
func (p ErrorList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p ErrorList) Less(i, j int) bool {
	e := &p[i].Pos
	f := &p[j].Pos
	// Note that it is not sufficient to simply compare file offsets because
	// the offsets do not reflect modified line information (through //line
	// comments).
	if e.Filename < f.Filename {
		return true
	}
	if e.Filename == f.Filename {
		if e.Line < f.Line {
			return true
		}
		if e.Line == f.Line {
			return e.Column < f.Column
		}
	}
	return false
}

// Sort sorts an ErrorList. *ErrorWithPosition entries are sorted by position,
// other errors are sorted by error message, and before any *ErrorWithPosition
// entry.
//
func (p ErrorList) Sort() {
	sort.Sort(p)
}

// RemoveMultiples sorts an ErrorList and removes all but the first error per line.
func (p *ErrorList) RemoveMultiples() {
	sort.Sort(p)
	var last Position // initial last.Line is != any legal error line
	i := 0
	for _, e := range *p {
		if e.Pos.Filename != last.Filename || e.Pos.Line != last.Line {
			last = e.Pos
			(*p)[i] = e
			i++
		}
	}
	(*p) = (*p)[0:i]
}

// Return the number of errors
//
func (p ErrorList) ErrorCount() int {
	return len(p)
}

// An ErrorList implements the error interface.
//
func (p ErrorList) Error() string {
	switch len(p) {
	case 0:
		return "no errors"
	case 1:
		return p[0].Error()
	}
	return fmt.Sprintf("%s (and %d more errors)", p[0], len(p)-1)
}

// Err returns an error equivalent to this error list.
// If the list is empty, Err returns nil.
func (p ErrorList) Err() error {
	if len(p) == 0 {
		return nil
	}
	return p
}

// PrintError is a utility function that prints a list of errors to w,
// one error per line, if the err parameter is an ErrorList. Otherwise
// it prints the err string.
//
func PrintError(w io.Writer, err error) {
	if list, ok := err.(ErrorList); ok {
		for _, e := range list {
			fmt.Fprintf(w, "%s\n", e)
		}
	} else {
		fmt.Fprintf(w, "%s\n", err)
	}
}
