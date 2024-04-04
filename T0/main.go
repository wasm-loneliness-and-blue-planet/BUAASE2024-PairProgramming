package main

var state int

func main() {
	state = 1
}

//export add
func add(x int, y int) int {
	return x + y
}

func Mul(x int, y int) int {
	return x * y
}

//export getState
func getState() int {
	return state
}

//export setState
func setState(x int) bool {
	state = x
	return true
}
