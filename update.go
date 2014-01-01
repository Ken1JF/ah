/*
 *  File:		src/gitHub.com/Ken1JF/ahgo/ah/update.go
 *  Project:	abst-hier
 *
 *  Created by Ken Friedenbach on 2/10/10.
 *  Copyright 2010-2014 Ken Friedenbach. All rights reserved.
 *
 *	This file implements the updating of Abstraction Hierarchies.
 */

package ah

import (
	"os"
	"fmt"
	"sort"
	"strconv"
	"unsafe"
)

type NodeStatus uint16


var ( // TODO: make this a per AbstHier variable
	TraceAH bool = false
)

func PrintNodeLoc(gl GraphLevel, nl NodeLoc, str string) {
	fmt.Print(str)
	if gl == PointLevel {
		c, r := GetColRow(nl)
		fmt.Print("(", c)
		fmt.Print(",", r)
		fmt.Print(")")
	} else {
		fmt.Print(nl)
	}
}

// trace and un provide tracing capability.
//
func trace(msg string, gl GraphLevel, nl NodeLoc, gl2 GraphLevel, nl2 NodeLoc, ai ArcIdx, ns NodeStatus) string {
	if TraceAH {
		fmt.Print("Entering ", msg, " Level ", int(gl))
		if nl != NilNodeLoc {
			PrintNodeLoc(gl, nl, " Node ")
		}
		if nl2 != NilNodeLoc {
			PrintNodeLoc(gl2, nl2, " Node2 ")
		}
		if ai != nilArc {
			fmt.Print(" ArcIdx ", ai)
		}
		if ns != 0xFFFF {
			fmt.Print(" NodeState ", ns)
		}
		fmt.Println(".")
	}
	return msg
}

func un(msg string) {
	if TraceAH {
		fmt.Println("Leaving", msg)
	}
}

// SetAHTrace turns tracing on and off
//
func SetAHTrace(newSt bool) {
	TraceAH = newSt
}

// GetAHTrace turns tracing on and off
//
func GetAHTrace() bool {
	return TraceAH
}

// GraphLevel is the type for levels of the Go Abstraction Hierarchy.
//
type GraphLevel uint8

const (
	PointLevel	GraphLevel = iota
	StringLevel	
	BlockLevel
	GroupLevel
	RegionLevel
	AreaLevel
)

// AbstHier is the Abstraction Hierarchy structure
//
type AbstHier struct {
	// exported temporarily for test_ahgo.go
	Graphs		[6] Graph 
//	strs	Graph	// Strings - capture, ko, laddars
//	blks	Graph	// Blocks - connections, cuts
//	grps	Graph	// Groups - eyes, life and death
//	rgns	Graph	// Regions
//	aras	Graph	// Areas
	// other Board info:
	Board
	// The depth of updating being done:
	updtLev		GraphLevel
	// variables used to control the amount of trace information printed
	black_nodes, black_members, white_nodes, white_members, unocc_nodes, unocc_members int
	moves_printed int

}

// GetNodeCounts counts nodes and members
//
func (abhr *AbstHier) GetNodeCounts(gl GraphLevel) (b_nodes int, b_members int, w_nodes int, w_members int, u_nodes int, u_members int) {
	var color uint16
	var count int
	var memCount int
	CountNode := func (g *Graph, nl NodeLoc) {
		if (g.Nodes[nl].lowState == color) || ((PointStatus(color) == Unocc)&&(g.Nodes[nl].lowState != uint16(Black)) && (g.Nodes[nl].lowState != uint16(White))) {
			count+=1
			abhr.EachMember(gl, nl,
				func (memNl NodeLoc) {
					memCount+=1
				})
		}
	}
	color = uint16(Black)
	memCount = 0
	abhr.EachNode(gl, CountNode)
	b_nodes = count
	b_members = memCount

	count = 0
	memCount = 0
	color = uint16(White)
	abhr.EachNode(gl, CountNode)
	w_nodes = count
	w_members = memCount

	count = 0
	memCount = 0
	color = uint16(Unocc)
	abhr.EachNode(gl, CountNode)
	u_nodes = count
	u_members = memCount
//	fmt.Println("Black nodes = ", b_nodes, " Black members = ", b_members);
//	fmt.Println("White nodes = ", w_nodes, " White members = ", w_members);
//	fmt.Println("Unocc nodes = ", u_nodes, " Unocc members = ", u_members);
	
	return
}

