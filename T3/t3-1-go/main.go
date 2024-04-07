package main

import (
	"errors"
	"fmt"
	"syscall/js"
)

func main() {
	wait := make(chan struct{}, 0)
	js.Global().Set("mancalaBoard", js.FuncOf(MancalaBoardWrapper))
	<-wait
}

func NewMancalaGame(firstHand int) mancalaGame {
	return mancalaGame{
		WhoseTurn: firstHand - 1,
		Boards:    [2]mancalaBoard{NewMancalaBoard(), NewMancalaBoard()},
		FirstHand: firstHand,
		IsEnd:     false,
	}
}

func NewMancalaBoard() mancalaBoard {
	return mancalaBoard{
		Holes: [6]int{4, 4, 4, 4, 4, 4},
		Store: 0,
	}
}

func MancalaBoardWrapper(this js.Value, args []js.Value) interface{} {
	flag := args[0].Int()
	size := args[2].Int()
	seq := make([]int, size)
	// println(flag, size, seq)
	for i := 0; i < size; i++ {
		seq[i] = args[1].Index(i).Int()
	}
	result := mancalaMidBoard(flag, seq, size)
	return js.ValueOf(result)
}

type mancalaBoard struct {
	Holes [6]int
	Store int
}

type mancalaGame struct {
	Boards    [2]mancalaBoard
	WhoseTurn int // 此处 0 指先手，1 指后手，不是玩家 1 / 2
	FirstHand int // 此处 1 指玩家 1，2 指玩家 2
	IsEnd     bool
}

func (m *mancalaGame) print() {

	fmt.Println("-----------------------------")
	fmt.Print(m.Boards[1].Store)
	fmt.Print("  ")
	reverse := make([]int, 6)
	for i := 0; i < 6; i++ {
		reverse[i] = m.Boards[1].Holes[5-i]
	}
	fmt.Println(reverse)
	fmt.Print("-----------------------------\n   ")
	fmt.Println(m.Boards[0])
	fmt.Println("-----------------------------")
}

func reverseArray(arr []int) {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}

// simulates the process of sowing stones in the mancala game.
//
// Parameters:
// - player: 0 / 1
// - fromHoleIndex: 0 ~ 5
// Returns:
// - int: the next player, 0 / 1
// - error
func (m *mancalaGame) sow(player, fromHoleIndex int) (int, error) {
	// 1. Get the number of stones in the hole
	totalCount := m.Boards[player].Holes[fromHoleIndex]
	if totalCount == 0 {
		return -1, errors.New("no stones in the hole")
	}
	m.Boards[player].Holes[fromHoleIndex] = 0
	boardIndex := player
	holeIndex := fromHoleIndex + 1
	for count := totalCount; count > 0; count-- {
		board := &m.Boards[boardIndex]
		if holeIndex == 6 {
			if boardIndex == player { // 对于己方计分洞需要分配棋子
				board.Store++
				if count == 1 { // 如果正好落在己方记分洞中，则当前玩家再下一次
					return boardIndex, nil
				}
			} // 对于对方计分洞直接跳过
			boardIndex = 1 - boardIndex // 切换到对方
			holeIndex = 0
		} else if count == 1 && board.Holes[holeIndex] == 0 && boardIndex == player {
			// **取子：**如果**最后播撒下的一颗棋子(count = 1)**落在**己方<u>无棋子</u>的棋洞**中，且该棋洞**正对面对方的棋洞**中**<u>有棋子</u>**，则将**这最后一颗棋子**和**正对面对方坑洞内的<u>所有</u>棋子**都放入**己方计分洞**，也就是说全部成为自己的得分。
			opponentBoard := &m.Boards[1-boardIndex]
			board.Holes[holeIndex] += opponentBoard.Holes[5-holeIndex] + 1
			opponentBoard.Holes[5-holeIndex] = 0
		} else {
			board.Holes[holeIndex]++
			holeIndex++
		}
	}
	return 1 - player, nil
}

func (m *mancalaGame) playOneStep(step int) error {
	player := step/10 - 1
	// if m.FirstHand == 1 {
	// 	// 1, 2 => 0, 1
	// 	player -= 1
	// } else /* m.FirstHand = 2 */ if player == 2 {
	// 	// 1, 2 => 1, 0
	// 	player = 0
	// }
	if player != m.WhoseTurn {
		return errors.New("t3-1-go Wrong player")
	}
	fromHoleIndex := step%10 - 1
	// println("Table playing from ", player, fromHoleIndex)
	nextPlayer, err := m.sow(player, fromHoleIndex)
	if err != nil {
		return err
	}
	m.WhoseTurn = nextPlayer
	return nil
}

func (m *mancalaGame) checkEnd() bool {
	// 6. **游戏结束：**有一方的所有棋洞中都没有棋子时，游戏结束。此时，**所有玩家不能再进行操作**。另一方的棋洞中仍有棋子，**这些棋子全部放到己方的计分洞中**，即作为**仍有棋子的这一方的得分**的一部分。
	isEnd := false
	for i := 0; i < 2; i++ {
		empty := true
		for _, count := range m.Boards[i].Holes {
			if count > 0 {
				empty = false
				break
			}
		}
		if empty {
			m.print()
			println("player", i+1, "has no stones.")
			isEnd = true
			break
		}
	}
	if !isEnd {
		return false
	}
	// 计算得分
	score := [2]int{0, 0}
	for i := 0; i < 2; i++ {
		for j := 0; j < 6; j++ {
			score[i] += m.Boards[i].Holes[j]
		}
	}
	if score[0] == 0 {
		m.Boards[0].Store += score[1]
	} else {
		m.Boards[1].Store += score[0]
	}
	m.IsEnd = true
	return true
}

// 返回值数组的元素个数为 15，每位数字的含义如下：
//
//	| 位 0 - 5                | 位 6                    | 位 7 - 12               | 位 13                   | 位 14  |
//	| ----------------------- | ----------------------- | ----------------------- | ----------------------- | ------ |
//	| 棋洞 11 - 16 中的棋子数 | 选手 1 计分洞中的棋子数 | 棋洞 21 - 26 中的棋子数 | 选手 2 计分洞中的棋子数 | 数据位 |
func mancalaMidBoard(flag int, seq []int, size int) []interface{} {
	firstHand := seq[0] / 10
	m := NewMancalaGame(firstHand)
	// return a JS array
	res := make([]interface{}, 15)
	invalid := false
	for i, v := range seq {
		// 1. Detect if the seq is following the rules
		err := m.playOneStep(v)
		if err != nil {
			println(err)
			invalid = true
			break
		}

		// 2. Check if the game ends
		if m.checkEnd() {
			if i == size-1 {
				// ending
			} else {
				invalid = true
			}
			break
		}
	}
	m.print()
	for j := 0; j < 6; j++ {
		res[j] = m.Boards[0].Holes[j]
		res[j+7] = m.Boards[1].Holes[j]
	}
	res[6] = m.Boards[0].Store
	res[13] = m.Boards[1].Store

	if m.IsEnd && !invalid {
		netScore := m.Boards[0].Store - m.Boards[1].Store
		println("Game end, net score", netScore)
		res[14] = 200 + netScore
	} else if !invalid { // not ending
		res[14] = m.getHander()
	} else /* invalid = true */ if flag == 1 {
		res[14] = 200 + 2*m.Boards[0].Store - 48
	} else {
		res[14] = 200 + 48 - 2*m.Boards[1].Store
	}
	return res
}

func (m *mancalaGame) getHander() int {
	return m.WhoseTurn + 1
}
