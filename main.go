package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct {
	radius float64
}

type Rectangle struct {
	width, height float64
}

func (r Rectangle) Area() float64 {
	return r.width * r.height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.width + r.height)
}

func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.radius
}

func PrintShapeInfo(s Shape) {
	fmt.Printf("Area: %f, Perimeter: %f\n", s.Area(), s.Perimeter())
}

func main() {
	r := Rectangle{width: 2, height: 3}
	c := Circle{radius: 1}

	PrintShapeInfo(r) // Output: Area: 6.000000, Perimeter: 10.000000
	PrintShapeInfo(c) // Output: Area: 3.141593, Perimeter: 6.283185

}