// PrintGraph prints the graph
//
func (abhr *AbstHier) PrintGraph(gl GraphLevel, printUnocc bool) {
	var color uint16
	var count int
	var memCount int
	PrintNode := func (g *Graph, nl NodeLoc) {
		var thisMemCount int
		// TODO: this is not sufficient for higher levels...
		if (g.Nodes[nl].lowState == color) || ((PointStatus(color) == Unocc)&&(g.Nodes[nl].lowState != uint16(Black))&&(g.Nodes[nl].lowState != uint16(White) )) {
			count+=1
			fmt.Print(nl)
			fmt.Print(": ", g.Nodes[nl].lowState)
			fmt.Print(",", g.Nodes[nl].highState)
			abhr.EachMember(gl, nl,
				func(memN1 NodeLoc) {
					thisMemCount+=1
				})
			fmt.Print(" ", thisMemCount)
			fmt.Print("-mem: ")
			var lo1, hi1 uint16
			loG := &abhr.Graphs[gl-1]
			abhr.EachMember(gl, nl,
				func (memNl NodeLoc) {
					memCount+=1
					if g.gLevel == StringLevel {
						c, r := GetColRow(memNl)
						fmt.Print("(", c)
						fmt.Print(",", r)
						fmt.Print(")")
					} else {
						fmt.Print(memNl)
					}
					lo := loG.Nodes[memNl].lowState
					hi := loG.Nodes[memNl].highState
					if (lo != lo1) || (hi != hi1) {
						fmt.Print(":", lo)
						fmt.Print(":", hi)
						lo1 = lo
						hi1 = hi
					}
					fmt.Print(", ")
				})
			fmt.Print(" adj: ")
			g.EachIncidentArc(nl, 
				func (a ArcIdx) {
					fNod := g.arcs[a].fromNode
					if fNod == nl {
						fmt.Print(g.arcs[a].toNode)
					} else {
						fmt.Print(fNod)
					}
					fmt.Print("(", g.arcs[a].imageCount)
					fmt.Print("), ")
				})
			fmt.Println()
		}
	}
	fmt.Println("Black nodes")
	color = uint16(Black)
	memCount = 0
	abhr.EachNode(gl, PrintNode)
	fmt.Println(" Total ", count, " nodes, with ", memCount, " members")

	count = 0
	memCount = 0
	fmt.Println("White nodes")
	color = uint16(White)
	abhr.EachNode(gl, PrintNode)
	fmt.Println(" Total ", count, " nodes, with ", memCount, " members")

	if printUnocc {
		count = 0
		memCount = 0
		fmt.Println("Unocc nodes")
		color = uint16(Unocc)
		abhr.EachNode(gl, PrintNode)
		fmt.Println(" Total ", count, " nodes, with ", memCount, " members")
	}
}

// DecrementMovesPrinted is used when backtracking, i.e. by UndoBoardMove
//
func (abhr *AbstHier) DecrementMovesPrinted() {
	if (abhr.moves_printed > 0) {
		abhr.moves_printed -= 1
	}
}

// SetBackMovesPrinted is used when moves have changed, i.e. when doing and undoing captures.
//
func (abhr *AbstHier) SetBackMovesPrinted(changed int16) {
	if (abhr.moves_printed > (int(changed) -1) ) {
		abhr.moves_printed = int(changed) - 1
	}
}

// PrintAbstHier prints the StringLevel and higher graphs,
// upto and including the updtLev.
//
func (abhr *AbstHier) PrintAbstHier(str string, printUnocc bool) {
	new_black_nodes, new_black_members, new_white_nodes, new_white_members, new_unocc_nodes, new_unocc_members := abhr.GetNodeCounts(StringLevel)
	if ( (new_black_nodes != abhr.black_nodes) || (new_black_members != abhr.black_members) ||
		 (new_white_nodes != abhr.white_nodes) || (new_white_members != abhr.white_members) ||
		 (new_unocc_nodes != abhr.unocc_nodes) || (new_unocc_members != abhr.unocc_members) ) {
		var lev GraphLevel
		fmt.Println("Abstraction Hierarchy: ", str)
		for lev = StringLevel; lev <= abhr.updtLev; lev++ {
			fmt.Println("  Level", lev)
			abhr.PrintGraph(lev, printUnocc)
		}
		for i, m := range abhr.movs {
			if (i >= abhr.moves_printed) {
				abhr.moves_printed = i
				m.PrintMove(i)
			}
		}
		abhr.black_nodes, abhr.black_members, abhr.white_nodes, abhr.white_members, abhr.unocc_nodes, abhr.unocc_members = new_black_nodes, new_black_members, new_white_nodes, new_white_members, new_unocc_nodes, new_unocc_members
	}
}

// 
//
func (abhr *AbstHier) initHoshiPts() {
	// Set the hoshi points:
	if (abhr.colSize >= 13) && (abhr.rowSize >= 13) {
		// Set hoshi points on the fourth line.
		//
		// Set the corner hoshi points:
		abhr.addHoshiPt(3, 3)
		abhr.addHoshiPt(abhr.colSize-4, 3)
		abhr.addHoshiPt(abhr.colSize-4, abhr.rowSize-4)
		abhr.addHoshiPt(3, abhr.rowSize-4)
		// Set the center hoshi point, if there is one:
		if ((abhr.colSize & 1) == 1) && ((abhr.rowSize & 1) == 1) {
			abhr.addHoshiPt(abhr.colSize/2, abhr.rowSize/2)
		}
		// Set the Side hoshi points, if there are any:
		if ((abhr.colSize >= 15) && ((abhr.colSize & 1) == 1)) {
			// Set the upper and lower side hoshi points:
			abhr.addHoshiPt(abhr.colSize/2, 3)
			abhr.addHoshiPt(abhr.colSize/2, abhr.rowSize-4)
		}
		if ((abhr.rowSize >= 15) && ((abhr.rowSize & 1) == 1)) {
			// Set the left and right side hoshi points:
			abhr.addHoshiPt(3, abhr.rowSize/2)
			abhr.addHoshiPt(abhr.colSize-4, abhr.rowSize/2)
		}
	}
	if (abhr.colSize <= 11) && (abhr.rowSize <= 11) {
		// Set hoshi points on the third line.
		if (abhr.colSize >= 7) && (abhr.rowSize >= 7) {
			// Set the corner points:
			abhr.addHoshiPt(2, 2)
			abhr.addHoshiPt(abhr.colSize-3, 2)
			abhr.addHoshiPt(abhr.colSize-3, abhr.rowSize-3)
			abhr.addHoshiPt(2, abhr.rowSize-3)
		}
		// Set the center hoshi point if there is one:
		if (abhr.colSize >= 9) && (abhr.rowSize >= 9) {
			if ((abhr.colSize & 1) == 1) && ((abhr.rowSize & 1) == 1) {
				abhr.addHoshiPt(abhr.colSize/2, abhr.rowSize/2)
			}
		} else if (abhr.colSize == 5) && (abhr.rowSize == 5) {
			abhr.addHoshiPt(abhr.colSize/2, abhr.rowSize/2)
		}
	}
	if len(abhr.hoshiPts) > 0 {
		sort.Sort(abhr.hoshiPts)
	}
}

