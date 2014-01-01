/*
 *	File:       src/gitHub.com/Ken1JF/ahgo/ah/board.go
 *  Project:    ahgo
 *
 *  Created by Ken Friedenbach on 12/08/09.
 *  Copyright 2009-2014 Ken Friedenbach. All rights reserved.
 *
 *	This package implements storage of Go boards,
 *  and provides the base types for Abstraction Hierarchies.
 */

package ah

import (
	"fmt"
	"strconv"
)

// Constants for Board Size and NodeLoc encoding.
//
// The NodeLoc encoding follows the SGF FF[4] convention,
// as illustrated in this diagram:
//		http://www.red-bean.com/sgf/images/TA2.gif
// The first coordinate is the column, the second is the row.
// This ordering of Column, Row is maintained throughout most
// of the code: function arguments, return values, etc.
//
// When the coordinates are used to access points of the Board array,
// the order is reversed:
//		Nodes[nl]
// The r value selects a row, and c selects a column within the row.
//
const (
	MaxBoardSize uint8 = 19
	MaxRows      uint8 = MaxBoardSize
	MaxCols      uint8 = 32 // for efficiency of access
	RowShift     uint8 = 5
	ColMask      int   = 0x1F                                  // lower 5 bits
	MaxListLen   int   = int(MaxBoardSize) * int(MaxBoardSize) // to stop infinite cycles
)

// These Direction definitions refer to the canonical board orientation,
// which may be rotated or reflected for display, etc.
//
type Direction uint8

// These costants provide a basis for classifying points in an N X M board,
// indicating that a point has an adjacent point in the given direction.
//
const (
	NoDir    Direction = iota            // singleton point
	UpperDir Direction = 1 << (iota - 1) // r-1
	LeftDir                              // c-1
	LowerDir                             // r+1
	RightDir                             // c+1
)

// PointType is a static classification of all board points
// used in non-empty Graphs embedded on a rectangular grid.
//
// The basic codes are singleton_pt through CenterPt.
// They are formed by adding together the
// directions in which the point has a connection.
//
// There are also advanced codes (specialized _cetner_pt codes)

type PointType uint8

const (

	// SingletonPt has no connection.
	// It occurs only on a 1 x 1 board.
	//
	SingletonPt PointType = PointType(NoDir)

	// An EndPt has only one connection.
	// They occur on the ends of 1 x N and N x 1 boards.
	//
	LowerEndPt PointType = PointType(UpperDir)
	RightEndPt PointType = PointType(LeftDir)
	UpperEndPt PointType = PointType(LowerDir)
	LeftEndPt  PointType = PointType(RightDir)

	// A BridgePt has two opposite connections.
	// They occur in the internal points of 1 x N and N x 1 boards.
	//
	UpperLowerBridgePt PointType = PointType(LowerDir + UpperDir)
	LeftRightBridgePt  PointType = PointType(RightDir + LeftDir)

	// A CornerPt has two adjacent connections.
	// They occur on N x M boards, both N and M >= 2.
	//
	UpperLeftCornerPt  PointType = PointType(LowerDir + RightDir)
	UpperRightCornerPt PointType = PointType(LowerDir + LeftDir)
	LowerRightCornerPt PointType = PointType(UpperDir + LeftDir)
	LowerLeftCornerPt  PointType = PointType(UpperDir + RightDir)

	// An EdgePt has three adjacent connections.
	// They occur on boards with either N or M >= 3.
	//
	UpperEdgePt PointType = PointType(LeftDir + LowerDir + RightDir)
	LeftEdgePt  PointType = PointType(UpperDir + LowerDir + RightDir)
	LowerEdgePt PointType = PointType(UpperDir + LeftDir + RightDir)
	RightEdgePt PointType = PointType(UpperDir + LeftDir + LowerDir)

	// A CenterPt has four adjacent connections.
	//
	CenterPt PointType = PointType(UpperDir + LeftDir + LowerDir + RightDir)

	// A HoshiPt is a specially marked CenterPt.
	//
	HoshiPt PointType = PointType(CenterPt + (1 << 4))

	// Corner N-N and Line N points
	//
	Corner_2_2_Pt PointType = PointType(CenterPt + (2 << 4))
	Line_2_Pt     PointType = PointType(CenterPt + (3 << 4))
	Corner_3_3_Pt PointType = PointType(CenterPt + (4 << 4))
	Line_3_Pt     PointType = PointType(CenterPt + (5 << 4))
	Corner_4_4_Pt PointType = PointType(CenterPt + (6 << 4))
	Line_4_Pt     PointType = PointType(CenterPt + (7 << 4))
	Corner_5_5_Pt PointType = PointType(CenterPt + (8 << 4))
	Line_5_Pt     PointType = PointType(CenterPt + (9 << 4))
	Corner_6_6_Pt PointType = PointType(CenterPt + (10 << 4))
	Line_6_Pt     PointType = PointType(CenterPt + (11 << 4))
	Corner_7_7_Pt PointType = PointType(CenterPt + (12 << 4))
	Line_7_Pt     PointType = PointType(CenterPt + (13 << 4))
	// TODO: do we need UninitializedPt?
	// Or is SingletonPt the correct "zero" state? (i.e. unconnected)
	//
//	UninitializedPt	PointType = PointType(CenterPt + (14 << 4))
)

