package main

import (
	"errors"
	"fmt"
	"syscall/js"
)

func main() {
	wait := make(chan struct{}, 0)
	js.Global().Set("mancalaOperator", js.FuncOf(MancalaOperator))
	<-wait
}

func MancalaOperator(this js.Value, args []js.Value) interface{} {
	/*
	   - 函数名：`mancalaOperator()` 或者 `mancala_operator()` etc.，根据你选择的编程语言选择合适的命名格式；

	   - 参数：

	     - 一个 i32 类型数字 `flag`，为 `1` 或 `2`：代表应该为哪位选手思考行棋决策；

	     - 一个 i32 类型数组 `status`：

	     - 数组的元素个数为 14，每位数字的含义如下：

	       | 位 0 - 5                | 位 6                    | 位 7 - 12               | 位 13                   |
	       | ----------------------- | ----------------------- | ----------------------- | ----------------------- |
	       | 棋洞 11 - 16 中的棋子数 | 选手 1 计分洞中的棋子数 | 棋洞 21 - 26 中的棋子数 | 选手 2 计分洞中的棋子数 |

	   - 返回值：一个 i32 类型数字，代表**根据目前的棋盘状况，为了取得更大的净胜棋数，选手 `flag` 应当分配哪个棋洞中的棋子**。
	*/
	flag := args[0].Int()
	if flag != 1 && flag != 2 {
		return js.ValueOf("invalid flag")
	}
	status := make([]int, 14)
	for i := 0; i < 14; i++ {
		element := args[1].Index(i)
		if element.Type() != js.TypeNumber {
			// return js.ValueOf("Mancala operator has invalid status: status is not a 14-element arr")
		}
		status[i] = element.Int()
	}
	println("flag: ", flag, " status: ", status)
	return js.ValueOf(4)
	nextStep := mancalaOperator(flag, status)
	return js.ValueOf(nextStep)
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

func (m *mancalaGame) sow(player, fromHoleIndex int) (error, int) {
	// 1. Get the number of stones in the hole
	count := m.Boards[player].Holes[fromHoleIndex]
	if count == 0 {
		return errors.New("No stones in the hole"), -1
	}
	boardIndex := player
	holeIndex := fromHoleIndex + 1
	for count > 0 {
		board := &m.Boards[boardIndex]
		if holeIndex == 6 {
			if boardIndex == player { // 对于己方计分洞需要分配棋子
				board.Store++
				if count == 0 { // 如果正好落在己方记分洞中，则当前玩家再下一次
					return nil, boardIndex
				}
			} // 对于对方计分洞直接跳过
			boardIndex = 1 - boardIndex // 切换到对方
			holeIndex = 0
		} else {
			opponentBoard := &m.Boards[1-boardIndex]
			if board.Holes[holeIndex] == 0 && opponentBoard.Holes[holeIndex] > 0 {
				board.Store += 1 + opponentBoard.Holes[holeIndex]
				opponentBoard.Holes[holeIndex] = 0
			} else {
				board.Holes[holeIndex]++
			}
			holeIndex++
		}
		count--
	}
	return nil, 1 - player
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
	println("Playing from ", player, fromHoleIndex)
	err, nextPlayer := m.sow(player, fromHoleIndex)
	if err != nil {
		return err
	}
	m.WhoseTurn = nextPlayer
	return nil
}

func (m *mancalaGame) makeMove(player, fromHoleIndex int) (error, mancalaGame) {
	cloned := *m
	err, _ := cloned.sow(player, fromHoleIndex)
	return err, cloned
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

type moveIterator struct {
	game   *mancalaGame
	player int
	index  int
}

func (it *moveIterator) Next() (mancalaGame, bool) {
	for i := it.index; i < 6; {
		it.index++
		if it.game.Boards[it.player].Holes[i] > 0 {
			err, game := it.game.makeMove(it.player, i)
			if err == nil {
				return game, true
			}
		}
	}
	return mancalaGame{}, false
}

func (m *mancalaGame) evaluate(player int) int {
	return m.Boards[player].Store - m.Boards[1-player].Store
}

func (m *mancalaGame) getWinner() (int, int) {
	netScore := m.Boards[0].Store - m.Boards[1].Store
	if netScore > 0 {
		return 0, netScore
	} else {
		return 1, -netScore
	}
}

func mancalaOperator(flag int, status []int) int {
	game := mancalaGame{
		Boards: [2]mancalaBoard{
			{
				Holes: [6]int{status[0], status[1], status[2], status[3], status[4], status[5]},
				Store: status[6],
			},
			{
				Holes: [6]int{status[7], status[8], status[9], status[10], status[11], status[12]},
				Store: status[13],
			},
		},
		WhoseTurn: flag,
		FirstHand: 0,
		IsEnd:     false,
	}
	eval, nextStep := minimax(&game, flag, 5, -999, 999)
	println("Eval: ", eval, " Step: ", nextStep)
	return nextStep
}

func minimax(m *mancalaGame, player int, depth int, alpha int, beta int) (int, int) {
	// 终止条件: 如果到达指定的深度或游戏结束
	if depth == 0 || m.checkEnd() {
		// 返回当前状态的评估分数
		return m.evaluate(player), -1
	}

	maxEval := -9999
	// 遍历所有可能的移动
	it := moveIterator{m, player, 0}
	nextStep := 0
	for game, ok := it.Next(); ok; game, ok = it.Next() {
		eval, _ := minimax(&game, player, depth-1, alpha, beta)
		maxEval = max(maxEval, eval)
		if maxEval == eval {
			nextStep = it.index - 1
		}
		alpha = max(alpha, eval)
		if beta <= alpha {
			break // 剪枝
		}
	}
	return maxEval, nextStep
}
