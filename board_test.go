package ah_test

import (
	//    "testing"
	. "."
	"fmt"
)

func ExampleDirection() {
	for d := 0; d < int(RightDir)*2; d++ {
		fmt.Println(Direction(d))
	}
	// Output:
	// No Direction
	// Upper Direction
	// Left Direction
	// Upper Left Directions
	// Lower Direction
	// Upper Lower Directions
	// Left Lower Directions
	// Upper Left Lower Directions
	// Right Direction
	// Upper Right Directions
	// Left Right Directions
	// Upper Left Right Directions
	// Lower Right Directions
	// Upper Lower Right Directions
	// Left Lower Right Directions
	// Upper Left Lower Right Directions
}

func ExamplePointType() {
	// Check  PointType definitions:
loop:
	for pt := SingletonPt; pt <= Line_7_Pt; pt++ {
		switch pt {
		case SingletonPt:
			fmt.Println(int(pt), "SingletonPt:", pt)

		case LowerEndPt:
			fmt.Println(int(pt), "LowerEndPt:", pt)
		case RightEndPt:
			fmt.Println(int(pt), "RightEndPt:", pt)
		case UpperEndPt:
			fmt.Println(int(pt), "UpperEndPt:", pt)
		case LeftEndPt:
			fmt.Println(int(pt), "LeftEndPt:", pt)

		case UpperLowerBridgePt:
			fmt.Println(int(pt), "UpperLowerBridgePt:", pt)
		case LeftRightBridgePt:
			fmt.Println(int(pt), "LeftRightBridgePt:", pt)

		case UpperLeftCornerPt:
			fmt.Println(int(pt), "UpperLeftCornerPt:", pt)
		case UpperRightCornerPt:
			fmt.Println(int(pt), "UpperRightCornerPt:", pt)
		case LowerRightCornerPt:
			fmt.Println(int(pt), "LowerRightCornerPt:", pt)
		case LowerLeftCornerPt:
			fmt.Println(int(pt), "LowerLeftCornerPt:", pt)

		case UpperEdgePt:
			fmt.Println(int(pt), "UpperEdgePt:", pt)
		case LeftEdgePt:
			fmt.Println(int(pt), "LeftEdgePt:", pt)
		case LowerEdgePt:
			fmt.Println(int(pt), "LowerEdgePt:", pt)
		case RightEdgePt:
			fmt.Println(int(pt), "RightEdgePt:", pt)

		case CenterPt:
			fmt.Println(int(pt), "CenterPt:", pt)

		case HoshiPt:
			fmt.Println(int(pt), "HoshiPt:", pt)

		case Corner_2_2_Pt:
			fmt.Println(int(pt), "Corner_2_2_Pt:", pt)
		case Line_2_Pt:
			fmt.Println(int(pt), "Line_2_Pt:", pt)
		case Corner_3_3_Pt:
			fmt.Println(int(pt), "Corner_3_3_Pt:", pt)
		case Line_3_Pt:
			fmt.Println(int(pt), "Line_3_Pt:", pt)
		case Corner_4_4_Pt:
			fmt.Println(int(pt), "Corner_4_4_Pt:", pt)
		case Line_4_Pt:
			fmt.Println(int(pt), "Line_4_Pt:", pt)
		case Corner_5_5_Pt:
			fmt.Println(int(pt), "Corner_5_5_Pt:", pt)
		case Line_5_Pt:
			fmt.Println(int(pt), "Line_5_Pt:", pt)
		case Corner_6_6_Pt:
			fmt.Println(int(pt), "Corner_6_6_Pt:", pt)
		case Line_6_Pt:
			fmt.Println(int(pt), "Line_6_Pt:", pt)
		case Corner_7_7_Pt:
			fmt.Println(int(pt), "Corner_7_7_Pt:", pt)
		case Line_7_Pt:
			fmt.Println(int(pt), "Line_7_Pt:", pt)

			break loop // loop will not terminate 255++ => 0

		// case UninitializedPt:
		//      fmt.Println("UninitializedPt:", pt)

		default:
			//				fmt.Println("skipping:", pt)
		}
	}
	// Output:
	// 0 SingletonPt: ·
	// 1 LowerEndPt: ╹
	// 2 RightEndPt: ╸
	// 3 LowerRightCornerPt: ┛
	// 4 UpperEndPt: ╻
	// 5 UpperLowerBridgePt: ┃
	// 6 UpperRightCornerPt: ┓
	// 7 RightEdgePt: ┨
	// 8 LeftEndPt: ╺
	// 9 LowerLeftCornerPt: ┗
	// 10 LeftRightBridgePt: ━
	// 11 LowerEdgePt: ┷
	// 12 UpperLeftCornerPt: ┏
	// 13 LeftEdgePt: ┠
	// 14 UpperEdgePt: ┯
	// 15 CenterPt: ╋
	// 31 HoshiPt: ◘
	// 47 Corner_2_2_Pt: ╬
	// 63 Line_2_Pt: ┼
	// 79 Corner_3_3_Pt: ╬
	// 95 Line_3_Pt: ┼
	// 111 Corner_4_4_Pt: ╬
	// 127 Line_4_Pt: ┼
	// 143 Corner_5_5_Pt: ╬
	// 159 Line_5_Pt: ┼
	// 175 Corner_6_6_Pt: ╬
	// 191 Line_6_Pt: ┼
	// 207 Corner_7_7_Pt: ╬
	// 223 Line_7_Pt: ┼
}

