/*
 *	File:		src/gitHub.com/Ken1JF/ah//gostrings.go
 *  Project:	abst-hier
 *
 *  Created by Ken Friedenbach on 03/28/2010.
 *  Copyright 2010-2014 Ken Friedenbach. All rights reserved.
 *
 *	This file supports the strings level of the Go Abstraction Hierarchy.
 *
 */

package ah

type StringStatus uint16

const (
	UndefinedStringStatus StringStatus = iota
	BlackString
	WhiteString
	UnoccupiedString
)