var PtTypeNames = [255]string{
	SingletonPt: "·",

	LowerEndPt: "╹",
	RightEndPt: "╸",
	UpperEndPt: "╻",
	LeftEndPt:  "╺",

	UpperLowerBridgePt: "┃",
	LeftRightBridgePt:  "━",

	UpperLeftCornerPt:  "┏",
	UpperRightCornerPt: "┓",
	LowerRightCornerPt: "┛",
	LowerLeftCornerPt:  "┗",

	UpperEdgePt: "┯",
	LeftEdgePt:  "┠",
	LowerEdgePt: "┷",
	RightEdgePt: "┨",

	CenterPt: "╋",

	HoshiPt: "◘",

	Corner_2_2_Pt: "╬",
	Line_2_Pt:     "┼",
	Corner_3_3_Pt: "╬",
	Line_3_Pt:     "┼",
	Corner_4_4_Pt: "╬",
	Line_4_Pt:     "┼",
	Corner_5_5_Pt: "╬",
	Line_5_Pt:     "┼",
	Corner_6_6_Pt: "╬",
	Line_6_Pt:     "┼",
	Corner_7_7_Pt: "╬",
	Line_7_Pt:     "┼",
	// TODO: do we need UninitializedPt?
	// Or is SingletonPt the correct "zero" state? (i.e. unconnected)
	//
	//	UninitializedPt:	"¿",
}

// ColValue and RowValue are types that represent column and row coordinates, respectively.
// By making them different types, type checking can help avoid ordering errors.
//
type ColValue uint8
type RowValue uint8

// TODO:(align) if alignment in 6g compiler is improved, go back to using this:
//
//type LocValue uint16

// NodeLoc is a type overload in Abstractions Hierarchies.
// At the Board level, NodeLoc is an col, row pair corresponding
// to board coordinates.
// At Go String and higher abstractions, NodeLoc is an index into
// an array of Nodes.
//
type NodeLoc uint16

// TODO:(align) if alignment in 6g compiler is improved, go back to using this:
//
//type NodeLoc struct {
//	Location LocValue
//}

type NodeLocList []NodeLoc

// NodeLocList satisfies the sort interface:
//
func (p NodeLocList) Len() int {
	return len(p)
}

func (p NodeLocList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p NodeLocList) Less(i, j int) bool {
	if p[i] < p[j] {
		return true
	}
	return false
}

// GetColRow returns the Col and Row portions of a NodeLoc
//
func GetColRow(loc NodeLoc) (col ColValue, row RowValue) {
	row = RowValue((int(loc) >> RowShift) & ColMask)
	col = ColValue(int(loc) & ColMask)
	return col, row
}

// MakeNodeLoc returns a NodeLoc with the given Col and Row values
//
func MakeNodeLoc(col ColValue, row RowValue) (nl NodeLoc) {
	nl = NodeLoc(((int(row) & ColMask) << RowShift) | (int(col) & ColMask))
	return nl
}

// Some useful global variable values
//
var ( // TODO: redo as const???
	PassNodeLoc    NodeLoc
	IllegalNodeLoc NodeLoc
)

// Initialization function for this file
//
func init() {
	// TODO: change the representation of PassNodeLoc, so we can support larger boards
	PassNodeLoc = MakeNodeLoc(ColValue(MaxBoardSize), RowValue(MaxBoardSize))
	IllegalNodeLoc = MakeNodeLoc(ColValue(MaxBoardSize+1), RowValue(MaxBoardSize+1))
}