var PrintAbstHierInit bool = false	// TODO: make these parser instance variables?
var saveAHTrace bool	// TODO: make these parser instance variables?

// setupAbstHier either clears an existing AH, or allocates and initializes a new AH.
// It is called by InitAbstHier
//
func (abhr *AbstHier) setupAbstHier(colSize ColValue, rowSize RowValue, upLev GraphLevel, doPlay bool,
		brdCHi CompStateFunc, brdCNw CompStateFunc,
		strCHi CompStateFunc, strCNw CompStateFunc,
		blkCHi CompStateFunc, blkCNw CompStateFunc,
		grpCHi CompStateFunc, grpCNw CompStateFunc,
		rgnCHi CompStateFunc, rgnCNw CompStateFunc,
		araCHi CompStateFunc, araCNw CompStateFunc) (*AbstHier) {
	defer un(trace("setupAbstHier", 0, NilNodeLoc, 0, NilNodeLoc, nilArc, 0xFFFF))
	var g,hiG *Graph
	if abhr != nil { // ClearAbstHier
		if TraceAH {
			fmt.Println("Clearing AbstHier", abhr);
		}
		abhr.Graphs[PointLevel].clearGraph(PointLevel, brdCHi, brdCNw, NodeStatus(UndefinedPointStatus))
		abhr.Graphs[StringLevel].clearGraph(StringLevel, strCHi, strCNw, NodeStatus(UndefinedStringStatus))
		abhr.Graphs[BlockLevel].clearGraph(BlockLevel, blkCHi, blkCNw, NodeStatus(0))
		abhr.Graphs[GroupLevel].clearGraph(GroupLevel, grpCHi, grpCNw, NodeStatus(0))
		abhr.Graphs[RegionLevel].clearGraph(RegionLevel, rgnCHi, rgnCNw, NodeStatus(0))
		abhr.Graphs[AreaLevel].clearGraph(AreaLevel, araCHi, araCNw, NodeStatus(0))
		abhr.setSize(colSize, rowSize)
		abhr.updtLev = upLev
	} else {
		if TraceAH {
			fmt.Println("New AbstHier");
		}
		abhr = new(AbstHier)
		abhr.setSize(colSize, rowSize)
		abhr.updtLev = upLev
		abhr.Graphs[PointLevel].initGraph(PointLevel, brdCHi, brdCNw, NodeStatus(UndefinedPointStatus))
		abhr.Graphs[StringLevel].initGraph(StringLevel, strCHi, strCNw, NodeStatus(UndefinedStringStatus))
		abhr.Graphs[BlockLevel].initGraph(BlockLevel, blkCHi, blkCNw, NodeStatus(0))
		abhr.Graphs[GroupLevel].initGraph(GroupLevel, grpCHi, grpCNw, NodeStatus(0))
		abhr.Graphs[RegionLevel].initGraph(RegionLevel, rgnCHi, rgnCNw, NodeStatus(0))
		abhr.Graphs[AreaLevel].initGraph(AreaLevel, araCHi, araCNw, NodeStatus(0))
	}
	if doPlay {
		switch upLev {
			case AreaLevel:
				g = &abhr.Graphs[AreaLevel]
				g.initNode = g.AddGraphNode(0)
				// TODO: set highState, as below
				fallthrough
			case RegionLevel:
				hiG = g
				g = &abhr.Graphs[RegionLevel]
				g.initNode = g.AddGraphNode(0)
				// TODO: set highState, as below
				if RegionLevel < upLev {
					abhr.AddMember(RegionLevel, g.initNode, hiG.initNode)
				}
				fallthrough
			case GroupLevel:
				hiG = g
				g = &abhr.Graphs[GroupLevel]
				g.initNode = g.AddGraphNode(0)
				// TODO: set highState, as below
				if GroupLevel < upLev {
					abhr.AddMember(GroupLevel, g.initNode, abhr.Graphs[RegionLevel].initNode)
				}
				fallthrough
			case BlockLevel:
				hiG = g
				g = &abhr.Graphs[BlockLevel]
				g.initNode = g.AddGraphNode(0)
				// TODO: set highState, as below
				if BlockLevel < upLev {
					abhr.AddMember(BlockLevel, g.initNode, abhr.Graphs[GroupLevel].initNode)
				}
				fallthrough
			case StringLevel:
				hiG = g
				g = &abhr.Graphs[StringLevel]
				g.initNode = g.AddGraphNode(uint16(UnoccupiedString))
				hiSt := g.CompHigh(abhr, StringLevel, g.initNode , uint16(UnoccupiedString))
				g.Nodes[g.initNode].highState = uint16(hiSt)
				if StringLevel < upLev {
					abhr.AddMember(StringLevel, g.initNode, abhr.Graphs[BlockLevel].initNode)
				}
				fallthrough
			case PointLevel:
				g = &abhr.Graphs[PointLevel]
				hiG = &abhr.Graphs[StringLevel]	// make sure hiG is valid
				// Add all the points to initNode
	//			saveAHTrace = TraceAH // don't trace this
	//			TraceAH = false
				var c ColValue
				var r RowValue
				for r = 0; r < rowSize; r++ {
					for c = 0; c < colSize; c++ {
						nl := MakeNodeLoc(c, r)
						abhr.AddMember(PointLevel, nl, hiG.initNode)
					}
				}
	//			TraceAH = saveAHTrace	// restore
	//			if PrintAbstHierInit == true {
	//				abhr.PrintAbstHier("Before Changing nodes to initial states", true)
	//				SetAHTrace(true)
	//			} else {	// If Printing is not on, turn off tracing during init
	//				saveAHTrace = TraceAH
	//				TraceAH = false
	//			}
				// Change the nodes to their initial states
				for r = 0; r < rowSize; r++ {
					for c = 0; c < colSize; c++ {
						nl := MakeNodeLoc(c, r)
						hiSt := brdCHi(abhr, PointLevel, nl, uint16(Unocc))
						abhr.ChangeNodeState(PointLevel, nl, NodeStatus(hiSt), true)
						if TraceAH == true {
							abhr.PrintAbstHier("After Changing row "+ strconv.Itoa(int(r)) + " col " + strconv.Itoa(int(c)), false )
						}
					}
				}
	//			if PrintAbstHierInit == true {
	//				PrintAbstHierInit = false
	//				SetAHTrace(true)	// TODO: need this?
	//			} else {
	//				TraceAH = saveAHTrace	// restore to global setting
	//			}
		}
	}

	return abhr
}

