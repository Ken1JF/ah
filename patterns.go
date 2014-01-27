/*
 *	File:		src/gitHub.com/Ken1JF/ah/patterns.go
 *  Project:	abst-hier
 *
 *  Created by Ken Friedenbach on 12/25/2010.
 *  Copyright 2010-2011 Ken Friedenbach. All rights reserved.
 *
 *	This file supports the building, reading, and writing of Libraries of Patterns.
 */

package ah

type PatternType uint8

// Define the types of patterns on a Board:
const (
	UNKN_PATTERN PatternType = iota
	WHOLE_BOARD_PATTERN
	CORNER_PATTERN
	HALF_BOARD_PATTERN
	SIDE_PATTERN
	CENTER_PATTERN
)

// names for printing:
var PatternNames = [6]string{
	"Unknown PatternType",
	"Whole_Board",
	"Corner",
	"Half_Board",
	"Side",
	"Center",
}

type PointSymmetry uint8

// Define the types of symmetry of each point:
const (
	UNKN_SYMMETRY        PointSymmetry = iota
	CENTER_SYMMETRY                    // Center point of odd sized square board, all eight transformations leave point unchanged
	DIAGONAL_SYMMETRY                  // Points on diagonal line, one transformation leaves point unchanged, four images
	CENTER_LINE_SYMMETRY               // Points on one center line, one transformation leaves point unchanged, four images
	NO_SYMMETRY                        // All trnaformations move point, eight images
)

type HandicapSymmetry uint8

const (
	UNKN_HANDI_SYMMETRY      HandicapSymmetry = iota
	EIGHT_WAY_SYMMETRY                        // one eighth of board + 7 imaage
	TWO_DIAGONAL_SYMMETRY                     // right side of board + 3 images
	TWO_CENTER_LINE_SYMMETRY                  // upper right corrner + 3 images
	FLIP_VERT_SYMMETRY                        // one half of board + 1 image
	FLIP_BACK_SYMMETRY                        // one half of board + 1 image
)

var BoardHandicapSymmetry = []HandicapSymmetry{
	EIGHT_WAY_SYMMETRY,       // HA = 0
	UNKN_HANDI_SYMMETRY,      // HA = 1, not used
	TWO_DIAGONAL_SYMMETRY,    // HA = 2
	FLIP_BACK_SYMMETRY,       // HA = 3, upper left corner empty (programs differ)
	EIGHT_WAY_SYMMETRY,       // HA = 4
	EIGHT_WAY_SYMMETRY,       // HA = 5
	TWO_CENTER_LINE_SYMMETRY, // HA = 6
	TWO_CENTER_LINE_SYMMETRY, // HA = 7
	EIGHT_WAY_SYMMETRY,       // HA = 8
	EIGHT_WAY_SYMMETRY,       // HA = 9
}

var TransPreservesHandicapPattern = [6][8]bool{
	/* UNKN_HANDI_SYMMETRY (not used) */
	{true /* T_IDENTITY */, true /* T_ROTA_090 */, true /* T_ROTA_180 */, true /* T_ROTA_270 */, true /* T_FLP_SLAS */, true /* T_FLP_VERT */, true /* T_FLP_BACK */, true /* T_FLP_HORI */},
	/* EIGHT_WAY_SYMMETRY */
	{true /* T_IDENTITY */, true /* T_ROTA_090 */, true /* T_ROTA_180 */, true /* T_ROTA_270 */, true /* T_FLP_SLAS */, true /* T_FLP_VERT */, true /* T_FLP_BACK */, true /* T_FLP_HORI */},
	/* TWO_DIAGONAL_SYMMETRY */
	{true /* T_IDENTITY */, false /* T_ROTA_090 */, true /* T_ROTA_180 */, false /* T_ROTA_270 */, true /* T_FLP_SLAS */, false /* T_FLP_VERT */, true /* T_FLP_BACK */, false /* T_FLP_HORI */},
	/* TWO_CENTER_LINE_SYMMETRY */
	{true /* T_IDENTITY */, false /* T_ROTA_090 */, true /* T_ROTA_180 */, false /* T_ROTA_270 */, false /* T_FLP_SLAS */, true /* T_FLP_VERT */, false /* T_FLP_BACK */, true /* T_FLP_HORI */},
	/* FLIP_VERT_SYMMETRY */
	{true /* T_IDENTITY */, false /* T_ROTA_090 */, false /* T_ROTA_180 */, false /* T_ROTA_270 */, false /* T_FLP_SLAS */, true /* T_FLP_VERT */, false /* T_FLP_BACK */, false /* T_FLP_HORI */},
	/* FLIP_BACK_SYMMETRY */
	{true /* T_IDENTITY */, false /* T_ROTA_090 */, false /* T_ROTA_180 */, false /* T_ROTA_270 */, false /* T_FLP_SLAS */, false /* T_FLP_VERT */, true /* T_FLP_BACK */, false /* T_FLP_HORI */},
}

// GetPointSymmetry returns the symmetry of a point:
//func (abhr *AbstHier) GetPointSymmetry(nl NodeLoc) (PointSymmetry) {
//	ncols, nrows := abhr.GetSize()
//	col, row := GetColRow(nl)
//}

// IsCanonical tests if a point is in the 1/8 th of the board with lowest coordinates:
func (abhr *AbstHier) IsCanonical(nl NodeLoc, hs HandicapSymmetry) (ret bool) {
	ncols, nrows := abhr.GetSize()
	midRow := RowValue(nrows-1) / 2
	midCol := ColValue(ncols-1) / 2
	col, row := GetColRow(nl)
	switch hs {
	case EIGHT_WAY_SYMMETRY: // one eighth of board + 7 imaage
		ret = (row <= midRow) && (col >= ColValue(nrows-RowSize(row+1))) // top half of board AND right of diagonal
	case TWO_DIAGONAL_SYMMETRY: // right side of board + 3 images
		if row <= midRow {
			ret = (col >= ColValue(nrows-RowSize(row+1)))
		} else {
			ret = (col >= ColValue(row))
		}
	case TWO_CENTER_LINE_SYMMETRY: // upper right corrner + 3 images
		ret = (row <= midRow) && (col >= midCol) // top half of board AND right half of board
	case FLIP_VERT_SYMMETRY: // one half of board + 1 image
		ret = (col >= ColValue(nrows-RowSize(row+1)))
	case FLIP_BACK_SYMMETRY:
		ret = (col >= ColValue(row))
	}
	return ret
}

// FindCanonicalRep returns the Canonical point and the first transformation
// that takes a point into its canonical representation
func (abhr *AbstHier) FindCanonicalRep(nl NodeLoc, hs HandicapSymmetry) (nlc NodeLoc, trans BoardTrans) {
	for trans := T_FIRST; trans <= T_LAST; trans += 1 {
		if TransPreservesHandicapPattern[hs][trans] {
			col, row := GetColRow(nl)
			nlc := abhr.TransNodeLoc(trans, col, row)
			if abhr.IsCanonical(nlc, hs) {
				return nlc, trans
			}
		}
	}
	FatalAbstHierError("FindCanonicalRep - not found")
	return nlc, trans
}