// A Move is an action that changes the state of the board or game.
// Enough detail is stored so thet the move can be "undone".
// This allows an array to serve as a stack for use in traversing a tree of Move actions.
// The color of the moveLoc before and after the move is known from the moveType
// And information about captures and kos are stored so they can be restored.
//
type MoveRecord struct {
	moveLoc     NodeLoc     // Location where the change occurs (off board for Pass)
	moveType    PointStatus // Black, White, AB_U, AB_W, AE_B, AE_W, AW_B, and AW_U
	moveNum     uint16      // move number in the game (handicaps and setup actions are 0)
	capturedBy  int16       // index of the move that captures this move
	nextCapture int16       // next move captured by the same move
	firstCap    int16       // first move captured by this move
	koPoint     NodeLoc     // koPoint before this move
}

func (mv *MoveRecord) PrintMove(idx int) {
	if idx >= 0 {
		fmt.Print("Move[", idx, "]:")
	} else {
		fmt.Print("Move:")
	}
	c, r := GetColRow(mv.moveLoc)
	fmt.Print(" Loc: (", c, ",", r, ") Type: ", mv.moveType, " Num:", mv.moveNum)
	if mv.capturedBy != nilMovNum {
		fmt.Print(" CapdBy:", mv.capturedBy)
	}
	if mv.nextCapture != nilMovNum {
		fmt.Print(" NextCap:", mv.nextCapture)
	}
	if mv.firstCap != nilMovNum {
		fmt.Print(" FirstCap:", mv.firstCap)
	}
	if mv.koPoint != NilNodeLoc {
		fmt.Print(" Ko:", mv.koPoint)
	}
	fmt.Println(" ")
}

const nilMovNum int16 = -1

// GetPointType retrieves the static PointType of a Point (stored in inList).
//
func (bp *GraphNode) GetPointType() PointType {
	return PointType(bp.inList)
}

// GetPointColRow retrieves the NodeLoc of a Point (stored in firstMem),
// and returns the column and row components.
//
func (bp *GraphNode) GetPointColRow() (ColValue, RowValue) {
	nl := bp.firstMem
	c, r := GetColRow(nl)
	return c, r
}

// MoveHashList records the move index and stone counts associated with a board position
//
type MoveHashList struct {
	moveIdx        uint16 // index into movs[]
	bCount, wCount uint8  // black and white stone stone count, mode 256
	nextMoveHash   *MoveHashList
}

// Board holds the basic elements to record a Go game
//
type Board struct {
	// hoshi points on the board:
	hoshiPts NodeLocList
	// moves of the game, including setup
	movs []MoveRecord
	// board size
	colSize ColValue
	rowSize RowValue
	// terms of play
	handi uint8
	komi  float32
	// Ko point (point where immediate capture is illegal)
	KoPt NodeLoc
	// Who made the first move?
	mov1    PointStatus
	mov1Set bool
	// the number of moves, NOT the number of entries in movs[], 
	// which includes handicap and setup stones.
	numMoves        uint16
	numbBlackStones uint16
	numbWhiteStones uint16
	// board orientation
	brdTrans BoardTrans
	// board zCode
	brdZCode [NUM_TRANS]ZobristCode
	// check for duplicates:
	zHashTable map[ZobristCode]*MoveHashList
	// TODO: add some others: toMove, rules, players, etc.
	// TODO: is this too much? or not enough?
}

// GetNMoves returns the number of moves currently stored in the movs array
//
func (brd *Board) GetNMoves() int {
	return int(brd.numMoves)
}

// setSize sets the size of the board, and calls initBoardPoints
//
func (abhr *AbstHier) setSize(csz ColValue, rsz RowValue) {
	abhr.colSize = csz
	abhr.rowSize = rsz
    for t := T_FIRST; t <= T_LAST; t += 1 {
        abhr.brdZCode[t] = 0   // empty board and BlackToPlayZKey == 0 
    }
	abhr.zHashTable = make(map[ZobristCode]*MoveHashList, 100)
	abhr.KoPt = NilNodeLoc
	abhr.Graphs[PointLevel].initBoardPoints(csz, rsz)
	abhr.initHoshiPts()
}

// addHoshiPt set the Board point, and puts the point in the hoshiPts array.
//
func (abhr *AbstHier) addHoshiPt(c ColValue, r RowValue) {
	nl := MakeNodeLoc(c, r)
	abhr.Graphs[PointLevel].Nodes[nl].inList = ArcIdx(HoshiPt)
	abhr.hoshiPts = append(abhr.hoshiPts, nl)
}

// OnBoard returns true if the point is on the board
//
func (brd *Board) OnBoard(nl NodeLoc) (ret bool) {
	c, r := GetColRow(nl)
        //    fmt.Printf(" c = %d, brd.colSize = %d. r = %d, brd.rowSize = %d\n", c, brd.colSize, r, brd.rowSize)
	if (c < brd.colSize) && (r < brd.rowSize) {
		ret = true
	}
	return ret
}