// FindMoveNumber searches the Moves to find the most recent move at a point
//
func (abhr *AbstHier) FindMoveNumber(nl NodeLoc) int16 {
	var movNum int16 = nilMovNum
	var i int16
	for i = int16(len(abhr.movs))-1; (i >= 0) && (movNum == nilMovNum); i-- {
		if abhr.movs[i].moveLoc == nl {
			movNum = i
		}
	}
	return movNum
}

// NumElements returns the number of elements in a partition
//
func (abhr *AbstHier) NumElements(gl GraphLevel, aMem NodeLoc) (ret int) {
	g := &abhr.Graphs[gl]
	hiG := &abhr.Graphs[gl+1]
	hiNod := g.Nodes[aMem].memberOf
	mem := hiG.Nodes[hiNod].firstMem
	// TODO: add check for infinite loop???
	for mem != NilNodeLoc {
		ret += 1
		mem = g.Nodes[mem].nextSame
	}
	return ret
}

// AddMember adda a member node to a higher node member list
//
func (abhr *AbstHier) AddMember(gl GraphLevel, mem NodeLoc, hiNod NodeLoc) {
	defer un(trace("AddMember", gl, mem, gl+1, hiNod, nilArc, 0xFFFF))
	g := &abhr.Graphs[gl]
	hiG := &abhr.Graphs[gl+1]
	g.Nodes[mem].memberOf = hiNod
	g.Nodes[mem].nextSame = hiG.Nodes[hiNod].firstMem
	hiG.Nodes[hiNod].firstMem = mem
}

// DeleteMember deletes a member node from a higer node member list
//
func (abhr *AbstHier) DeleteMember(gl GraphLevel, mem NodeLoc, hiNod NodeLoc) {
	defer un(trace("DeleteMember", gl, mem, gl+1, hiNod, nilArc, 0xFFFF))
	g := &abhr.Graphs[gl]
	hiG := &abhr.Graphs[gl+1]
	m := hiG.Nodes[hiNod].firstMem
	if m == mem {
		hiG.Nodes[hiNod].firstMem = g.Nodes[mem].nextSame
		if hiG.Nodes[hiNod].firstMem == NilNodeLoc {
			hiG.DeleteNode(hiNod)
		}
	} else {
Loop:
		for {
			prevM := m
			m =  g.Nodes[m].nextSame
			if m == mem {
				g.Nodes[prevM].nextSame = g.Nodes[m].nextSame
				break Loop
			}
			if m == NilNodeLoc {
				FatalAbstHierError("DeleteMember: member not found")
			}
		}
	}
}

