package rect

import "fmt"

type rect struct {
    width int
    height int
}

func (r *rect) area() int {
    return r.width * r.height
}

func Print() {
    r := rect{width: 10, height: 5}
    fmt.Println("area:", r.area()) // this adds an extra space
}