func ExamplePointStatus() {
	for ps := UndefinedPointStatus; ps <= LastPointStatus; ps++ {
		switch ps {

		case UndefinedPointStatus,
			Black, White,
			AB_U, AB_W, AE_B, AE_W, AW_B, AW_U,
			// Unoccupied, generic value
			Unocc,
			// No Adjacent Stones:
			B0W0,
			// Single Adjacent Stone:
			W1, B1,
			// Two Adjacent Stones:
			W2, B1W1, W1B1, B2,
			// Three Adjacent Stones:
			B3, B2W1, B1W2, W3, WBB, WBW, BWB, W2B1,
			// Four Adjacent Stones:
			B4, B3W1, BWBW, BBWW, B1W3, W4,

			LastPointStatus:
			fmt.Println(ps, int(ps))

		default:
			//				fmt.Println("skipping:", ps)
		}
	}
	// Output:
	// UndefinedPointStatus 0
	// Black 1
	// White 2
	// AB_U 3
	// AB_W 4
	// AE_B 5
	// AE_W 6
	// AW_B 7
	// AW_U 8
	// Unocc 9
	// B0W0 10
	// W1 11
	// B1 12
	// W2 13
	// B1W1 14
	// W1B1 15
	// B2 16
	// B3 17
	// B2W1 18
	// B1W2 19
	// W3 20
	// WBB 21
	// WBW 22
	// BWB 23
	// W2B1 24
	// B4 25
	// B3W1 26
	// BWBW 27
	// BBWW 28
	// B1W3 29
	// W4 30
	// LastPointStatus 31
}

func ExampleAhTypeSizes() {
	PrintAhTypeSizes()
	// Output:
	// Type Direction size 1 alignment 1
	// Type PointType size 1 alignment 1
	// Type ColValue size 1 alignment 1
	// Type RowValue size 1 alignment 1
	// Type ColSize size 1 alignment 1
	// Type RowSize size 1 alignment 1
	// Type NodeLoc size 2 alignment 2
	// Type NodeLocList size 24 alignment 8
	// Type MoveRecord size 14 alignment 2
	// Type Board size 112 alignment 8
	// Type SearchStack size 40 alignment 8
	// Type GraphNodeLocFunc size 8 alignment 8
	// Type NodeLocFuncBool size 8 alignment 8
	// Type BoardTrans size 1 alignment 1
	// Type ArcIdx size 2 alignment 2
	// Type GraphMark size 1 alignment 1
	// Type GraphNode size 16 alignment 2
	// Type GraphArc size 10 alignment 2
	// Type CompStateFunc size 8 alignment 8
	// Type ChangeRequest size 4 alignment 2
	// Type Graph size 104 alignment 8
	// Type NodeLocFunc size 8 alignment 8
	// Type ArcFunc size 8 alignment 8
	// Type PointStatus size 2 alignment 2
	// Type GraphLevel size 1 alignment 1
	// Type AbstHier size 800 alignment 8
	// Type NodeStatus size 2 alignment 2
	// Type StringStatus size 2 alignment 2
}
