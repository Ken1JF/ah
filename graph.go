/*
 *  File:		src/gitHub.com/Ken1JF/ahgo/ah/graph.go
 *  Project:	abst-hier
 *
 *  Created by Ken Friedenbach on 2/10/10.
 *  Copyright 2010-2014 Ken Friedenbach. All rights reserved.
 *
 *	This file implements the data structures for Directed Graphs,
 *	used by the Abstraction Hierarchy algorithms.
 *
 */

package ah

import (
	"fmt"
)

type ArcIdx uint16

const nilArc ArcIdx = 0xFFFF

const MAX_ARC_LIST = 500

// NilNodeLoc is a list terminator value for NodeLoc lists
// 0 cannot be used, because it is C,R = 0,0
// TODO: only needs to be 10 bits.
// Can use 4 or 6 upper bits for other purposes...
const NilNodeLoc NodeLoc = 0xFFFF

type GraphMark uint8

const (
	OnFreeList GraphMark = (1 << iota)
	BFSMark
	CRMark
)

// GraphNode is the node type for Strings, Groups, Areas, etc.
type GraphNode struct {
	// base node fields, common to all Graphs: Points, Strings, and higher
	highState uint16
	lowState  uint16
	memberOf  NodeLoc
	nextSame  NodeLoc // next member with same memberOf parent
	mark      GraphMark
	unused    uint8
	// higher node fields, used by Strings, and higher graphs. Special uses for Points
	inList   ArcIdx  // used for PointType at Board level
	outList  ArcIdx  // unused at Board level
	firstMem NodeLoc // used for c, r at Board level
}

// Accessor functions: if needed, i.e. outside ah package
func (gn *GraphNode) SetNodeLowState(ns uint16) {
	gn.lowState = ns
}

// GetNodeLowState returns the lowState
func (gn *GraphNode) GetNodeLowState() uint16 {
	return gn.lowState
}

// IsBFSMarked returns true if the BFSMark bit is set
func (gn *GraphNode) IsBFSMarked() bool {
	return (gn.mark & BFSMark) > 0
}

// MarkBFSNode sets the BFSMark bit
func (gn *GraphNode) MarkBFSNode() {
	gn.mark |= BFSMark
}

// ClearBFSMark clears the BFSMark bit
func (gn *GraphNode) ClearBFSMark() {
	gn.mark &^= BFSMark
}

// IsCRMarked returns true if the CRMark bit is set
func (gn *GraphNode) IsCRMarked() bool {
	return (gn.mark & CRMark) > 0
}

// MarkCRNode sets the CRMark bit
func (gn *GraphNode) MarkCRNode() {
	gn.mark |= CRMark
}

// ClearCRMark clears the CRMark bit
func (gn *GraphNode) ClearCRMark() {
	gn.mark &^= CRMark
}

// SetNodeHighState set the highState field
func (gn *GraphNode) SetNodeHighState(ns uint16) {
	gn.highState = ns
}

// GetNodeHighState return the highState field
func (gn *GraphNode) GetNodeHighState() uint16 {
	return gn.highState
}

// GraphArc is the arc type for connecting two Nodes.
type GraphArc struct {
	fromNode   NodeLoc
	toNode     NodeLoc
	inNext     ArcIdx
	outNext    ArcIdx
	imageCount uint16
}

type CompStateFunc func(*AbstHier, GraphLevel, NodeLoc, uint16) uint16

// ChangeRequest is the type for registering objects to check for changes
type ChangeRequest struct {
	chgNL  NodeLoc
	chgLev GraphLevel
}

// RequestChange records the request to check for a change:
func (g *Graph) RequestChange(nl NodeLoc, adjNL NodeLoc) {
	var memNL NodeLoc
	switch g.gLevel {
	// TODO: make this work for higher levels
	case AreaLevel:
		fallthrough
	case RegionLevel:
		fallthrough
	case GroupLevel:
		fallthrough
	case BlockLevel:
		fallthrough
	case StringLevel:
		FatalAbstHierError("RequestChange - not implemented above PointLevel")
	case PointLevel:
		memNL = nl
	}
	if g.Nodes[memNL].IsCRMarked() == false {
		var newReq ChangeRequest
		newReq.chgLev = g.gLevel
		newReq.chgNL = memNL
		g.Nodes[memNL].MarkCRNode()
		g.changeList = append(g.changeList, newReq)
	}
}

