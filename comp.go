/*
 *	File:		src/gitHub.com/Ken1JF/ahgo/ah/comp.go
 *  Project:	abst-hier
 *
 *  Created by Ken Friedenbach on 3/03/2010.
 *  Copyright 2010-2014 Ken Friedenbach. All rights reserved.
 *
 *	This file implements the state computation functions,
 *  for Abstraction Hierarchies.
 *
 */

package ah

import (
	"fmt"
)

func computeNewAreaState(abhr *AbstHier, gl GraphLevel, n NodeLoc, newLoSt uint16) (ret uint16) {
	ret = newLoSt
	return ret
}

func computeAreaHighState(abhr *AbstHier, gl GraphLevel, n NodeLoc, newLoSt uint16) (ret uint16) {
	ret = newLoSt
	return ret
}

func computeNewRegionState(abhr *AbstHier, gl GraphLevel, n NodeLoc, newLoSt uint16) (ret uint16) {
	ret = newLoSt
	return ret
}

func computeRegionHighState(abhr *AbstHier, gl GraphLevel, n NodeLoc, newLoSt uint16) (ret uint16) {
	ret = newLoSt
	return ret
}

func computeNewGroupState(abhr *AbstHier, gl GraphLevel, n NodeLoc, newLoSt uint16) (ret uint16) {
	ret = newLoSt
	return ret
}

func computeGroupHighState(abhr *AbstHier, gl GraphLevel, n NodeLoc, newLoSt uint16) (ret uint16) {
	ret = newLoSt
	return ret
}

func computeNewBlockState(abhr *AbstHier, gl GraphLevel, n NodeLoc, newLoSt uint16) (ret uint16) {
	ret = newLoSt
	return ret
}

func computeBlockHighState(abhr *AbstHier, gl GraphLevel, n NodeLoc, newLoSt uint16) (ret uint16) {
	ret = newLoSt
	return ret
}

func computeNewStringState(abhr *AbstHier, gl GraphLevel, n NodeLoc, newLoSt uint16) (ret uint16) {
	nwSt := StringStatus(newLoSt)
	switch nwSt {
	case BlackString:
		ret = uint16(BlackString)
	case WhiteString:
		ret = uint16(WhiteString)
	default:
		ret = uint16(UnoccupiedString)
	}
	return ret
}

func computeStringHighState(abhr *AbstHier, gl GraphLevel, n NodeLoc, newLoSt uint16) (ret uint16) {
	nwSt := StringStatus(newLoSt)
	switch nwSt {
	case BlackString:
		ret = uint16(BlackString)
	case WhiteString:
		ret = uint16(WhiteString)
	default:
		ret = uint16(UnoccupiedString)
	}
	return ret
}

// PointStatus is either the color of a stone occupying the point,
// or a count/pattern of the stones adjacent to an unoccupied point.
// Liberty points are made unique, by shifting the Col and Row values
// into the upper 11 bits, leaving 5 bits for "raw" PointStatus values.
type PointStatus uint16

// Color and Liberty values
const (
	// Undefined, used by ChangeNodeState, as a temporary value

	UndefinedPointStatus PointStatus = iota

	// Occupied, stone color

	Black
	White

	// Move types, for AB, AW, and AE properties:
	// Note: the "no-op" moves: AB_B, AW_W, and AE_U are not supported.

	AB_U // Add Black, was Unoccupied
	AB_W // Add Black, was White
	AE_B // Add Empty, was Black
	AE_W // Add Empty, was White
	AW_B // Add White, was Black
	AW_U // Add White, was Unoccupied

	// Unoccupied if PointStatus >= Unocc

	Unocc // generic unoccupied, unknown adj. status

	B0W0 // no adj. stones, known to be a non-liberty point

	// Liberty if PointStatus > B0W0

	// Single Adjacent Stone:
	W1
	B1

	// Two Adjacent Stones:
	W2
	B1W1
	W1B1 // non-canonical value
	B2

	// Three Adjacent Stones:
	B3 // 3 adj. stones
	B2W1
	B1W2
	W3
	WBB // non-canonical value
	WBW // non-canonical value
	BWB // non-canonical value
	W2B1

	// Four Adjacent Stones:
	B4
	B3W1
	BWBW
	BBWW
	B1W3
	W4

	LastPointStatus // for checking values
)