// GetMovDepth returns the current size of Board.movs
//
func (brd *Board) GetMovDepth() (ret int16) {
	ret = int16(len(brd.movs))
	return ret
}

// AddMove supports a variable sized array
//
func (abhr *AbstHier) AddMove(nl NodeLoc, color PointStatus, n uint16, ko NodeLoc) int {
	ln := len(abhr.movs)
	var newMove MoveRecord
	abhr.movs = append(abhr.movs, newMove)
	abhr.movs[ln].moveLoc = nl
	abhr.movs[ln].moveType = color
	abhr.movs[ln].moveNum = n
	abhr.movs[ln].capturedBy = nilMovNum
	abhr.movs[ln].nextCapture = nilMovNum
	abhr.movs[ln].firstCap = nilMovNum
	abhr.movs[ln].koPoint = ko
    if n > 0 {
        // toggle who to play, always
        for t := T_FIRST; t <= T_LAST; t +=1 {
            abhr.brdZCode[t] = abhr.brdZCode[t] ^ WhiteToPlayZKey
        }
    }
	c, r := GetColRow(nl)
	// record the stone, if one
	if nl != PassNodeLoc {
		if color == Black {
            for t := T_FIRST; t <= T_LAST; t +=1 {
                t_nl := abhr.TransNodeLoc(t, c, r)
                t_c, t_r := GetColRow(t_nl)
                abhr.brdZCode[t] = abhr.brdZCode[t] ^ BlackZKey[t_c][t_r]
            }
			abhr.numbBlackStones += 1
		} else {
            for t := T_FIRST; t <= T_LAST; t +=1 {
                t_nl := abhr.TransNodeLoc(t, c, r)
                t_c, t_r := GetColRow(t_nl)
                abhr.brdZCode[t] = abhr.brdZCode[t] ^ WhiteZKey[t_c][t_r]
            }
			abhr.numbWhiteStones += 1
		}
	}
	// remove the old koPoint, if there was one
	if ln > 0 {
		if abhr.movs[ln-1].koPoint != NilNodeLoc {
            c, r = GetColRow(abhr.movs[ln-1].koPoint)
            for t := T_FIRST; t <= T_LAST; t += 1 {
                t_nl := abhr.TransNodeLoc(t, c, r)
                t_c, t_r := GetColRow(t_nl)
                abhr.brdZCode[t] = abhr.brdZCode[t] ^ KoZKey[t_c][t_r]
            }
		}
	}
	// add the new koPoint, if there is one
	if ko != NilNodeLoc {
		c, r = GetColRow(ko)
        for t := T_FIRST; t <= T_LAST; t += 1 {
            t_nl := abhr.TransNodeLoc(t, c, r)
            t_c, t_r := GetColRow(t_nl)
            abhr.brdZCode[t] = abhr.brdZCode[t] ^ KoZKey[t_c][t_r]
        }
	}
	return ln
}

// SetPoint is used ONLY by test_ahgo.go,
// in setting up test boards for testing trans.go, etc.
// TODO: change test_ahgo.go so SetPoint is not needed...
//
func (abhr *AbstHier) SetPoint(nl NodeLoc, color PointStatus) (err ErrorList) {
	if abhr.OnBoard(nl) {
		abhr.Graphs[PointLevel].Nodes[nl].SetNodeLowState(uint16(color))
		if (color == Black) || (color == White) { // add a move record, for possible future capture
			_ = abhr.AddMove(nl, color, 0, NilNodeLoc) // pre-placed are move "0"
		}
	} else {
		err.Add(NoPos, "point not on board")
	}
	return err
}

// NumLiberties returns the number of vacant points adjacent to a point
//
func (abhr *AbstHier) NumLiberties(nl NodeLoc) int {
	nLibs := 0
	g := &abhr.Graphs[PointLevel]
	hiG := &abhr.Graphs[StringLevel]
	hiNL := g.Nodes[nl].memberOf
	abhr.EachAdjNode(StringLevel, hiNL,
		func(adjNl NodeLoc) {
			if !IsOccupied(PointStatus(g.Nodes[hiG.Nodes[adjNl].firstMem].GetNodeLowState())) {
				nLibs += 1
			}
		})
	return nLibs
}

const RulesAllowSuicide bool = false // TODO: make this a function of Rules being used

