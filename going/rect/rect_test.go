package rect_test

import "testing"

import(
  "fmt"
  "rect"
)

func TestRect(t *testing.T){
  fmt.Println("In test rect")
  rect.Print()
}

func ExamplePrefix() {
  fmt.Println("hello")
  rect.Print()
  // Output:
  // hello
  // area: 50
}