// Graph is the type for levels of a Go Abstraction Hierarchy
type Graph struct {
	// storage for Nodes and arcs
	// exported temporarily for test_ahgo.go
	Nodes []GraphNode
	arcs  []GraphArc
	// storage for change requests
	changeList []ChangeRequest
	// functions to compute State value
	CompHigh CompStateFunc
	compNew  CompStateFunc
	// lists of available Nodes and arcs
	freeNodes NodeLoc
	freeArcs  ArcIdx
	// the level of this graph in the hierarchy
	gLevel      GraphLevel
	undefStatus NodeStatus
	// an initial node, for use in initialization:
	initNode NodeLoc
	// TODO: what about MoveTreeNodes, ListItems, Analysis Nodes, and Deferred requests ?
	// (see GraphOps.P in GoMacAppConv project for Pascal example code.)
}

// initGraph must be called before a graph can be used
func (g *Graph) initGraph(gl GraphLevel, cHi CompStateFunc, cNw CompStateFunc, udef NodeStatus) {
	g.CompHigh = cHi
	g.compNew = cNw
	g.gLevel = gl
	g.freeNodes = NilNodeLoc
	g.freeArcs = nilArc
	g.undefStatus = udef
}

// clearGraph is called when a graph is being reset, e.g. board change size
func (g *Graph) clearGraph(gl GraphLevel, cHi CompStateFunc, cNw CompStateFunc, udef NodeStatus) {
	g.gLevel = gl
	g.CompHigh = cHi
	g.compNew = cNw
	g.freeNodes = NilNodeLoc
	g.freeArcs = nilArc
	g.undefStatus = udef
	// TODO: check if these are necessary?
	if g.gLevel > PointLevel {
		g.Nodes = g.Nodes[0:0]
		g.arcs = g.arcs[0:0]
	}
}