// EachMember visits each member of node in an abstraction hierarchy
//
func (abhr *AbstHier) EachMember(gl GraphLevel, nl NodeLoc, visit NodeLocFunc) {
	// TODO: check that gl >= StringLevel?
	g := &abhr.Graphs[gl]
	loG := &abhr.Graphs[gl-1]
	mem := g.Nodes[nl].firstMem
	// TODO: add check for infinite loop???
	for mem != NilNodeLoc {
		nextMem := loG.Nodes[mem].nextSame
		visit(mem)
		mem = nextMem
	}
}

// AddNodeLow adds a graph node into an abstraction hierarchy
// It returns the hiNod, and a bool to indicate if the hiNod is new
//
func (abhr *AbstHier) AddNodeLow(gl GraphLevel, lowSt uint16) (hiNod NodeLoc, isNew bool) {
	defer un(trace("AddNodeLow", 0, NilNodeLoc, 0, NilNodeLoc, nilArc, NodeStatus(lowSt)))
	g := &abhr.Graphs[gl]
	hiNod = g.AddGraphNode(lowSt)
	hiSt := g.compNew(abhr, gl, hiNod, lowSt)
	abhr.Graphs[gl].Nodes[hiNod].highState = uint16(hiSt)
	isNew = abhr.AddNodeHigh(gl, hiNod, hiSt)
	return hiNod, isNew
}

// AddNodeHigh adds a node to a graph
//
func (abhr *AbstHier) AddNodeHigh(gl GraphLevel, chgNl NodeLoc, newSt uint16) ( bool) {
	defer un(trace("AddNodeHigh", gl, chgNl, 0, NilNodeLoc, nilArc,  NodeStatus(newSt)))
	var sameNod NodeLoc = NilNodeLoc
	var hiNod NodeLoc
	var isNew bool
	abhr.EachAdjNode(gl, chgNl,
		func (adjNl NodeLoc) { // FindSamePtType
			adjPt := &abhr.Graphs[gl].Nodes[adjNl]
			if adjPt.GetNodeHighState() == newSt { 
				sameNod = adjNl 
			}
		})
	if sameNod == NilNodeLoc {
		if abhr.updtLev <= (gl+1) {
			hiG := &abhr.Graphs[gl+1]
			hiNod = hiG.AddGraphNode(uint16(newSt))
			hiSt := hiG.compNew(abhr, gl+1, hiNod, newSt)
			hiG.Nodes[hiNod].highState = uint16(hiSt)
		} else {
			hiNod, isNew = abhr.AddNodeLow(gl+1, uint16(newSt))
		}
	} else {
		hiNod = abhr.Graphs[gl].Nodes[sameNod].memberOf
	}
	abhr.AddMember(gl, chgNl, hiNod)
	return isNew
}

// CheckMerge is called after an arc is added to see if it causes a merge
//
func (abhr *AbstHier) CheckMerge(gl GraphLevel, n1 NodeLoc, n2 NodeLoc) {
	defer un(trace("CheckMerge", gl, n1, gl, n2, nilArc, 0xFFFF))
	g := &abhr.Graphs[gl]
	high1 := g.Nodes[n1].memberOf
	high2 := g.Nodes[n2].memberOf
	MergeNodes := func(nod NodeLoc, hiNod NodeLoc) {
		defer un(trace(" MergeNodes", gl, nod, gl+1, hiNod, nilArc, 0xFFFF))
		var mrgStk *SearchStack
		saveSt := g.Nodes[nod].highState
		mrgStk = abhr.BreadthFirstSearch(gl, nod, 
			func(inNL NodeLoc) (ret bool) { // Find target, ==> false
				return false
			})
		if mrgStk != nil {
			mrgStk.UnMarkNodes()
		}
		// set the nodes to undefStatus
		for _,p := range mrgStk.nods {
			if TraceAH {
				PrintNodeLoc(gl, p, " Setting to undefStatus")
				fmt.Println()
			}
			g.Nodes[p].highState = uint16(g.undefStatus)
		}
		// change them back to saveSt, without checking for Merge/Split
		for _,p := range mrgStk.nods {
			if TraceAH {
				PrintNodeLoc(gl, p, " Calling to change to saveState")
				fmt.Print(" ")
			}
			abhr.ChangeNodeState(gl, p, NodeStatus(saveSt), false)
		}
	}
	if high1 != high2 { // n1 and n2 are members of different nodes
		hiSt1 := g.Nodes[n1].highState
		hiSt2 := g.Nodes[n2].highState
		if hiSt1 == hiSt2 { // but high states are the same, so MergeNodes
			// BiggerPart: true if member list for n1 longer than n2
			big := false
			// TODO: why check for NilNodeLoc, aren't these error conditions?
			if high1 == NilNodeLoc {
				// done
			} else if high2 == NilNodeLoc {
				big = true
			} else { // need to look
				hiG := &abhr.Graphs[gl+1]
				mem1 := hiG.Nodes[high1].firstMem
				mem2 := hiG.Nodes[high2].firstMem
				for (mem1 != NilNodeLoc) && (mem2 != NilNodeLoc) {
					mem1 = g.Nodes[mem1].nextSame
					mem2 = g.Nodes[mem2].nextSame
				}
				if mem1 != NilNodeLoc {
					big = true
				}
			}
			if big {
				MergeNodes(n2, high2)
			} else {
				MergeNodes(n1, high1)
			}
		}
	}
}

