package main

import (
	"errors"
	"syscall/js"
)

func main() {
	wait := make(chan struct{}, 0)
	js.Global().Set("mancalaOperator", js.FuncOf(MancalaOperator))
	<-wait
}

func MancalaOperator(this js.Value, args []js.Value) interface{} {
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
	// println("flag: ", flag, " status: ", status)
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

func (m *mancalaGame) makeMove(fromHoleIndex int) (mancalaGame, error) {
	clonedGame := *m
	// println("It's ", clonedGame.WhoseTurn, "'s turn.")
	player, err := clonedGame.sow(clonedGame.WhoseTurn, fromHoleIndex)
	if err != nil {
		return clonedGame, err
	}
	clonedGame.WhoseTurn = player
	// println("It'll be player ", clonedGame.WhoseTurn, "'s turn.")
	return clonedGame, err
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

func (it *moveIterator) Next() (mancalaGame, int, bool) {
	for i := it.index; i < 6; i = it.index {
		// println("iteration", i, "pieces", it.game.Boards[it.player].Holes[i])
		it.index++
		if it.game.Boards[it.player].Holes[i] > 0 {
			game, err := it.game.makeMove(i)
			if err == nil {
				// println("Returning")
				return game, i, true
			}
		}
	}
	return mancalaGame{}, -1, false
}

func (m *mancalaGame) evaluate(player int) int {
	score := m.Boards[player].Store - m.Boards[1-player].Store
	return score
}

func mancalaOperator(flag int, status []int) int {
	player := flag - 1
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
		WhoseTurn: player,
		FirstHand: 1, // unnecessary
		IsEnd:     false,
	}

	// eval, nextStep := naiveMinimax(&game, player, 4)
	_, nextStep := αβMinimax(&game, player, 6, -9999, 9999)
	nextStep = flag*10 + nextStep + 1
	return nextStep
}

// The `player` parameter is always the same — from which player's perspective for decision, i.e. the maximizing player. We determine the current move during simulation by `WhoseTurn`.
// Returns the best score and the best move.
func naiveMinimax(m *mancalaGame, player int, depth int) (int, int) {
	// 终止条件: 如果到达指定的深度或游戏结束
	if depth == 1 || m.checkEnd() {
		// 返回当前状态的评估分数
		return m.evaluate(player), -1
	}

	maxEval := -9999
	// 遍历所有可能的移动
	it := moveIterator{m, player, 0}
	nextStep := 0
	// m.print()
	for game, i, ok := it.Next(); ok; game, i, ok = it.Next() {
		eval, _ := naiveMinimax(&game, player, depth-1)
		// println("$$ Depth:", depth, "Step: ", nextStep, "Eval", eval, "MaxEval: ", maxEval)
		if eval > maxEval {
			maxEval = eval
			nextStep = i
		}
	}
	if maxEval == -9999 {
		maxEval = m.evaluate(player)
	}
	return maxEval, nextStep
}

// The `player` parameter is always the same — from which player's perspective for decision, i.e. the maximizing player. We determine the current move during simulation by `WhoseTurn`.
// Returns the best score and the best move.
// https://oi-wiki.org/search/alpha-beta/
func αβMinimax(m *mancalaGame, player int, depth int, alpha int, beta int) (int, int) {
	if depth == 1 || m.checkEnd() {
		return m.evaluate(player), -1
	}
	nextStep := -1
	if m.WhoseTurn == player { // maximizing player
		maxEval := -9999
		it := moveIterator{m, player, 0}
		for game, i, ok := it.Next(); ok; game, i, ok = it.Next() {
			countLeft := 48 - (game.Boards[0].Store + game.Boards[1].Store)
			eval, _ := αβMinimax(&game, player, depth-1, alpha, beta)
			if game.WhoseTurn == player {
				eval += eval
			}
			if countLeft < 13 {
				eval += eval / 2
			}

			if eval > maxEval {
				maxEval = eval
				nextStep = i
			}
			alpha = max(alpha, eval)
			if beta <= alpha {
				break
			}
		}
		if maxEval == -9999 {
			maxEval = m.evaluate(player)
		}
		return maxEval, nextStep
	} else {
		minEval := 9999
		it := moveIterator{m, player, 0}
		for game, i, ok := it.Next(); ok; game, i, ok = it.Next() {
			countLeft := 48 - (game.Boards[0].Store + game.Boards[1].Store)
			eval, _ := αβMinimax(&game, player, depth-1, alpha, beta)
			if game.WhoseTurn == player {
				eval += eval / 2
			}
			if countLeft < 13 {
				eval += eval / 3
			}
			if eval < minEval {
				minEval = eval
				nextStep = i
			}
			beta = min(beta, eval)
			if beta <= alpha {
				break
			}
		}
		if minEval == 9999 {
			minEval = m.evaluate(player)
		}
		return minEval, nextStep
	}
}