// initBoardPoints() must be called before using a Board.
// currently called by setSize().
// initBoardPoints sets the static point type in inList
func (brd *Graph) initBoardPoints(colSize ColValue, rowSize RowValue) {
	var c ColValue
	var r RowValue
	nl := MakeNodeLoc(colSize-1, rowSize-1)
	brd.Nodes = make([]GraphNode, int(nl+1))
	// Set the default point type and ptCol, ptRow:
	for r = 0; r < rowSize; r++ {
		for c = 0; c < colSize; c++ {
			nl = MakeNodeLoc(c, r)
			brd.Nodes[nl].inList = ArcIdx(CenterPt)
			brd.Nodes[nl].firstMem = nl
		}
	}
	// Check for special Board sizes:
	if rowSize == 1 {
		if colSize == 1 {
			// 1 by 1 Board
			brd.Nodes[MakeNodeLoc(0, 0)].inList = ArcIdx(SingletonPt)
		} else { // colSize >= 2
			// 1 by N Board
			brd.Nodes[MakeNodeLoc(0, 0)].inList = ArcIdx(LeftEndPt)
			for c = 1; c < (colSize - 1); c++ {
				brd.Nodes[MakeNodeLoc(c, 0)].inList = ArcIdx(LeftRightBridgePt)
			}
			brd.Nodes[MakeNodeLoc(colSize-1, 0)].inList = ArcIdx(RightEndPt)
		}
	} else if colSize == 1 { // rowSize >= 2
		// N by 1 Board
		brd.Nodes[MakeNodeLoc(0, 0)].inList = ArcIdx(UpperEndPt)
		for r = 1; r < (rowSize - 1); r++ {
			brd.Nodes[MakeNodeLoc(0, r)].inList = ArcIdx(UpperLowerBridgePt)
		}
		brd.Nodes[MakeNodeLoc(0, rowSize-1)].inList = ArcIdx(LowerEndPt)
	} else { // N by M Board, N and M >= 2
		// Set the corner points
		brd.Nodes[MakeNodeLoc(0, 0)].inList = ArcIdx(UpperLeftCornerPt)
		brd.Nodes[MakeNodeLoc(colSize-1, 0)].inList = ArcIdx(UpperRightCornerPt)
		brd.Nodes[MakeNodeLoc(colSize-1, rowSize-1)].inList = ArcIdx(LowerRightCornerPt)
		brd.Nodes[MakeNodeLoc(0, rowSize-1)].inList = ArcIdx(LowerLeftCornerPt)
		// Set the edge points
		for r = 1; r < rowSize-1; r++ {
			brd.Nodes[MakeNodeLoc(0, r)].inList = ArcIdx(LeftEdgePt)
			brd.Nodes[MakeNodeLoc(colSize-1, r)].inList = ArcIdx(RightEdgePt)
		}
		for c = 1; c < colSize-1; c++ {
			brd.Nodes[MakeNodeLoc(c, 0)].inList = ArcIdx(UpperEdgePt)
			brd.Nodes[MakeNodeLoc(c, rowSize-1)].inList = ArcIdx(LowerEdgePt)
		}
		if (rowSize >= 5) && (colSize >= 5) {
			// Set the 2-2 points
			brd.Nodes[MakeNodeLoc(1, 1)].inList = ArcIdx(Corner_2_2_Pt)
			brd.Nodes[MakeNodeLoc(colSize-2, 1)].inList = ArcIdx(Corner_2_2_Pt)
			brd.Nodes[MakeNodeLoc(colSize-2, rowSize-2)].inList = ArcIdx(Corner_2_2_Pt)
			brd.Nodes[MakeNodeLoc(1, rowSize-2)].inList = ArcIdx(Corner_2_2_Pt)
			// Set the Line 2 points
			for r = 2; r < rowSize-2; r++ {
				brd.Nodes[MakeNodeLoc(1, r)].inList = ArcIdx(Line_2_Pt)
				brd.Nodes[MakeNodeLoc(colSize-2, r)].inList = ArcIdx(Line_2_Pt)
			}
			for c = 2; c < colSize-2; c++ {
				brd.Nodes[MakeNodeLoc(c, 1)].inList = ArcIdx(Line_2_Pt)
				brd.Nodes[MakeNodeLoc(c, rowSize-2)].inList = ArcIdx(Line_2_Pt)
			}
		}
		if (rowSize >= 7) && (colSize >= 7) {
			// Set the 3-3 points
			brd.Nodes[MakeNodeLoc(2, 2)].inList = ArcIdx(Corner_3_3_Pt)
			brd.Nodes[MakeNodeLoc(colSize-3, 2)].inList = ArcIdx(Corner_3_3_Pt)
			brd.Nodes[MakeNodeLoc(colSize-3, rowSize-3)].inList = ArcIdx(Corner_3_3_Pt)
			brd.Nodes[MakeNodeLoc(2, rowSize-3)].inList = ArcIdx(Corner_3_3_Pt)
			// Set the Line 3 points
			for r = 3; r < rowSize-3; r++ {
				brd.Nodes[MakeNodeLoc(2, r)].inList = ArcIdx(Line_3_Pt)
				brd.Nodes[MakeNodeLoc(colSize-3, r)].inList = ArcIdx(Line_3_Pt)
			}
			for c = 3; c < colSize-3; c++ {
				brd.Nodes[MakeNodeLoc(c, 2)].inList = ArcIdx(Line_3_Pt)
				brd.Nodes[MakeNodeLoc(c, rowSize-3)].inList = ArcIdx(Line_3_Pt)
			}
		}
		if (rowSize >= 9) && (colSize >= 9) {
			// Set the 4-4 points
			brd.Nodes[MakeNodeLoc(3, 3)].inList = ArcIdx(Corner_4_4_Pt)
			brd.Nodes[MakeNodeLoc(colSize-4, 3)].inList = ArcIdx(Corner_4_4_Pt)
			brd.Nodes[MakeNodeLoc(colSize-4, rowSize-4)].inList = ArcIdx(Corner_4_4_Pt)
			brd.Nodes[MakeNodeLoc(3, rowSize-4)].inList = ArcIdx(Corner_4_4_Pt)
			// Set the Line 4 points
			for r = 4; r < rowSize-4; r++ {
				brd.Nodes[MakeNodeLoc(3, r)].inList = ArcIdx(Line_4_Pt)
				brd.Nodes[MakeNodeLoc(colSize-4, r)].inList = ArcIdx(Line_4_Pt)
			}
			for c = 4; c < colSize-4; c++ {
				brd.Nodes[MakeNodeLoc(c, 3)].inList = ArcIdx(Line_4_Pt)
				brd.Nodes[MakeNodeLoc(c, rowSize-4)].inList = ArcIdx(Line_4_Pt)
			}
		}
		if (rowSize >= 11) && (colSize >= 11) {
			// Set the 5-5 points
			brd.Nodes[MakeNodeLoc(4, 4)].inList = ArcIdx(Corner_5_5_Pt)
			brd.Nodes[MakeNodeLoc(colSize-5, 4)].inList = ArcIdx(Corner_5_5_Pt)
			brd.Nodes[MakeNodeLoc(colSize-5, rowSize-5)].inList = ArcIdx(Corner_5_5_Pt)
			brd.Nodes[MakeNodeLoc(4, rowSize-5)].inList = ArcIdx(Corner_5_5_Pt)
			// Set the Line 5 points
			for r = 5; r < rowSize-5; r++ {
				brd.Nodes[MakeNodeLoc(4, r)].inList = ArcIdx(Line_5_Pt)
				brd.Nodes[MakeNodeLoc(colSize-5, r)].inList = ArcIdx(Line_5_Pt)
			}
			for c = 5; c < colSize-5; c++ {
				brd.Nodes[MakeNodeLoc(c, 4)].inList = ArcIdx(Line_5_Pt)
				brd.Nodes[MakeNodeLoc(c, rowSize-5)].inList = ArcIdx(Line_5_Pt)
			}
		}
		if (rowSize >= 13) && (colSize >= 13) {
			// Set the 6-6 points
			brd.Nodes[MakeNodeLoc(5, 5)].inList = ArcIdx(Corner_6_6_Pt)
			brd.Nodes[MakeNodeLoc(colSize-6, 5)].inList = ArcIdx(Corner_6_6_Pt)
			brd.Nodes[MakeNodeLoc(colSize-6, rowSize-6)].inList = ArcIdx(Corner_6_6_Pt)
			brd.Nodes[MakeNodeLoc(5, rowSize-6)].inList = ArcIdx(Corner_6_6_Pt)
			// Set the Line 6 points
			for r = 6; r < rowSize-6; r++ {
				brd.Nodes[MakeNodeLoc(5, r)].inList = ArcIdx(Line_6_Pt)
				brd.Nodes[MakeNodeLoc(colSize-6, r)].inList = ArcIdx(Line_6_Pt)
			}
			for c = 6; c < colSize-6; c++ {
				brd.Nodes[MakeNodeLoc(c, 5)].inList = ArcIdx(Line_6_Pt)
				brd.Nodes[MakeNodeLoc(c, rowSize-6)].inList = ArcIdx(Line_6_Pt)
			}
		}
		if (rowSize >= 15) && (colSize >= 15) {
			// Set the 7-7 points
			brd.Nodes[MakeNodeLoc(6, 6)].inList = ArcIdx(Corner_7_7_Pt)
			brd.Nodes[MakeNodeLoc(colSize-7, 6)].inList = ArcIdx(Corner_7_7_Pt)
			brd.Nodes[MakeNodeLoc(colSize-7, rowSize-7)].inList = ArcIdx(Corner_7_7_Pt)
			brd.Nodes[MakeNodeLoc(6, rowSize-7)].inList = ArcIdx(Corner_7_7_Pt)
			// Set the Line 7 points
			for r = 7; r < rowSize-7; r++ {
				brd.Nodes[MakeNodeLoc(6, r)].inList = ArcIdx(Line_7_Pt)
				brd.Nodes[MakeNodeLoc(colSize-7, r)].inList = ArcIdx(Line_7_Pt)
			}
			for c = 7; c < colSize-7; c++ {
				brd.Nodes[MakeNodeLoc(c, 6)].inList = ArcIdx(Line_7_Pt)
				brd.Nodes[MakeNodeLoc(c, rowSize-7)].inList = ArcIdx(Line_7_Pt)
			}
		}
	}
}

