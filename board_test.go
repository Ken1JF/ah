package ah

import (
    //    "testing"
    "fmt"
)

func ExampleDirection() {
    for d := 0; d < int(RightDir)*2 ; d++ {
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
                //			case Black_Occ_Pt:
                //				fmt.Println("Black_Occ_Pt:", pt)
                //			case White_Occ_Pt:
                //				fmt.Println("White_Occ_Pt:", pt)
			break loop // loop will not terminate 255++ => 0
            
                       //			case UninitializedPt:
                       //				fmt.Println("UninitializedPt:", pt)
            
            default:
                //				fmt.Println(" skipping:", pt)
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
    // 223 Line_7_Pt: ¿
}

func ExamplePointStatus() {
    
    
}
