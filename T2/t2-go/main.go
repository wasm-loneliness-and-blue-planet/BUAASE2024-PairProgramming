package main

import (
	"errors"
	"fmt"
	"syscall/js"
)

func main() {
	wait := make(chan struct{}, 0)
	js.Global().Set("mancalaResult", js.FuncOf(MancalaResult))
	<-wait
}

func NewMancalaGame(firstHand int) mancalaGame {
	return mancalaGame{
		WhoseTurn: 0,
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

func MancalaResult(this js.Value, args []js.Value) interface{} {
	flag := args[0].Int()
	size := args[2].Int()
	seq := make([]int, size)
	for i := 0; i < size; i++ {
		seq[i] = args[1].Index(i).Int()
	}
	result := mancalaResult(flag, seq, size)
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
	player := step / 10
	if m.FirstHand == 1 {
		// 1, 2 => 0, 1
		player -= 1
	} else if player == 2 {
		// 1, 2 => 1, 0
		player = 0
	}
	if player != m.WhoseTurn {
		return errors.New("Wrong player")
	}
	fromHoleIndex := step%10 - 1
	nextPlayer, err := m.sow(player, fromHoleIndex)
	if err != nil {
		return err
	}
	m.WhoseTurn = nextPlayer
	return nil
}

func (m *mancalaGame) checkEnd() bool {
	// 6. **游戏结束：**有一方的所有棋洞中都没有棋子时，游戏结束。此时，**所有玩家不能再进行操作**。另一方的棋洞中仍有棋子，**这些棋子全部放到己方的计分洞中**，即作为**仍有棋子的这一方的得分**的一部分。
	playersNoEmptyHoles := make([]bool, 2)
	isEnd := false
	for i := 0; i < 2; i++ {
		for j := 0; j < 6; j++ {
			if m.Boards[i].Holes[j] > 0 {
				playersNoEmptyHoles[i] = true
				break
			}
		}
		if !playersNoEmptyHoles[i] {
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

func (m *mancalaGame) getWinner() (int, int) {
	netScore := m.Boards[0].Store - m.Boards[1].Store
	if netScore > 0 {
		return 0, netScore
	} else {
		return 1, -netScore
	}
}

func mancalaResult(flag int, seq []int, size int) int {
	m := NewMancalaGame(flag)
	for i := 0; i < size; i++ {
		// 1. Detect if the seq is following the rules
		m.print()
		err := m.playOneStep(seq[i])
		if err != nil {
			println(err)
			return 30000 + i
		}
		// 2. Check if the game ends
		if m.checkEnd() {
			if i == size-1 {
				return 15000 + m.Boards[0].Store - m.Boards[1].Store
			} else {
				return 30000 + i + 1
			}
		}
	}
	return 20000 + m.Boards[0].Store
}