// UndoBoardMove is used to move up the game tree.
//
func (abhr *AbstHier) UndoBoardMove(doPlay bool) {
	defer un(trace("UndoBoardMove", PointLevel, NilNodeLoc, 0, NilNodeLoc, nilArc, 0))
	movIdx := len(abhr.movs) - 1

	if movIdx >= 0 {
		//		fmt.Printf("Undoing Move index: %d moveNum\n", movIdx, abhr.movs[movIdx].moveNum)
		// get the move to undo
		var undoMove = abhr.movs[movIdx]
		if TraceAH {
			undoMove.PrintMove(movIdx)
		}
		// remove it from abhr
		abhr.numMoves -= 1
		abhr.movs = abhr.movs[0:movIdx]
		abhr.DecrementMovesPrinted()
		// undo the effects of the move:
		// 1. Retore the previous Ko point
		abhr.KoPt = undoMove.koPoint
		unDoLoc := undoMove.moveLoc
		if doPlay && (unDoLoc != PassNodeLoc) { // nothing to undo for a pass (TODO: excpet for who is to play?)
			var ncaps = 0
			// 2. Restore the move loc to its previous color
			prevColor := PreviousColor(undoMove.moveType)
			if prevColor == Unocc {
				prevColor = PointStatus(abhr.Graphs[PointLevel].CompHigh(abhr, PointLevel, unDoLoc, uint16(prevColor)))
			}
			abhr.ChangeNodeState(PointLevel, unDoLoc, NodeStatus(prevColor), true)
			// 3. Undo captures if there are any:
			cap := undoMove.firstCap
			if cap != nilMovNum {
				// 3A. restore captured stones to opposite color
				oppColor := OppositeColor(undoMove.moveType)
				for cap != nilMovNum {
					ncaps += 1
					abhr.ChangeNodeState(PointLevel, abhr.movs[cap].moveLoc, NodeStatus(oppColor), true)
					cap = abhr.movs[cap].nextCapture
				}
				// 3B. remove the moves from the capture list
				cap = undoMove.firstCap
				for cap != nilMovNum {
					nxtCap := abhr.movs[cap].nextCapture
					abhr.movs[cap].capturedBy = nilMovNum
					abhr.movs[cap].nextCapture = nilMovNum
					abhr.SetBackMovesPrinted(cap)
					cap = nxtCap
				}
			}
		}
	}
	if TraceAH {
		abhr.PrintAbstHier(" When leaving UndoBoardMove", true)
	}
}