// String() string function for fmt.Printf
func (ptst PointStatus) String() string {
	switch ptst {
	case UndefinedPointStatus:
		return "UndefinedPointStatus"
	case Black:
		return "Black"
	case White:
		return "White"

		// Move types, for AB, AW, and AE properties:
		// Note: the "no-op" moves: AB_B, AW_W, and AE_U are not supported.

	case AB_U:
		return "AB_U"
	case AB_W:
		return "AB_W"
	case AE_B:
		return "AE_B"
	case AE_W:
		return "AE_W"
	case AW_B:
		return "AW_B"
	case AW_U:
		return "AW_U"

		// Unoccupied if PointStatus >= Unocc

	case Unocc:
		return "Unocc"

	case B0W0:
		return "B0W0"

		// Liberty if PointStatus > B0W0

		// Single Adjacent Stone:
	case W1:
		return "W1"
	case B1:
		return "B1"

		// Two Adjacent Stones:
	case W2:
		return "W2"
	case B1W1:
		return "B1W1"
	case W1B1:
		return "W1B1"
	case B2:
		return "B2"

		// Three Adjacent Stones:
	case B3:
		return "B3"
	case B2W1:
		return "B2W1"
	case B1W2:
		return "B1W2"
	case W3:
		return "W3"
	case WBB:
		return "WBB"
	case WBW:
		return "WBW"
	case BWB:
		return "BWB"
	case W2B1:
		return "W2B1"

		// Four Adjacent Stones:
	case B4:
		return "B4"
	case B3W1:
		return "B3W1"
	case BWBW:
		return "BWBW"
	case BBWW:
		return "BBWW"
	case B1W3:
		return "B1W3"
	case W4:
		return "W4"

	case LastPointStatus:
		return "LastPointStatus"
	}
	return fmt.Sprintf("UnknownPointSatus: %d", int(ptst))
}

// no computeNewPointState, no new points

func computePointHighState(abhr *AbstHier, gl GraphLevel, nl NodeLoc, newLoSt uint16) uint16 {
	var ret PointStatus
	switch PointStatus(newLoSt) {
	case Black:
		ret = Black
	case White:
		ret = White
	default: // unoccupied
		ret = Unocc
		g := &abhr.Graphs[PointLevel]
		abhr.EachAdjNode(PointLevel, nl,
			func(adjNl NodeLoc) {
				adjN := &g.Nodes[adjNl]
				switch ret {
				// No adjacent stones:
				case Unocc, B0W0:
					switch PointStatus(adjN.GetNodeLowState()) {
					case Black:
						ret = B1
					case White:
						ret = W1
					default:
						ret = B0W0
					}
				// One adjacent stone:
				case B1:
					switch PointStatus(adjN.GetNodeLowState()) {
					case Black:
						ret = B2
					case White:
						ret = B1W1
					default:
						ret = B1
					}
				case W1:
					switch PointStatus(adjN.GetNodeLowState()) {
					case Black:
						ret = W1B1
					case White:
						ret = W2
					default:
						ret = W1
					}
				// Two adjacent stones:
				case B2:
					switch PointStatus(adjN.GetNodeLowState()) {
					case Black:
						ret = B3
					case White:
						ret = B2W1
					default:
						ret = B2
					}
				case B1W1:
					switch PointStatus(adjN.GetNodeLowState()) {
					case Black:
						ret = BWB
					case White:
						ret = B1W2
					default:
						ret = B1W1
					}
				case W1B1:
					switch PointStatus(adjN.GetNodeLowState()) {
					case Black:
						ret = WBB
					case White:
						ret = WBW
					default:
						ret = W1B1
					}
				case W2:
					switch PointStatus(adjN.GetNodeLowState()) {
					case Black:
						ret = W2B1
					case White:
						ret = W3
					default:
						ret = W2
					}
				// Three adjacent stones:
				case B3:
					switch PointStatus(adjN.GetNodeLowState()) {
					case Black:
						ret = B4
					case White:
						ret = B3W1
					default:
						ret = B3
					}
				case B2W1:
					switch PointStatus(adjN.GetNodeLowState()) {
					case Black:
						ret = B3W1
					case White:
						ret = BBWW
					default:
						ret = B2W1
					}
				case BWB:
					switch PointStatus(adjN.GetNodeLowState()) {
					case Black:
						ret = B3W1
					case White:
						ret = BWBW
					default:
						ret = BWB
					}
				case WBB:
					switch PointStatus(adjN.GetNodeLowState()) {
					case Black:
						ret = B3W1
					case White:
						ret = BBWW
					default:
						ret = WBB
					}
				case B1W2:
					switch PointStatus(adjN.GetNodeLowState()) {
					case Black:
						ret = BBWW
					case White:
						ret = B1W3
					default:
						ret = B1W2
					}
				case W2B1:
					switch PointStatus(adjN.GetNodeLowState()) {
					case Black:
						ret = BBWW
					case White:
						ret = B1W3
					default:
						ret = W2B1
					}
				case WBW:
					switch PointStatus(adjN.GetNodeLowState()) {
					case Black:
						ret = BWBW
					case White:
						ret = B1W3
					default:
						ret = B1W2
					}
				case W3:
					switch PointStatus(adjN.GetNodeLowState()) {
					case Black:
						ret = B1W3
					case White:
						ret = W4
					default:
						ret = W3
					}
					// Four adjacent stones: (all done...)
				}
			})
		if ret == B0W0 { // unoccupied, non-Liberty
			// Make the B0W0 unique up to PointType
			// TODO: look at adjacencies, and allow HoshiPt and AdjHoshi
			// to "break up" the unoccupied board regions.
			// TODO: recognize the 1-2, 2-3, 3-4, etc. point types (as part
			// of static board intiialization.)
			typ := abhr.Graphs[PointLevel].Nodes[nl].GetPointType()
			ret = (PointStatus(typ) << UnoccPtTypeShift) | ret
		} else { // Liberty
			// adjust the non-canonical values
			switch ret {
			case BWB:
				ret = B2W1
			case W1B1:
				ret = B1W1
			case W2B1:
				ret = B1W2
			case WBB:
				ret = B2W1
			case WBW:
				ret = B1W2
			default:
			}
			// make the Liberty return value unique:
			ret = PointStatus((int(nl) << UnoccPtTypeShift) | int(ret))
		}
	}
	return uint16(ret)
}

