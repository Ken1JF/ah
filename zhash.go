/*
 *	File:		src/github.com/Ken1JF/ah/zhash.go
 *  Project:	abst-hier
 *
 *  Created by Ken Friedenbach on 12/24/2010.
 *  Copyright 2010-2014 Ken Friedenbach. All rights reserved.
 *
 *	This file supports Zobrist hashing of Go boards.
 *
 *	Currently uses a primary hash code of 32 bits, and no secondary hash code.
 *	If too many clashss are encountered, the following could be tried:
 *		Use a primary hash code of 64-N bits.
 *		Use a secondary code of N bits.
 *	For N=19, the secondary code could be XOR of B and W stones on each row/col.
 */

package ah

import (
	"math/rand"
)

type ZobristCode uint32

var ZSeed int64

var BlackZKey [MaxBoardSize][MaxBoardSize]ZobristCode
var WhiteZKey [MaxBoardSize][MaxBoardSize]ZobristCode
var KoZKey [MaxBoardSize][MaxBoardSize]ZobristCode

// var BlackToPlayZKey ZobristCode don't need, redundant
var WhiteToPlayZKey ZobristCode

// should be private, but printed in zhash_test.go
var KeyCount int64

var checkKeys map[ZobristCode]int64

const checkKeysNeeded int = 1085

func newZobristKey() ZobristCode {
	var nzc ZobristCode
	var present bool
	present = true
	for present {
		nzc = ZobristCode((rand.Uint32() << 24) ^ (rand.Uint32() << 16) ^ (rand.Uint32() << 8) ^ rand.Uint32())
		_, present = checkKeys[nzc]
		if present {
			panic("Zobrist Duplicate Keys")
		}
	}
	KeyCount += 1
	checkKeys[nzc] = KeyCount
	return nzc
}

// Initialize the Zobrist Keys
func init() {
	checkKeys = make(map[ZobristCode]int64, checkKeysNeeded)
	ZSeed = (2 * 3 * 5 * 7 * 11) - 1
	rand.Seed(ZSeed)
	var i, j uint8
	for i = 0; i < MaxBoardSize; i += 1 {
		for j = 0; j < MaxBoardSize; j += 1 {
			BlackZKey[i][j] = newZobristKey()
			WhiteZKey[i][j] = newZobristKey()
			KoZKey[i][j] = newZobristKey()
		}
	}
	//	BlackToPlayZKey = newZobristKey()
	WhiteToPlayZKey = newZobristKey()
}