// DoBoardMove is used by the SGF parser (sgf: W, B, etc.)
// also used by AddTeachingPattern and other funtions that re-walk the SGF tree.
//
func (abhr *AbstHier) DoBoardMove(nl NodeLoc, color PointStatus, doPlay bool) (moveNum int, err ErrorList) {
	defer un(trace("DoBoardMove", PointLevel, nl, 0, NilNodeLoc, nilArc, NodeStatus(color)))
	var nextKoPoint, firstCapture NodeLoc = NilNodeLoc, NilNodeLoc
	var numCapture, numFriendly, numEnemy, numVacant int

	// Record the color of the first move:
	if abhr.mov1Set == false {
		abhr.mov1Set = true
		abhr.mov1 = color
	}

	// Check the legality of the move

	// Check if NodeLoc is on the board
	if !abhr.OnBoard(nl) { // not on board,
		if nl != PassNodeLoc { //	check if a Pass move
			err.Add(NoPos, "bad move, not on Board")
		} else { // OK, record pass
			abhr.numMoves += 1
			_ = abhr.AddMove(PassNodeLoc, color, abhr.numMoves, NilNodeLoc)
		}
	} else { // OnBoard, check for captures and legality
		pt := &abhr.Graphs[PointLevel].Nodes[nl]
		if IsOccupied(PointStatus(pt.GetNodeLowState())) { // Check if Occupied
			// don't make an error if not playing
			if doPlay {
				err.Add(NoPos, "move at occupied point")
			}
		}
		if nl == abhr.KoPt { // check for illegal Ko recapture
			err.Add(NoPos, "illegal retake of Ko at move " + strconv.Itoa(int(abhr.numMoves+1)))
		}
		// Count Friendly, Enemy, Vacant, and Captures among Adjacent
		abhr.EachAdjNode(PointLevel, nl, // ChkBlack/WhiteAdjs
			func(nl2 NodeLoc) {
				pt2 := abhr.Graphs[PointLevel].Nodes[nl2]
				pt2St := pt2.GetNodeLowState()
				if IsOccupied(PointStatus(pt2St)) {
					if pt2St == uint16(color) {
						numFriendly += 1
					} else {
						numEnemy += 1
						if abhr.NumLiberties(nl2) == 1 {
							numCapture += 1
							if firstCapture == NilNodeLoc {
								firstCapture = nl2
							}
						}
					}
				} else {
					numVacant += 1
				}
			})
		if numVacant == 0 && doPlay { // potential Ko or Suicide, if playing
			if (numCapture == 1) && (numFriendly == 0) {
				if abhr.NumElements(PointLevel, firstCapture) == 1 {
					nextKoPoint = firstCapture
				}
			} else if numCapture == 0 {
				if !RulesAllowSuicide { // TODO: xxx make this a function call
					abhr.EachAdjNode(PointLevel, nl, // ChkBlk/WhiteLibs
						func(nl2 NodeLoc) { // Lower Estimate of gained Liberties
							pt2 := abhr.Graphs[PointLevel].Nodes[nl2]
							pt2St := pt2.GetNodeLowState()
							if pt2St == uint16(color) {
								if abhr.NumLiberties(nl2) > 1 {
									numVacant += 1 // at least!
								}
							}
						})
					if numVacant == 0 {
						err.Add(NoPos, "illegal suicide")
					}
				}
			}
		}
		// record the move, and make it
		abhr.numMoves += 1
		moveNum = int(abhr.numMoves)
		ptMovNum := abhr.AddMove(nl, color, abhr.numMoves, abhr.KoPt)
		ptMov := &abhr.movs[ptMovNum]
		if doPlay {
			abhr.ChangeNodeState(PointLevel, nl, NodeStatus(color), true)
			// Make captures:
			if numCapture > 0 {
				capCount := 0
				if TraceAH {
					fmt.Println(" Starting to make ", numCapture, " captures.")
				}
				capColor := OppositeColor(color)
				abhr.EachAdjNode(PointLevel, nl,
					func(adjNL NodeLoc) {
						adjColor := abhr.Graphs[PointLevel].Nodes[adjNL].lowState
						if (adjColor == uint16(capColor)) && (abhr.NumLiberties(adjNL) == 0) { // string to be captured
							abhr.EachMember(StringLevel, abhr.Graphs[PointLevel].Nodes[adjNL].memberOf,
								// Put the members on the capturedBy list:
								func(capMem NodeLoc) {
									//									capMemPt := &abhr.Graphs[PointLevel].Nodes[capMem]
									capMemMovNum := abhr.FindMoveNumber(capMem)
									if capMemMovNum == nilMovNum {
										err.Add(NoPos, "move number of capture not found")
									} else {
										capMemMov := &abhr.movs[capMemMovNum]
										if capMemMov.capturedBy == nilMovNum {
											capCount += 1
											capMemMov.capturedBy = int16(ptMovNum)
											capMemMov.nextCapture = ptMov.firstCap
											ptMov.firstCap = capMemMovNum
										}
									}
								})
						}
					})
				// Set captured points to Unocc
				if TraceAH {
					fmt.Println(" Setting", capCount, "captured stones to Unocc.")
				}
				capMovNum := abhr.movs[ptMovNum].firstCap
				for capMovNum != nilMovNum {
					abhr.SetBackMovesPrinted(capMovNum)
					capNL := abhr.movs[capMovNum].moveLoc
					c, r := GetColRow(capNL)
					if capColor == Black {
                        for t := T_FIRST; t <= T_LAST; t +=1 {
                            t_nl := abhr.TransNodeLoc(t, c, r)
                            t_c, t_r := GetColRow(t_nl)
                            abhr.brdZCode[t] = abhr.brdZCode[t] ^ BlackZKey[t_c][t_r]
                        }
						abhr.numbBlackStones -= 1
					} else {
                        for t := T_FIRST; t <= T_LAST; t +=1 {
                            t_nl := abhr.TransNodeLoc(t, c, r)
                            t_c, t_r := GetColRow(t_nl)
                            abhr.brdZCode[t] = abhr.brdZCode[t] ^ WhiteZKey[t_c][t_r]
                        }
						abhr.numbWhiteStones -= 1
					}
					abhr.Graphs[PointLevel].Nodes[capNL].lowState = uint16(Unocc)
					capMovNum = abhr.movs[capMovNum].nextCapture
				}
				// Change captured points to new state
				if TraceAH {
					fmt.Println(" Changing", capCount, "captured stones to newstates.")
				}
				capMovNum = abhr.movs[ptMovNum].firstCap
				for capMovNum != nilMovNum {
					capNL := abhr.movs[capMovNum].moveLoc
					newStatus := abhr.Graphs[PointLevel].CompHigh(abhr, PointLevel, capNL, uint16(Unocc))
					abhr.ChangeNodeState(PointLevel, capNL, NodeStatus(newStatus), true)
					capMovNum = abhr.movs[capMovNum].nextCapture
				}
				// Check for changes in adjacent (Liberty) points:
				abhr.EachAdjNode(PointLevel, nl,
					func(adjNl NodeLoc) {
						curSt := abhr.Graphs[PointLevel].Nodes[adjNl].GetNodeHighState()
						lowSt := abhr.Graphs[PointLevel].Nodes[adjNl].GetNodeLowState()
						newSt := abhr.Graphs[PointLevel].CompHigh(abhr, PointLevel, adjNl, lowSt)
						if newSt != curSt {
							abhr.ChangeNodeState(PointLevel, adjNl, NodeStatus(newSt), true)
						}
					})
			}
		}
	}
	abhr.KoPt = nextKoPoint
	// Record this board position, and check if it is a repetition:
	if nl != PassNodeLoc { // don't worry about repetitions after a pass move
		var v *MoveHashList
		v, ok := abhr.zHashTable[abhr.brdZCode[T_IDENTITY]]
		if ok {
			// check if this is a false match:
			if (v.bCount == uint8(abhr.numbBlackStones)) && (v.wCount == uint8(abhr.numbWhiteStones)) {
				err.Add(NoPos, "repeat position, move " + strconv.Itoa(moveNum) + " repeats position after " + strconv.Itoa(int(v.moveIdx)))
			} // else {
			// TODO: remove after checking
			//	err.Add(NoPos, "false repeat position")
			// }
			// find the end of the list
			for v.nextMoveHash != nil {
				v = v.nextMoveHash
			}
			v.nextMoveHash = new(MoveHashList)
			v.nextMoveHash.moveIdx = uint16(moveNum)
			v.bCount = uint8(abhr.numbBlackStones)
			v.wCount = uint8(abhr.numbWhiteStones)
		} else {
			// Record this position
			v = new(MoveHashList)
			v.moveIdx = uint16(moveNum)
			v.nextMoveHash = nil
			v.bCount = uint8(abhr.numbBlackStones)
			v.wCount = uint8(abhr.numbWhiteStones)
//			abhr.zHashTable[abhr.brdZCode] = v, true
            for t := T_FIRST; t <= T_LAST; t +=1 {
                // TODO: how to record transformation?
                abhr.zHashTable[abhr.brdZCode[t]] = v
            }
		}
	}
	return moveNum, err
}