const (
	RawStatusMask    PointStatus = 0x1F // last five bits are raw value
	LibertyColShift  uint16      = 10
	UnoccPtTypeShift uint16      = 6
	LibertyRowShift  uint16      = 5
)

// GetRawPointStatus removes Col and Row bits
func GetRawPointStatus(pst PointStatus) PointStatus {
	return (pst & RawStatusMask)
}

// IsOccupied returns true if PointStatus is Black or White
func IsOccupied(c PointStatus) bool {
	return (c == Black) || (c == White)
}

// IsLiberty returns true if PointStatus is > B0W0
func IsLiberty(c PointStatus) bool {
	return c > B0W0
}

// CurrentColor returns the color of the point after this move is made.
func CurrentColor(mTyp PointStatus) (ret PointStatus) {
	if mTyp == Black {
		ret = Black
	} else if mTyp == White {
		ret = White
	} else {
		switch mTyp {
		case AB_U, AB_W:
			ret = Black
		case AE_B, AE_W:
			ret = Unocc
		case AW_B, AW_U:
			ret = White
		default:
			ret = Unocc // ??? The "no-op" move types are the only ones not explicitly represented
		}
	}
	return ret
}

// PreviousColor returns the color of the point before the move was made.
func PreviousColor(mTyp PointStatus) (ret PointStatus) {
	if mTyp == Black {
		ret = Unocc
	} else if mTyp == White {
		ret = Unocc
	} else {
		switch mTyp {
		case AB_U, AW_U:
			ret = Unocc
		case AB_W, AE_W:
			ret = White
		case AE_B, AW_B:
			ret = Black
		default:
			ret = CurrentColor(mTyp) // The "no-op" move types are the only ones not explicitly represented
		}
	}
	return ret
}

// OppositeColor returns the opposite color for Black and White.
// For other values, the generic Unocc is returned.
func OppositeColor(c PointStatus) (ret PointStatus) {
	if c == Black {
		ret = White
	} else if c == White {
		ret = Black
	} else {
		ret = Unocc
	}
	return ret
}

// InitAbstHier must be called before a abstraction hierarchy can be used
func (abhr *AbstHier) InitAbstHier(c ColValue, r RowValue, upLev GraphLevel, doPlay bool) *AbstHier {
	defer un(trace("InitAbstHier", 0, NilNodeLoc, 0, NilNodeLoc, nilArc, 0xFFFF))
	// turn off tracing during InitAbstHi
	saveTrace := TraceAH
	TraceAH = false
	abhr = abhr.setupAbstHier(c, r, upLev, doPlay,
		computePointHighState, nil, // TODO: or would it be useful for initialization?
		computeStringHighState, computeNewStringState,
		computeBlockHighState, computeNewBlockState,
		computeGroupHighState, computeNewGroupState,
		computeRegionHighState, computeNewRegionState,
		computeAreaHighState, computeNewAreaState)
	if TraceAH {
		abhr.PrintAbstHier("after setupAbstHier", true)
	}
	TraceAH = saveTrace
	return abhr
}