// AddGraphNode uses a variable sized array of GraphNodes,
// together with the freeNodes list of deleted Nodes.
func (g *Graph) AddGraphNode(ns uint16) (newN NodeLoc) {
	if g.freeNodes == NilNodeLoc {
		var newNode GraphNode
		ln := len(g.Nodes)
		g.Nodes = append(g.Nodes, newNode)
		newN = NodeLoc(ln)
		/*
			if ln == cap(g.Nodes) { // reallocate
				newSz := (ln + 1) * 2
				if newSz < 32 { // avoid small allocations
					newSz = 32
				}
				newN := make([]GraphNode, newSz)
				for i, curN := range g.Nodes {
					newN[i] = curN
				}
				g.Nodes = newN
			}
			g.Nodes = g.Nodes[0: ln+1]
		*/
	} else {
		newN = g.freeNodes
		g.freeNodes = g.Nodes[newN].memberOf
		g.Nodes[newN].mark = 0
	}
	g.Nodes[newN].lowState = ns
	g.Nodes[newN].highState = 0
	g.Nodes[newN].memberOf = NilNodeLoc
	g.Nodes[newN].nextSame = NilNodeLoc
	g.Nodes[newN].inList = nilArc
	g.Nodes[newN].outList = nilArc
	g.Nodes[newN].firstMem = NilNodeLoc
	g.Nodes[newN].mark = 0
	return newN
}