// GetSize returns the column and row sizes of the Board
//
func (abhr *AbstHier) GetSize() (ColValue, RowValue) {
	return abhr.colSize, abhr.rowSize
}

// SetHandicap sets the handicap on the Board
// Note: AB must be used to add the handicap stones.
//
func (brd *Board) SetHandicap(n int) {
	brd.handi = uint8(n)
}

// GetHandicap returns the handicap set for the Board.
//
func (brd *Board) GetHandicap() int {
	return int(brd.handi)
}

// SetKomi sets the komi for the Board.
//
func (brd *Board) SetKomi(k float32) {
	brd.komi = k
}

// GetKomi returns the komi set for the Board.
//
func (brd *Board) GetKomi() float32 {
	return brd.komi
}

// SetMov1 sets the player who moves first.
//
func (brd *Board) SetMov1(c PointStatus) {
	if brd.mov1Set == false {
		brd.mov1Set = true
		brd.mov1 = c
	}
}

// GetMov1 returns the player who moves first.
//
func (brd *Board) GetMov1() (PointStatus, bool) {
	return brd.mov1, brd.mov1Set
}

// SearchStack is type of value returned by Depth First Search.
//
type SearchStack struct {
	nods  []NodeLoc
	g     *Graph
	cur   uint16 // used in Breadth first Search: first pushed, first visited.
	found bool
}

// PushAndMark pushs a node on the search stack, and marks it
//
func (stk *SearchStack) PushAndMark(nl NodeLoc) *SearchStack {
	defer un(trace("PushAndMark", stk.g.gLevel, nl, 0, NilNodeLoc, nilArc, 0xFFFF))
	if !stk.g.Nodes[nl].IsBFSMarked() {
		if TraceAH {
			PrintNodeLoc(stk.g.gLevel, nl, "  Pushing: ")
			fmt.Println()
		}
		ln := len(stk.nods)
		if ln == cap(stk.nods) { // reallocate
			newLen := (ln + 1) * 2
			if newLen < 16 {
				newLen = 16
			}
			newPts := make([]NodeLoc, newLen)
			for i, cur_nod := range stk.nods {
				newPts[i] = cur_nod
			}
			stk.nods = newPts
		}
		stk.nods = stk.nods[0 : ln+1]
		stk.nods[ln] = nl
		stk.g.Nodes[nl].MarkBFSNode()
	} else {
		if TraceAH {
			PrintNodeLoc(stk.g.gLevel, nl, "  Not pushed, already marked: ")
			fmt.Println()
		}
	}
	return stk
}