// CheckSplit checks if removal of arc from nod1 to nod2 causes a split.
// If so, perform the split.
// Returns the original component, in case of multiple splits.
//
func (abhr *AbstHier) CheckSplit (gl GraphLevel, n1 NodeLoc, n2 NodeLoc) (ret NodeLoc) {
	defer un(trace("CheckSplit", gl, n1, gl, n2, nilArc, 0xFFFF))
	if TraceAH {
		abhr.PrintAbstHier("Entering CheckSplit", true)
	}
	var spltStk *SearchStack
	g := &abhr.Graphs[gl]
	high1 := g.Nodes[n1].memberOf
	targN := n1
	pathInPart := func(low1 NodeLoc, low2 NodeLoc) (found bool) {
		defer un(trace("pathInPart", gl, low1, gl, low2, nilArc, 0xFFFF))
		spltStk = abhr.BreadthFirstSearch(gl, low2,
				func(inNL NodeLoc) (bool) { // inPart
					return (inNL == targN)
				})
		if TraceAH {
			fmt.Println(" After BreadthFirstSearch, found = ", spltStk.found)
		}
		return spltStk.found
	}
	splitPoints := func(nod NodeLoc) {
		defer un(trace("splitPoints", gl, nod, 0, NilNodeLoc, nilArc, 0xFFFF))
		// save the current state
		saveSt := g.Nodes[nod].GetNodeLowState()
		// use ChangeNodeState (without Merge/Split checking) to change to UndefinedPointStatus
		for _, n := range spltStk.nods {
			if TraceAH {
				PrintNodeLoc(gl, n, " Calling to set to Undefined: ")
				fmt.Println()
			}
			abhr.ChangeNodeState(gl, n, g.undefStatus, false)
		}
		if TraceAH {
			fmt.Println(" finished set Undefined, set to saveState, len spltStk.nods = ", len(spltStk.nods))
		}
		// set the split Nodes back to saveState
		// for _,nn := range spltStk.nods {
		for _, nn := range spltStk.nods {
			if TraceAH {
				PrintNodeLoc(gl, nn, " computing the high value for node: ")
				fmt.Println("    nn = ", nn, " saveSt = ", saveSt);
			}
			if TraceAH {
				PrintNodeLoc(gl, nn, " Setting back to saveSt: ")
				fmt.Println(" lowState: ", saveSt)
			}
			g.Nodes[nn].lowState = saveSt
			g.Nodes[nn].highState = saveSt
			if abhr.updtLev <= (gl+1) {
				// TODO: verify this?
				abhr.Graphs[gl+1].Nodes[g.Nodes[nn].memberOf].lowState = saveSt
			} else {
				// TODO: need to ChangeNodeState for higher levels
				abhr.ChangeNodeState(gl+1, g.Nodes[nn].memberOf, NodeStatus(saveSt), false)	// TODO: or true?
			}
		}
		if TraceAH {
			fmt.Println(" finished set set to saveState, len spltStk.nods = ", len(spltStk.nods))
		}
		// For the above to work, make sure coloring functions satisfy the following:
		//	1. Each level has an "Undefined" state (say 0)
		//	2. Each coloring function maps "Undefined" to "Undefined".
		//	3. Coloring function is prepared to see "Undefined" in the Neighborhood.
		//	4. Rest of program does not expect (or use!) the "Undefined" state.
		
		// save the high node:
		hiNod := g.Nodes[spltStk.nods[0]].memberOf
		loSt := saveSt
		for lev := gl; lev < abhr.updtLev; lev++ {
			// TODO: is this highState or lowState being set? check and call compHi later?
			gLev := &abhr.Graphs[lev+1]
			if TraceAH {
				PrintNodeLoc(lev+1, hiNod, " Setting high node: ")
				fmt.Println()
			}
			gLev.Nodes[hiNod].highState = gLev.compNew(abhr, lev+1, hiNod, loSt)
			loSt = gLev.Nodes[hiNod].highState
			hiNod = gLev.Nodes[hiNod].memberOf
		}
	}
	ret = n1
	if high1 == g.Nodes[n2].memberOf {
		if !pathInPart(n1, n2) {
			spltStk.UnMarkNodes()
			pSize := 0;	// PartSize
			nod := abhr.Graphs[gl+1].Nodes[high1].firstMem
			for nod != NilNodeLoc {
				pSize += 1
				nod = g.Nodes[nod].nextSame
				if pSize > MaxListLen {
					FatalAbstHierError("CheckSplit: PartSize too big")
				}
			}
			if pSize >= 2*len(spltStk.nods) {
				splitPoints(n2)
			} else {
				targN = n2
				pathInPart(n2, n1)
				spltStk.UnMarkNodes()
				splitPoints(n1)
				ret = n2
			}
		} else {
			spltStk.UnMarkNodes()
		}
	}
	if TraceAH {
		abhr.PrintAbstHier("Leaving CheckSplit", true)
	}
	return ret
}

