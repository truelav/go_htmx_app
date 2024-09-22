package main

import (
	"strings"
)

func convertDoneToBool(isDone string) bool {
	var taskIsDone bool

	switch strings.ToLower(isDone) {
	case "yes", "on":
		taskIsDone = true
	case "no", "off":
		taskIsDone = false
	default:
		taskIsDone = false
	}

	return taskIsDone
}

// func experiment() {
// 	nums := []int{1, 2, 3}
// 	sum := 0

// 	for _, num := range nums {
// 		sum += num
// 	}

// 	fmt.Println(sum)

// 	for i, num := range nums {
// 		if num%2 == 0 {
// 			fmt.Println(i)
// 		}
// 	}

// 	kvs := map[string]string{"a": "1"}
// 	for k, v := range kvs {
// 		fmt.Printf("%s -> %s\n", k, v)
// 	}
// }

// func helloWorld() string {
// 	return "Hello World"
// }

// ////////////////////////////////////////////////

// type geometry interface {
// 	area() float64
// 	perimeter() float64
// }

// type circle struct {
// 	radius float64
// }

// type rectangle struct {
// 	width  float64
// 	height float64
// }

// func (r *rectangle) area() float64 {
// 	return r.height * r.width
// }

// func (r *rectangle) perimeter() float64 {
// 	return 2*r.height + 2*r.width
// }

// func (c *circle) area() float64 {
// 	return math.Pi * c.radius * c.radius
// }

// func (c *circle) perimeter() float64 {
// 	return 2 * math.Pi * c.radius
// }

// func measureGeometry(g geometry) {

// }
