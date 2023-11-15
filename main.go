package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Position struct {
	row, col int
}

type TicTacToe struct {
	movesPlayed int
	winState int
	board [3][3]int
}

func NewTicTacToe() TicTacToe {
	t := TicTacToe{}
	t.movesPlayed = 0
	t.winState = 0
	return t
}

func (t *TicTacToe) GetBoard() [3][3]int {
	return t.board
}

func (t *TicTacToe) PrintBoard() {
	m := map[int]string{
		0: ".",
		1: "X",
		2: "O",
	}
	for i := 0; i < 3; i++ {
		fmt.Fprint(os.Stderr, "[")
		for j := 0; j < 3; j++ {
			fmt.Fprint(os.Stderr, m[t.board[i][j]])
		}
		fmt.Fprintln(os.Stderr, "]")
	}
}

func (t *TicTacToe) GetAvailableMoves() []Position {
	availableMoves := []Position{}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if t.board[i][j] == 0 {
				availableMoves = append(availableMoves, Position{row: i, col: j})
			}
		}
	}
	return availableMoves
}

func (t *TicTacToe) SetBoard(board [3][3]int) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			t.board[i][j] = board[i][j]
			if board[i][j] > 0 {
				t.movesPlayed += 1
			}
		}
	}
}

func (t *TicTacToe) PlaceMove(row int, col int, player int) error {
	if row == -1 && col == -1 {
		return nil
	}
	if t.board[row][col] != 0 {
		return errors.New("Invalid move")
	}
	t.board[row][col] = player
	t.movesPlayed += 1
	t.winState = t.CheckWin()
	return nil
}

func (t *TicTacToe) CheckWin() int {
	if t.movesPlayed == 9 {
		return -1
	}
	for i := 0; i < 3; i++ {
		if t.board[i][0] > 0 && t.board[i][0] == t.board[i][1] && t.board[i][1] == t.board[i][2] {
			return t.board[i][0]
		}
		if t.board[0][i] > 0 && t.board[0][i] == t.board[1][i] && t.board[1][i] == t.board[2][i] {
			return t.board[0][i]
		}
	}
	if (t.board[0][0] > 0 && t.board[0][0] == t.board[1][1] && t.board[1][1] == t.board[2][2]) ||
		(t.board[0][2] > 0 && t.board[0][2] == t.board[1][1] && t.board[1][1] == t.board[2][0]) {
		return t.board[1][1]
	}
	return 0
}

type UltimateTicTacToe struct {
	boards [3][3]TicTacToe
	activeBoard Position
	isBoardOpen bool
	playerTurn int
	movesPlayed int
}

func NewUltimateTicTacToe() UltimateTicTacToe {
	ut := UltimateTicTacToe{}
	ut.playerTurn = 1
	ut.movesPlayed = 0
	ut.isBoardOpen = true
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			ut.boards[i][j] = NewTicTacToe()
		}
	}
	return ut
}

func (ut *UltimateTicTacToe) PrintBoards() {
	m := [3]string{
		0: ".",
		1: "X",
		2: "O",
	}
	if ut.isBoardOpen {
		m[0] = "-"
	}
	for r := 0; r < 9; r++ {
		fmt.Fprintf(os.Stderr,"[")
		if r % 3 == 0 {
			fmt.Fprintf(os.Stderr, "\n")
		}
		for c := 0; c < 9; c++ {
			subBoardPos := Position {
				row: int(r/3),
				col: int(c/3),
			}
			subRow := r % 3
			subCol := c % 3

			if c % 3 == 0 {
				fmt.Fprintf(os.Stderr, "\t")
			}
			mark := m[ut.boards[subBoardPos.row][subBoardPos.col].board[subRow][subCol]]
			if ut.activeBoard.row == subBoardPos.row && ut.activeBoard.col == subBoardPos.col && mark == "." {
				mark = "-"
			}
			fmt.Fprintf(os.Stderr,"%s",mark)

		}
		fmt.Fprintf(os.Stderr,"\n")
	}
	winStates, _, _ := ut.getWinStates()
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			fmt.Fprintf(os.Stderr, "%v", winStates[i][j])
		}
		fmt.Fprintf(os.Stderr, "\n")
	}
	fmt.Fprintf(os.Stderr,"\n")
}

func (ut *UltimateTicTacToe) PlaceMove(row int, col int) {
	if row == -1 && col == -1 {
		return
	}
	subBoardPos := Position{
		row: int(row / 3),
		col: int(col / 3),
	}
	subRow := row % 3
	subCol := col % 3

	subBoard := &ut.boards[subBoardPos.row][subBoardPos.col]
	err := subBoard.PlaceMove(subRow, subCol, ut.playerTurn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error placing move (%d, %d) on board (%d, %d): %v\n", subRow, subCol, subBoardPos.row, subBoardPos.col, err)
	}
	ut.movesPlayed +=1
	if ut.playerTurn == 1 {
		ut.playerTurn = 2
	} else {
		ut.playerTurn = 1
	}

	ut.activeBoard = Position {row: subRow, col: subCol}
	if (&ut.boards[subRow][subCol]).winState == 0 {
		ut.isBoardOpen = false
	} else {
		ut.isBoardOpen = true
	}
}

