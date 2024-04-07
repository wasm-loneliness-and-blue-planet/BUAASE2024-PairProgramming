package main

import (
	"fmt"
	// "syscall/js"
)

// func arrGetter(this js.Value, args []js.Value) interface{} {
// 	array := args[0]
// 	index := args[1].Int()
// 	return array.Index(index)
// }

// func objGetter(this js.Value, args []js.Value) interface{} {
// 	obj := args[0]
// 	key := args[1].String()
// 	return obj.Get(key)
// }

func main() {
	// wait := make(chan struct{}, 0)
	// js.Global().Set("arrGetter", js.FuncOf(arrGetter))
	// js.Global().Set("objGetter", js.FuncOf(objGetter))
	// <-wait
	m := mancalaGame{}
	m.print()
}

func NewMancalaGame() mancalaGame {
	return mancalaGame{
		WhoseTurn: 0,
		Boards:    [2]mancalaBoard{NewMancalaBoard(), NewMancalaBoard()},
	}
}

func NewMancalaBoard() mancalaBoard {
	return mancalaBoard{
		Holes: [6]int{4, 4, 4, 4, 4, 4},
		Store: 0,
	}
}

// func MancalaResult(this js.Value, args []js.Value) interface{} {
// 	flag := args[0].Int()
// 	size := args[2].Int()
// 	seq := make([]int, size)
// 	for i := 0; i < size; i++ {
// 		seq[i] = args[1].Index(i).Int()
// 	}
// 	result := mancalaResult(flag, seq, size)
// 	return js.ValueOf(result)
// }

type mancalaBoard struct {
	Holes [6]int
	Store int
}

type mancalaGame struct {
	Boards    [2]mancalaBoard
	WhoseTurn int
}

func (m *mancalaGame) print() {
	fmt.Println("------------------------------")
	fmt.Print(m.Boards[1].Store)
	fmt.Print("  ")
	reverse := make([]int, 6)
	for i := 0; i < 6; i++ {
		reverse[i] = m.Boards[1].Holes[5-i]
	}
	fmt.Println(reverse)
	fmt.Println("------------------------------")
	fmt.Print("  ")
	fmt.Println(m.Boards[0])
	fmt.Println("------------------------------")
}

func reverseArray(arr []int) {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func (m *mancalaGame) sow(player int, holeIndex int, count int) {
	for count > 0 {
		m.Boards[player].Holes[holeIndex]++
		count--
		if holeIndex == 5 {
			player = 1 - player
			holeIndex = 0
		} else {
			holeIndex++
		}
	}
}

func (m *mancalaGame) play(step int) error {
	// TODO: check if valid
	// the player choose one hole and sow the stones
	player := step / 10
	holeIndex := step % 10
	if holeIndex < 0 || holeIndex >= 6 {
		return fmt.Errorf("holeIndex out of range")
	}
	count := m.Boards[player].Holes[holeIndex]
	if count == 0 {
		return fmt.Errorf("holeIndex is empty")
	}
	m.Boards[player].Holes[holeIndex] = 0

	return nil
}

func mancalaResult(flag int, seq []int, size int) int {
	// 1. Detect if the seq is following the rules

	// 2. Check if the game ends

	return 0
}