// AddEdge uses a variable sized array of GraphArc,
// and the freeArcs list of deleted edges
func (g *Graph) AddEdge(n1 NodeLoc, n2 NodeLoc) ([]GraphArc, ArcIdx) {
	var newA ArcIdx = nilArc
	if TraceAH {
		fmt.Println("Adding Edge: Level", g.gLevel, "NodeLoc", n1, "NodeLoc2", n2)
	}
	if n1 > n2 {
		n1, n2 = n2, n1
	}
	if g.freeArcs == nilArc {
		var newArc GraphArc
		ln := len(g.arcs)
		g.arcs = append(g.arcs, newArc)
		/*
					if ln == cap(g.arcs) { // reallocate
						newSz := (ln + 1) * 2
						if newSz < 16 { // avoid small allocations
							newSz = 16
						}
						if TraceAH {
							fmt.Println("Reallocating arcs: ", newSz)
						}
						newA := make([]GraphArc, newSz)
						for i, curA := range g.arcs {
							newA[i] = curA
						}
						g.arcs = newA
					}
					if TraceAH {
			//			fmt.Println("Extending arcs: ", ln+1)
					}
					g.arcs = g.arcs[0: ln+1]
		*/
		newA = ArcIdx(ln)
	} else {
		if TraceAH {
			//			fmt.Println("Reusing freeArcs: ", g.freeArcs)
		}
		newA = g.freeArcs
		g.freeArcs = g.arcs[newA].inNext
		g.arcs[newA].inNext = nilArc // TODO: should not need?
	}
	if TraceAH {
		fmt.Println("Added:", newA)
	}
	g.arcs[newA].fromNode = n1
	g.arcs[newA].toNode = n2
	g.arcs[newA].inNext = g.Nodes[n2].inList
	g.Nodes[n2].inList = newA
	g.arcs[newA].outNext = g.Nodes[n1].outList
	g.Nodes[n1].outList = newA
	g.arcs[newA].imageCount = 1 // TODO: or leave as 0 based, i.e. excess count > 0
	return g.arcs, newA
}

// NodeLocFunc is the type of functions for EachAdjNode
type NodeLocFunc func(NodeLoc)

// ArcFunc is the type of functions for EachIncidentArc.
// TODO: consider adding a parameter indicating the arc direction?
type ArcFunc func(ArcIdx)

// EachIncidentArc visits each arc connected to a given node.
func (g *Graph) EachIncidentArc(n NodeLoc, visit ArcFunc) {
	a := g.Nodes[n].inList
	for a != nilArc {
		visit(a)
		a = g.arcs[a].inNext
	}
	a = g.Nodes[n].outList
	for a != nilArc {
		visit(a)
		a = g.arcs[a].outNext
	}
}