func (ut *UltimateTicTacToe) getWinStates() ([3][3]int, int, int) {
	winStates := [3][3]int{}
	player1wins := 0
	player2wins := 0
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			state := ut.boards[r][c].winState
			winStates[r][c] = state
			if state == 1 {
				player1wins += 1
			} else if state == 2 {
				player2wins += 1
			}
		}
	}
	return winStates, player1wins, player2wins
}

func (ut *UltimateTicTacToe) CheckWin() int {
	winStates, player1wins, player2wins := ut.getWinStates()

	for i := 0; i < 3; i++ {
		if winStates[i][0] > 0 && (winStates[i][0] == winStates[i][1] && winStates[i][1] == winStates[i][2]) {
			return winStates[i][0]
		}
		if winStates[0][i] > 0 && (winStates[0][i] == winStates[1][i] && winStates[1][i] == winStates[2][i]) {
			return winStates[0][i]
		}
	}
	if (winStates[0][0] > 0 && winStates[0][0] == winStates[1][1] && winStates[1][1] == winStates[2][2]) ||
		(winStates[0][2] > 0 &&
			winStates[0][2] == winStates[1][1] &&
			winStates[1][1] == winStates[2][0]) {
		return winStates[1][1]
	}

	if ut.movesPlayed == 81 {
		if player1wins > player2wins {
			return 1
		} else if player2wins > player1wins {
			return 2
		} else {
			return -1
		}
	}
	return 0
}

func (ut *UltimateTicTacToe) GetAvailableMoves() []Position {
	if !ut.isBoardOpen {
		moves := ut.boards[ut.activeBoard.row][ut.activeBoard.col].GetAvailableMoves()
		for i := 0; i < len(moves); i++ {
			moves[i].row += ut.activeBoard.row * 3
			moves[i].col += ut.activeBoard.col * 3
		}
		return moves
	}

	var totalMoves []Position
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			if ut.boards[r][c].winState != 0 {
				continue
			}
			subMoves := ut.boards[r][c].GetAvailableMoves()
			s := 3
			for i := 0; i < len(subMoves); i++ {
				subMoves[i].row += r * s
				subMoves[i].col += c * s
			}
			totalMoves = append(totalMoves, subMoves...)
		}
	}
	return totalMoves
}

func (ut *UltimateTicTacToe) Copy() UltimateTicTacToe {
	nut := NewUltimateTicTacToe()
	nut.activeBoard = ut.activeBoard
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			nut.boards[i][j].SetBoard(ut.boards[i][j].GetBoard())
		}
	}
	nut.isBoardOpen = ut.isBoardOpen
	nut.movesPlayed = ut.movesPlayed
	nut.playerTurn = ut.playerTurn
	nut.activeBoard.row = ut.activeBoard.row
	nut.activeBoard.col = ut.activeBoard.col
	return nut
}

func main() {
	ut := NewUltimateTicTacToe()
	for {
		var opponentRow, opponentCol int
		fmt.Scan(&opponentRow, &opponentCol)
		ut.PlaceMove(opponentRow, opponentCol)

		var validActionCount int
		fmt.Scan(&validActionCount)
		for i := 0; i < validActionCount; i++ {
			var row, col int
			fmt.Scan(&row, &col)
		}

		////////////////////////////////////////
		availableMoves := ut.GetAvailableMoves()
		l := availableMoves
		rand.Shuffle(len(l), func(i, j int) { l[i], l[j] = l[j], l[i] })

		winnerMoveScore := make([]int, len(availableMoves))
		start := time.Now()
		nMoves := len(availableMoves)
		for j := 0; j < 1000000; j++ {
			if time.Since(start).Milliseconds() > 110 {
				break
			}
			i := j % nMoves
			copyUt := ut
			r := l[i].row
			c := l[i].col
			targetWinner := copyUt.playerTurn
			var targetLoser int
			if targetWinner == 1 {
				targetLoser = 2
			} else {
				targetLoser = 1
			}
			copyUt.PlaceMove(r, c)
			for {
				if copyUt.CheckWin() > 0 {
					break
				}
				testAvail := copyUt.GetAvailableMoves()
				if len(testAvail) ==  0 {
					break
				}
				move := testAvail[rand.Intn(len(testAvail))]
				copyUt.PlaceMove(move.row, move.col)

			}
			winner := copyUt.CheckWin()
			if winner == targetWinner {
				winnerMoveScore[i] += 1
			} else if winner == targetLoser {
				winnerMoveScore[i] -= 1
			}
		}

		var nextWinnerMove = Position{row:-1, col:-1}
		winnerMaxScore := -1000000
		for i := 0; i < len(winnerMoveScore); i++ {
			if winnerMoveScore[i] >= winnerMaxScore {
				winnerMaxScore = winnerMoveScore[i]
				nextWinnerMove = availableMoves[i]
			}
		}
		nextMove := nextWinnerMove
		ut.PlaceMove(nextMove.row, nextMove.col)

		fmt.Printf("%d %d\n", nextMove.row, nextMove.col)
	}
}