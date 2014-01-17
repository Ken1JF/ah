/*
 *  File:		src/gitHub.com/Ken1JF/ahgo/ah/trans.go
 *  Project:	abst-hier
 *
 *  Created by Ken Friedenbach on 2/16/10.
 *  Copyright 2005-2010 Ken Friedenbach. All rights reserved.
 *
 *	This file implements transformations of the board.
 */

package ah

type BoardTrans uint8

// Define the eight symmetric transformations of an N by N board:
//
const (
	T_IDENTITY BoardTrans = iota // Identity (Zero) transformation
	T_ROTA_090                   // Rotate  90 degrees counter-clockwise
	T_ROTA_180                   // Rotate 180 degrees counter-clockwise
	T_ROTA_270                   // Rotate 270 degrees counter-clockwise
	T_FLP_SLAS                   // Flip about "/" axis (Slash)
	T_FLP_VERT                   // Flip about "|" axis (Vertical)
	T_FLP_BACK                   // Flip about "\" axis (Backslash)
	T_FLP_HORI                   // Flip about "-" axis (Horizontal)
)

// constants for iteration:
//
const (
	T_FIRST   BoardTrans = T_IDENTITY
	T_LAST    BoardTrans = T_FLP_HORI
	NUM_TRANS BoardTrans = T_LAST + 1
)

// name for printing:
//
var TransName = [8]string{
	"T_IDENTITY",
	"T_ROTA_090",
	"T_ROTA_180",
	"T_ROTA_270",
	"T_FLP_SLAS",
	"T_FLP_VERT",
	"T_FLP_BACK",
	"T_FLP_HORI",
}

// inverse transformations:
//
var InverseTrans = [8]BoardTrans{
	T_IDENTITY, // T_IDENTITY inverse
	T_ROTA_270, // T_ROTA_090 inverse
	T_ROTA_180, // T_ROTA_180 inverse
	T_ROTA_090, // T_ROTA_270 inverse
	T_FLP_SLAS, // T_FLP_SLAS inverse
	T_FLP_VERT, // T_FLP_VERT inverse
	T_FLP_BACK, // T_FLP_BACK inverse
	T_FLP_HORI, // T_FLP_HORI inverse
}

// define the result of composing two board transformations:
// board_trans_compose[A][B] is result of applying A, then B.
// Note: composition is not commutative	for all transforms.
//
var ComposeTrans = [8][8]BoardTrans{
	//									T_IDENTITY	T_ROTA_090	T_ROTA_180	T_ROTA_270	T_FLP_SLAS	T_FLP_VERT	T_FLP_BACK	T_FLP_HORI
	//									---------- 	----------	----------	----------	----------	----------	----------	----------
	/*	T_IDENTITY] = */ [8]BoardTrans{T_IDENTITY, T_ROTA_090, T_ROTA_180, T_ROTA_270, T_FLP_SLAS, T_FLP_VERT, T_FLP_BACK, T_FLP_HORI},
	/*	T_ROTA_090] = */ [8]BoardTrans{T_ROTA_090, T_ROTA_180, T_ROTA_270, T_IDENTITY, T_FLP_HORI, T_FLP_SLAS, T_FLP_VERT, T_FLP_BACK},
	/*	T_ROTA_180] = */ [8]BoardTrans{T_ROTA_180, T_ROTA_270, T_IDENTITY, T_ROTA_090, T_FLP_BACK, T_FLP_HORI, T_FLP_SLAS, T_FLP_VERT},
	/*	T_ROTA_270] = */ [8]BoardTrans{T_ROTA_270, T_IDENTITY, T_ROTA_090, T_ROTA_180, T_FLP_VERT, T_FLP_BACK, T_FLP_HORI, T_FLP_SLAS},
	/*	T_FLP_SLAS] = */ [8]BoardTrans{T_FLP_SLAS, T_FLP_VERT, T_FLP_BACK, T_FLP_HORI, T_IDENTITY, T_ROTA_090, T_ROTA_180, T_ROTA_270},
	/*	T_FLP_VERT] = */ [8]BoardTrans{T_FLP_VERT, T_FLP_BACK, T_FLP_HORI, T_FLP_SLAS, T_ROTA_270, T_IDENTITY, T_ROTA_090, T_ROTA_180},
	/*	T_FLP_BACK] = */ [8]BoardTrans{T_FLP_BACK, T_FLP_HORI, T_FLP_SLAS, T_FLP_VERT, T_ROTA_180, T_ROTA_270, T_IDENTITY, T_ROTA_090},
	/*	T_FLP_HORI] = */ [8]BoardTrans{T_FLP_HORI, T_FLP_SLAS, T_FLP_VERT, T_FLP_BACK, T_ROTA_090, T_ROTA_180, T_ROTA_270, T_IDENTITY},
}

func (abhr *AbstHier) TransNodeLoc(t BoardTrans, c ColValue, r RowValue) (newNode NodeLoc) {
	newNode = PassNodeLoc
	nc, nr := abhr.GetSize()
	if (c < nc) && (r < nr) { // check if OnBoard
		switch t {
		case T_IDENTITY:
			newNode = MakeNodeLoc(c, r)
		case T_ROTA_090:
			newNode = MakeNodeLoc(ColValue(r), RowValue(nc-(c+1)))
		case T_ROTA_180:
			newNode = MakeNodeLoc(nc-(c+1), nr-(r+1))
		case T_ROTA_270:
			newNode = MakeNodeLoc(ColValue(nr-(r+1)), RowValue(c))
		case T_FLP_SLAS:
			newNode = MakeNodeLoc(ColValue(nr-(r+1)), RowValue(nc-(c+1)))
		case T_FLP_VERT:
			newNode = MakeNodeLoc(nc-(c+1), r)
		case T_FLP_BACK:
			newNode = MakeNodeLoc(ColValue(r), RowValue(c))
		case T_FLP_HORI:
			newNode = MakeNodeLoc(c, nr-(r+1))
		}
	}
	return newNode
}

func (abhr *AbstHier) TransBoard(t BoardTrans) *AbstHier {
	var col ColValue
	var row RowValue
	brd := &abhr.Graphs[PointLevel]
	newAH := new(AbstHier)
	// TODO: support transformations of N by M boards, N != M:
	// nCols, nRows = abhr.TransColsRows(t BoardTrans)
	nCols, nRows := abhr.GetSize()
	newAH = newAH.InitAbstHier(nCols, nRows, abhr.updtLev, true)
	newBrd := &newAH.Graphs[PointLevel]
	for row = 0; row < nRows; row++ {
		for col = 0; col < nCols; col++ {
			newNode := abhr.TransNodeLoc(t, col, row)
			newBrd.Nodes[newNode].SetNodeLowState(brd.Nodes[MakeNodeLoc(col, row)].GetNodeLowState())
		}
	}
	return newAH
}