// FindEdge finds an edge (undirected arc) between two Nodes.
// edges are created with fromNode < toNode
func (g *Graph) FindEdge(frN NodeLoc, toN NodeLoc) ArcIdx {
	if TraceAH {
		//		fmt.Println("Finding: Level", g.gLevel, "NodeLoc", frN,
		//          "NodeLoc2", toN)
	}
	ret := nilArc
	if frN > toN { // reverse to assure fromNode < toNode
		frN, toN = toN, frN
		if TraceAH {
			//			fmt.Println("Swapped: Level", g.gLevel, "NodeLoc", frN,
			//              "NodeLoc2", toN)
		}
	}
	aIn := g.Nodes[toN].inList
	if aIn != nilArc {
		count := 0
		firstLoopArc := nilArc
	Loop:
		for {
			if aIn == nilArc {
				break Loop
			}
			nextAIn := g.arcs[aIn].inNext
			if frN == g.arcs[aIn].fromNode {
				ret = aIn
				break Loop
			}
			count++
			if count > MAX_ARC_LIST {
				fmt.Println("loop in arc list, aIn:", aIn,
					"frN:", g.arcs[aIn].fromNode,
					"toN:", g.arcs[aIn].toNode)
				if firstLoopArc == nilArc {
					firstLoopArc = aIn
				} else {
					if aIn == firstLoopArc {
						break Loop
					}
				}
			}
			aIn = nextAIn
		}
	}
	if TraceAH {
		fmt.Println("Found: ", ret)
	}
	return ret
}

// DeleteEdge deletes and edge from the graph,
// and puts it on the freeArcs list.
// The freeArcs list is linked by the inNext field.
func (g *Graph) DeleteEdge(frN NodeLoc, toN NodeLoc) {
	if frN > toN { // reverse to assure fromNode < toNode
		frN, toN = toN, frN
	}
	// delete from inList
	a := g.Nodes[toN].inList
	if g.arcs[a].fromNode == frN {
		g.Nodes[toN].inList = g.arcs[a].inNext
	} else {
	Loop:
		for {
			preA := a
			a = g.arcs[a].inNext
			if a == nilArc {
				FatalAbstHierError("DeleteEdge: arc not on inList")
			} else {
				if g.arcs[a].fromNode == frN {
					g.arcs[preA].inNext = g.arcs[a].inNext
					break Loop
				}
			}
		}
	}
	// delete from outList
	a = g.Nodes[frN].outList
	if g.arcs[a].toNode == toN {
		g.Nodes[frN].outList = g.arcs[a].outNext
	} else {
	Loop2:
		for {
			preA := a
			a = g.arcs[a].outNext
			if a == nilArc {
				FatalAbstHierError("DeleteEdge: arc not on outList")
			} else {
				if g.arcs[a].toNode == toN {
					g.arcs[preA].outNext = g.arcs[a].outNext
					break Loop2
				}
			}
		}
	}
	// put on avail list
	g.arcs[a].inNext = g.freeArcs
	g.freeArcs = a
	g.arcs[a].outNext = nilArc

}

// DeleteNode deletes an isolated node from the graph,
// and puts it on the freeNodes list.
// The freeNodes list is linked by the memberOf field.
func (g *Graph) DeleteNode(n NodeLoc) {
	if g.Nodes[n].inList != nilArc {
		FatalAbstHierError("DeleteNode: inList not empty")
	}
	if g.Nodes[n].outList != nilArc {
		FatalAbstHierError("DeleteNode: outList not empty")
	}
	if g.Nodes[n].mark != 0 {
		FatalAbstHierError("DeleteNode: freeing a marked node")
	}
	// put on avail list, and mark it as free
	g.Nodes[n].mark = OnFreeList
	g.Nodes[n].memberOf = g.freeNodes
	g.freeNodes = n
}

// GetPoint returns a pointer to the BoardPoint at [c][r]
func (gph *Graph) GetPoint(c ColValue, r RowValue) *GraphNode {
	// TODO: add a check for brd.OnBoard(c, r) ?
	// or assure that GetPoint is only called from context where c and r are valid ?
	nl := MakeNodeLoc(c, r)
	return &gph.Nodes[nl]
}
