package ah

import (
    //    "testing"
    "fmt"
)

func ExampleBoard() {
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