// DeleteArcHigh deletes the image of an arc from the high graph
//
func (abhr *AbstHier) DeleteArcHigh(gl GraphLevel, chgNod NodeLoc, adjNod NodeLoc) {
	defer un(trace("DeleteArcHigh", gl, chgNod, gl, adjNod, nilArc, 0xFFFF))
	g := &abhr.Graphs[gl]
	highChg := g.Nodes[chgNod].memberOf
	highAdj := g.Nodes[adjNod].memberOf
	if highChg != highAdj { // not in same high node
		hiG := &abhr.Graphs[gl+1]
		highArc := hiG.FindEdge(highChg, highAdj)
		if highArc == nilArc {
			FatalAbstHierError("DeleteArcHigh: no arc image")
		}
		hiG.arcs[highArc].imageCount -= 1
		if hiG.arcs[highArc].imageCount <= 0 {
			if abhr.updtLev <= hiG.gLevel  {
				hiG.DeleteEdge(highChg, highAdj)
			} else { // deleteArcLow
				// deleteArcLow(gl+1, highChg, highAdj)
				abhr.DeleteArcHigh(gl+1, highChg, highAdj)
				hiG.DeleteEdge(highChg, highAdj)
				abhr.CheckSplit(gl, chgNod, adjNod)
			}
		}
	}
}

// ChangeNodeState is called to change the NodeStatus of a GraphNode
//
func (abhr *AbstHier) ChangeNodeState(gl GraphLevel, chgNod NodeLoc, newState NodeStatus, doMrgSplt bool) {
	defer un(trace("ChangeNodeState", gl, chgNod, 0, NilNodeLoc, nilArc, newState))
	if TraceAH {
		fmt.Println(" Move: ", abhr.numMoves)
	}
	g := &abhr.Graphs[gl]
	// save the original State
	origState := g.Nodes[chgNod].GetNodeLowState()
	// delete the images of the edges
	abhr.EachAdjNode(gl, chgNod, 
		func (adjNod NodeLoc) {
			abhr.DeleteArcHigh(gl, chgNod, adjNod)
		}) 
	// delete the point from the string level graph
	// TODO: xxxx should call DeleteNodeHigh, not DeleteMember (for higher update levels)
	abhr.DeleteMember(gl, chgNod, g.Nodes[chgNod].memberOf)
	// change the State
	g.Nodes[chgNod].lowState = uint16(newState)
	g.Nodes[chgNod].highState = uint16(newState)
	// add the point back to the higher levels
	_ = abhr.AddNodeHigh(gl, chgNod, uint16(newState)) // needCheck
	// add the edges
	abhr.EachAdjNode(gl, chgNod,
		func(adjNod NodeLoc) {
			abhr.AddArcHigh(gl, chgNod, adjNod)
		})
	if doMrgSplt {
		// do the check for merges and splits
		// TODO: why are they pushed on a stack? 
		// (part of shard code?) (to reverse the order?) (or ???)
		var changeStack *SearchStack = new(SearchStack)
		changeStack.g = g
		var cs_split NodeLoc = NilNodeLoc
		abhr.EachAdjNode(gl, chgNod,
			func (adjNod NodeLoc) { // pushPoint
				changeStack = changeStack.PushAndMark(adjNod)
			})
		changeStack.UnMarkNodes()
		stkSz := len(changeStack.nods)
		for	stkSz > 0 { // chk_cs
			elem := changeStack.nods[stkSz-1]
			loSt := g.Nodes[elem].GetNodeLowState()
			if loSt == origState {
				if cs_split == NilNodeLoc {
					if TraceAH {
						PrintNodeLoc(changeStack.g.gLevel, elem, " chk_cs, setting cs_split: ")
						fmt.Println()
					}
					cs_split = elem
				} else {
					if TraceAH {
						PrintNodeLoc(changeStack.g.gLevel, cs_split, " chk_cs, calling CheckSplit: ")
						PrintNodeLoc(changeStack.g.gLevel, elem, ", ")
						fmt.Println()
					}
					cs_split = abhr.CheckSplit(gl, cs_split, elem)
				}
			} else if loSt == uint16(newState) { 
				if TraceAH {
					PrintNodeLoc(changeStack.g.gLevel, elem, " chk_cs, testing for CheckMerge: ")
					PrintNodeLoc(changeStack.g.gLevel, chgNod, ", ")
					fmt.Println()
				}
				if g.Nodes[elem].memberOf != g.Nodes[chgNod].memberOf {
					abhr.CheckMerge (gl, elem, chgNod) 
				}
			}
			stkSz -= 1
		}
	}
	if TraceAH {
		abhr.PrintAbstHier(" When leaving ChangeNodeState", true)
	}
}

// AddArcHigh is called to add an arc connecting two adjacent Nodes.
//
func (abhr *AbstHier) AddArcHigh(gl GraphLevel, chgNod NodeLoc, adjNod NodeLoc) {
	defer un(trace("AddArcHigh", gl, chgNod, gl, adjNod, nilArc, 0xFFFF))
	var highA ArcIdx = nilArc
	g := &abhr.Graphs[gl]
	highChg := g.Nodes[chgNod].memberOf
	highAdj := g.Nodes[adjNod].memberOf
	if highChg != highAdj {
		hiG := &abhr.Graphs[gl+1]
		highA = hiG.FindEdge(highChg, highAdj)
		if highA == nilArc {
			if abhr.updtLev <= hiG.gLevel {
				hiG.arcs, _ = hiG.AddEdge(highChg, highAdj)
			} else {
				abhr.AddArcLow(gl+1, highChg, highAdj)
			}
		} else {
			if TraceAH {
				fmt.Println(" Incrementing: ", highA)
			}
			hiG.arcs[highA].imageCount += 1
		}
	}
}