// UnMarkNodes clears the BFSMark on the nodes in SearchStack
//
func (stk *SearchStack) UnMarkNodes() {
	for _, nl := range stk.nods {
		stk.g.Nodes[nl].ClearBFSMark()
	}
}

// types for iteration and search functions
//
type GraphNodeLocFunc func(*Graph, NodeLoc)
type NodeLocFuncBool func(NodeLoc) bool

// EachNode visits all the graph Nodes.
//
func (abhr *AbstHier) EachNode(gl GraphLevel, Visit GraphNodeLocFunc) {
	if gl == PointLevel {
		var c ColValue
		var r RowValue
		nc, nr := abhr.GetSize()
		for r = 0; r < nr; r++ {
			for c = 0; c < nc; c++ {
				Visit(&abhr.Graphs[PointLevel], MakeNodeLoc(c, r))
			}
		}
	} else {
		g := &abhr.Graphs[gl]
		for i, nod := range g.Nodes {
			if (nod.mark & OnFreeList) == 0 {
				Visit(g, NodeLoc(i))
			}
		}
	}
}

// EachAdjNode visits each node connected to a given node.
//
func (abhr *AbstHier) EachAdjNode(gl GraphLevel, nl NodeLoc, Visit NodeLocFunc) {
	g := &abhr.Graphs[gl]
	if gl == PointLevel {
		c, r := GetColRow(nl)
		typ := g.Nodes[nl].GetPointType()
		if typ >= CenterPt { // center point visit all 4 neighbors
			// Visit point above
			Visit(MakeNodeLoc(c, r-1))
			// Visit point to left
			Visit(MakeNodeLoc(c-1, r))
			// Visit point below
			Visit(MakeNodeLoc(c, r+1))
			// Visit point to right
			Visit(MakeNodeLoc(c+1, r))
		} else { // less than four, test directions
			if (typ & PointType(UpperDir)) > 0 {
				Visit(MakeNodeLoc(c, r-1))
			}
			if (typ & PointType(LeftDir)) > 0 {
				Visit(MakeNodeLoc(c-1, r))
			}
			if (typ & PointType(LowerDir)) > 0 {
				Visit(MakeNodeLoc(c, r+1))
			}
			if (typ & PointType(RightDir)) > 0 {
				Visit(MakeNodeLoc(c+1, r))
			}
		}
	} else {
		a := g.Nodes[nl].inList
		for a != nilArc {
			fn := g.arcs[a].fromNode
			Visit(fn)
			a = g.arcs[a].inNext
		}
		a = g.Nodes[nl].outList
		for a != nilArc {
			tn := g.arcs[a].toNode
			Visit(tn)
			a = g.arcs[a].outNext
		}
	}
}

// BreadthFirstSearch performs Breadth First Search on a Board starting from a given point.
// To target of the BreadthFirstSearch is encapsulated in a parametric function.
//
func (abhr *AbstHier) BreadthFirstSearch(gl GraphLevel, startNl NodeLoc, IsTarget NodeLocFuncBool) *SearchStack {
	defer un(trace("BreadthFirstSearch", gl, startNl, 0, NilNodeLoc, nilArc, 0xFFFF))
	srchStk := new(SearchStack)
	g := &abhr.Graphs[gl]
	srchStk.g = g
	startPt := srchStk.g.Nodes[startNl]
	lookingFor := startPt.memberOf
	if TraceAH {
		PrintNodeLoc(gl+1, lookingFor, "  Looking for path in: ")
		fmt.Println()
	}
	srchStk = srchStk.PushAndMark(startNl)
	for srchStk.cur != uint16(len(srchStk.nods)) {
		abhr.EachAdjNode(gl, srchStk.nods[srchStk.cur],
			func(adjNl NodeLoc) {
				if srchStk.found == false { // don't do anything, once found
					if IsTarget(adjNl) {
						if TraceAH {
							PrintNodeLoc(gl, adjNl, "  IsTarget: ")
							fmt.Println()
						}
						srchStk.found = true
					}
					if g.Nodes[adjNl].IsBFSMarked() == false {
						if g.Nodes[adjNl].memberOf == lookingFor {
							if TraceAH {
								fmt.Println("  Found: ", lookingFor)
							}
							srchStk = srchStk.PushAndMark(adjNl)
						}
					}
				}
			})
		srchStk.cur += 1
	}
	return srchStk
}
