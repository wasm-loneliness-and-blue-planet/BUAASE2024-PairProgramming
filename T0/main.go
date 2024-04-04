package main

import (
	"strconv"
	"strings"
	"syscall/js"
)

var state int

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

func splitter(this js.Value, args []js.Value) interface{} {
	values := strings.Split(args[0].String(), ",")
	result := make([]interface{}, 0)
	for _, each := range values {
		integer, _ := strconv.Atoi(each)
		result = append(result, integer)
	}
	return js.ValueOf(result)
}

func splitter2(this js.Value, args []js.Value) interface{} {
	size := args[1].Int()
	values := strings.Split(args[0].String(), ",")
	result := make([]interface{}, size)
	for i, each := range values {
		integer, _ := strconv.Atoi(each)
		// result = append(result, integer)
		result[i] = integer
	}
	return js.ValueOf(result)
}

func arrGetter(this js.Value, args []js.Value) interface{} {
	array := args[0]
	index := args[1].Int()
	return array.Index(index)
}

func objGetter(this js.Value, args []js.Value) interface{} {
	obj := args[0]
	key := args[1].String()
	return obj.Get(key)
}

func main() {
	state = 233
	wait := make(chan struct{}, 0)
	js.Global().Set("splitter", js.FuncOf(splitter))
	js.Global().Set("splitter2", js.FuncOf(splitter2))
	js.Global().Set("arrGetter", js.FuncOf(arrGetter))
	js.Global().Set("objGetter", js.FuncOf(objGetter))
	<-wait
}