// AddArcLow is called to add an arc into an abstraction hierarchy
//
func (abhr *AbstHier) AddArcLow(gl GraphLevel, chgNod NodeLoc, adjNod NodeLoc) ( ArcIdx) {
	defer un(trace("AddArcLow", gl, chgNod, gl, adjNod, nilArc, 0xFFFF))
	newA := nilArc
	g := &abhr.Graphs[gl]
	g.arcs, newA = g.AddEdge(chgNod, adjNod)
	abhr.AddArcHigh(gl+1, g.Nodes[chgNod].memberOf, g.Nodes[adjNod].memberOf)
	abhr.CheckMerge(gl, chgNod, adjNod)
	return newA
}

// FatalAbstHierError is used to report "should not happen" errors,
// mostly related to nil pointer checks, for elements that should be
// present, but are not found.
// TODO: add more debug info?
// TODO: design a way to restart???
//
func FatalAbstHierError(str string) {
	fmt.Println("Fatal Error:", str)
	os.Exit(999)
}

// Some functions to print the size and alignment of types:
// TODO: delete when no longer needed (development aids)
//
//func printSizeAlign(s string, sz int, al int) {
func printSizeAlign(s string, sz uintptr, al uintptr) {
	fmt.Println("Type", s, "size", sz, "alignment", al)
}

func PrintAhStructSizes() {
// board.go
	var d Direction
	var pt PointType
	var cv ColValue
	var rv RowValue
	var nl NodeLoc
	var nll NodeLocList
	var m MoveRecord
//var bp BoardPoint
	var b Board
	var ss SearchStack
	var gnlf GraphNodeLocFunc
	var nlfb NodeLocFuncBool
	printSizeAlign("Direction", unsafe.Sizeof(d), unsafe.Alignof(d))
	printSizeAlign("PointType", unsafe.Sizeof(pt), unsafe.Alignof(pt))
	printSizeAlign("ColValue", unsafe.Sizeof(cv), unsafe.Alignof(cv))
	printSizeAlign("RowValue", unsafe.Sizeof(rv), unsafe.Alignof(rv))
	printSizeAlign("NodeLoc", unsafe.Sizeof(nl), unsafe.Alignof(nl))
	printSizeAlign("NodeLocList", unsafe.Sizeof(nll), unsafe.Alignof(nll))
	printSizeAlign("MoveRecord", unsafe.Sizeof(m), unsafe.Alignof(m))
//	printSizeAlign("BoardPoint", unsafe.Sizeof(bp), unsafe.Alignof(bp))
	printSizeAlign("Board", unsafe.Sizeof(b), unsafe.Alignof(b))
	printSizeAlign("SearchStack", unsafe.Sizeof(ss), unsafe.Alignof(ss))
	printSizeAlign("GraphNodeLocFunc", unsafe.Sizeof(gnlf), unsafe.Alignof(gnlf))
	printSizeAlign("NodeLocFuncBool", unsafe.Sizeof(nlfb), unsafe.Alignof(nlfb))
// trans.go
	var bt BoardTrans
	printSizeAlign("BoardTrans", unsafe.Sizeof(bt), unsafe.Alignof(bt))
// graph.go
	var ai ArcIdx
	var gm GraphMark
	var gn GraphNode
	var ga GraphArc
	var csf CompStateFunc
	var chreq ChangeRequest
	var g Graph
	var glnf NodeLocFunc
	var af ArcFunc
	printSizeAlign("ArcIdx", unsafe.Sizeof(ai), unsafe.Alignof(ai))
	printSizeAlign("GraphMark", unsafe.Sizeof(gm), unsafe.Alignof(gm))
	printSizeAlign("GraphNode", unsafe.Sizeof(gn), unsafe.Alignof(gn))
	printSizeAlign("GraphArc", unsafe.Sizeof(ga), unsafe.Alignof(ga))
	printSizeAlign("CompStateFunc", unsafe.Sizeof(csf), unsafe.Alignof(csf))
	printSizeAlign("ChangeRequest", unsafe.Sizeof(chreq), unsafe.Alignof(chreq))
	printSizeAlign("Graph", unsafe.Sizeof(g), unsafe.Alignof(g))
	printSizeAlign("NodeLocFunc", unsafe.Sizeof(glnf), unsafe.Alignof(glnf))
	printSizeAlign("ArcFunc", unsafe.Sizeof(af), unsafe.Alignof(af))
// comp.go
	var ps PointStatus
	printSizeAlign("PointStatus", unsafe.Sizeof(ps), unsafe.Alignof(ps))
// update.go
	var gl GraphLevel
	var ah AbstHier
	var ns NodeStatus
	printSizeAlign("GraphLevel", unsafe.Sizeof(gl), unsafe.Alignof(gl))
	printSizeAlign("AbstHier", unsafe.Sizeof(ah), unsafe.Alignof(ah))
	printSizeAlign("NodeStatus", unsafe.Sizeof(ns), unsafe.Alignof(ns))
// gostrings.go
	var sst StringStatus
	printSizeAlign("StringStatus", unsafe.Sizeof(sst), unsafe.Alignof(sst))
}
