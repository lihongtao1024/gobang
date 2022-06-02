package gobang

import (
	"bytes"
	"fmt"
	"gobang/ac"
	"gobang/algorithm"
	"gobang/component"
)

const boardWidth = 15
const (
	int32Min = -1000000
	int32Max = 1000000
)

var directionOffset = [][]*struct {
	x int8
	y int8
}{
	{
		{-5, 0}, {-4, 0}, {-3, 0}, {-2, 0}, {-1, 0},
		{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 0},
	},
	{
		{0, -5}, {0, -4}, {0, -3}, {0, -2}, {0, -1},
		{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {0, 5},
	},
	{
		{5, -5}, {4, -4}, {3, -3}, {2, -2}, {1, -1},
		{0, 0}, {-1, 1}, {-2, 2}, {-3, 3}, {-4, 4}, {-5, 5},
	},
	{
		{-5, -5}, {-4, -4}, {-3, -3}, {-2, -2}, {-1, -1},
		{0, 0}, {1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5},
	},
}

var boardPattern = []*struct {
	p string
	s int32
}{
	{"11111", 50000},
	{"011110", 4320},
	{"011100", 720},
	{"001110", 720},
	{"011010", 720},
	{"010110", 720},
	{"11110", 720},
	{"01111", 720},
	{"11011", 720},
	{"10111", 720},
	{"11101", 720},
	{"001100", 120},
	{"001010", 120},
	{"010100", 120},
	{"000100", 20},
	{"001000", 20},
}

type gobangImpl struct {
	algorithm.ACPattern
	algorithm.ACSearcher
	*possiblePositions
	boardDepth  int
	boardScore  [2]int32
	boardResult position
	playerSide  component.GobangColor
	gameWinner  component.GobagWinner
	boardPannel [boardWidth][boardWidth]component.GobangColor
}

func NewBoard(depth int, side component.GobangColor) component.Gobang {
	gobang := &gobangImpl{
		boardDepth:        depth,
		boardResult:       position{-1, -1, 0},
		playerSide:        side,
		possiblePositions: newPossiblePositions(),
	}

	gobang.ACPattern = ac.NewPattern()
	for index, pattern := range boardPattern {
		gobang.AddPattern(index, pattern.p)
	}
	gobang.Done()

	gobang.ACSearcher = ac.NewSearcher(gobang.ACPattern)
	return gobang
}

func (gobang *gobangImpl) isInBoard(x, y int8) bool {
	if x < 0 || x >= boardWidth || y < 0 || y >= boardWidth {
		return false
	}

	return true
}

func (gobang *gobangImpl) GetGameWinner() component.GobagWinner {
	return gobang.gameWinner
}

func (gobang *gobangImpl) Move(x, y int8) (int8, int8) {
	fmt.Println("move:", x, ",", y)
	npos := &gobang.boardResult
	npos.x = -1
	npos.y = -1

	pos := &gobang.boardPannel[y][x]
	if *pos != component.GobangNil {
		panic("system error.")
	}

	color2 := component.GobangNil
	if gobang.playerSide == component.GobangWhite {
		color2 = component.GobangBlack
	} else if gobang.playerSide == component.GobangBlack {
		color2 = component.GobangWhite
	} else {
		panic("system error.")
	}

	*pos = gobang.playerSide

	gobang.updateScore()
	if gobang.boardScore[0] >= boardPattern[0].s {
		gobang.gameWinner = component.GobagWinnerPlayer
	} else if gobang.boardScore[1] >= boardPattern[0].s {
		gobang.gameWinner = component.GobagWinnerComputer
	}

	if gobang.gameWinner != component.GobagWinnerNil {
		return npos.x, npos.y
	}

	gobang.addPossible(gobang, x, y)

	gobang.moveNext()
	fmt.Println("result:", npos.x, ",", npos.y)
	gobang.boardPannel[npos.y][npos.x] = color2

	gobang.updateScore()
	if gobang.boardScore[0] >= boardPattern[0].s {
		gobang.gameWinner = component.GobagWinnerPlayer
	} else if gobang.boardScore[1] >= boardPattern[0].s {
		gobang.gameWinner = component.GobagWinnerComputer
	}

	if gobang.gameWinner != component.GobagWinnerNil {
		return npos.x, npos.y
	}

	gobang.addPossible(gobang, npos.x, npos.y)
	//gobang.Draw()
	return npos.x, npos.y
}

func (gobang *gobangImpl) moveNext() {
	c := component.GobangNil
	if gobang.playerSide == component.GobangWhite {
		c = component.GobangBlack
	} else if gobang.playerSide == component.GobangBlack {
		c = component.GobangWhite
	} else {
		panic("system error.")
	}

	gobang.alphaBeta(gobang.boardDepth, int32Min, int32Max, c)
}

func (gobang *gobangImpl) fillPattern(color component.GobangColor,
	pbuf, cbuf *bytes.Buffer) {
	switch color {
	case gobang.playerSide:
		pbuf.WriteByte('1')
		cbuf.WriteByte('2')
	case component.GobangNil:
		pbuf.WriteByte('0')
		cbuf.WriteByte('0')
	default:
		pbuf.WriteByte('2')
		cbuf.WriteByte('1')
	}
}

func (gobang *gobangImpl) evaluationPosition(pos *position) {
	ppat, cpat := [4]string{}, [4]string{}
	pbuf, cbuf := bytes.NewBuffer([]byte{}), bytes.NewBuffer([]byte{})

	for i := 0; i < 4; i++ {
		for j := 0; j < 11; j++ {
			offset := directionOffset[i][j]
			if offset.x == 0 && offset.y == 0 {
				pbuf.WriteByte('1')
				cbuf.WriteByte('1')
			} else {
				offsetx, offsety := pos.x+offset.x, pos.y+offset.y
				if !gobang.isInBoard(offsetx, offsety) {
					continue
				}

				gobang.fillPattern(
					gobang.boardPannel[offsety][offsetx],
					pbuf,
					cbuf,
				)
			}
		}

		ppat[i] = pbuf.String()
		cpat[i] = cbuf.String()
		pbuf.Reset()
		cbuf.Reset()
	}

	for i := 0; i < 4; i++ {
		results := gobang.Match(ppat[i])
		for _, r := range results {
			pos.s += boardPattern[r].s
		}

		results = gobang.Match(cpat[i])
		for _, r := range results {
			pos.s += boardPattern[r].s
		}
	}
}

func (gobang *gobangImpl) evaluationScore(color component.GobangColor) int32 {
	if color == gobang.playerSide {
		return gobang.boardScore[0]
	} else if color == component.GobangNil {
		return 0
	}

	return gobang.boardScore[1]
}

func (gobang *gobangImpl) updateScore() {
	ppat, cpat := [4][21]string{}, [4][21]string{}
	pbuf, cbuf := bytes.NewBuffer([]byte{}), bytes.NewBuffer([]byte{})

	for y := 0; y < 15; y++ {
		for x := 0; x < 15; x++ {
			color := gobang.boardPannel[y][x]
			gobang.fillPattern(
				color,
				pbuf,
				cbuf,
			)
		}

		ppat[0][y] = pbuf.String()
		cpat[0][y] = cbuf.String()
		pbuf.Reset()
		cbuf.Reset()
	}

	for x := 0; x < 15; x++ {
		for y := 0; y < 15; y++ {
			color := gobang.boardPannel[y][x]
			gobang.fillPattern(
				color,
				pbuf,
				cbuf,
			)
		}

		ppat[1][x] = pbuf.String()
		cpat[1][x] = cbuf.String()
		pbuf.Reset()
		cbuf.Reset()
	}

	index := 0
	for x := 4; x < 15; x++ {
		for y := 0; y <= x; y++ {
			offsetx, offfsety := x-y, y
			color := gobang.boardPannel[offfsety][offsetx]
			gobang.fillPattern(
				color,
				pbuf,
				cbuf,
			)
		}

		ppat[2][index] = pbuf.String()
		cpat[2][index] = cbuf.String()
		pbuf.Reset()
		cbuf.Reset()
		index++
	}

	for y := 0; y < 10; y++ {
		for x := 0; x < 14-y; x++ {
			offsetx, offfsety := 14-x, x+y+1
			color := gobang.boardPannel[offfsety][offsetx]
			gobang.fillPattern(
				color,
				pbuf,
				cbuf,
			)
		}

		ppat[2][index] = pbuf.String()
		cpat[2][index] = cbuf.String()
		pbuf.Reset()
		cbuf.Reset()
		index++
	}

	index = 0
	for y := 10; y >= 0; y-- {
		for x := 0; x < 15-y; x++ {
			offsetx, offfsety := x, y+x
			color := gobang.boardPannel[offfsety][offsetx]
			gobang.fillPattern(
				color,
				pbuf,
				cbuf,
			)
		}

		ppat[3][index] = pbuf.String()
		cpat[3][index] = cbuf.String()
		pbuf.Reset()
		cbuf.Reset()
		index++
	}

	for x := 0; x < 10; x++ {
		for y := 0; y < 14-x; y++ {
			offsetx, offfsety := x+y+1, y
			color := gobang.boardPannel[offfsety][offsetx]
			gobang.fillPattern(
				color,
				pbuf,
				cbuf,
			)
		}
		ppat[3][index] = pbuf.String()
		cpat[3][index] = cbuf.String()
		pbuf.Reset()
		cbuf.Reset()
		index++
	}

	gobang.boardScore[0] = 0
	gobang.boardScore[1] = 0

	for i := 0; i < 4; i++ {
		for j := 0; j < 21; j++ {
			results := gobang.Match(ppat[i][j])
			for _, r := range results {
				gobang.boardScore[0] += boardPattern[r].s
			}

			results = gobang.Match(cpat[i][j])
			for _, r := range results {
				gobang.boardScore[1] += boardPattern[r].s
			}
		}
	}

	//fmt.Println("update score:", gobang.boardScore)
}

func (gobang *gobangImpl) alphaBeta(depth int, alpha, beta int32,
	color component.GobangColor) int32 {
	color2 := component.GobangNil
	if color == component.GobangWhite {
		color2 = component.GobangBlack
	} else if color == component.GobangBlack {
		color2 = component.GobangWhite
	} else {
		panic("system error.")
	}

	if depth == 0 {
		//fmt.Println("final score:", gobang.evaluationScore(color)-gobang.evaluationScore(color2))
		return gobang.evaluationScore(color) - gobang.evaluationScore(color2)
	}

	count := 0
	positions := gobang.getPossible(gobang.evaluationPosition)

	for len(positions) > 0 {
		pos := positions[0]
		positions = positions[1:]

		gobang.boardPannel[pos.y][pos.x] = color
		gobang.addPossible(gobang, pos.x, pos.y)

		if depth == 1 {
			gobang.updateScore()
		}

		score := -gobang.alphaBeta(depth-1, -beta, -alpha, color2)
		//if depth == gobang.boardDepth {
		//	fmt.Println("depth:", depth, "alpha:", alpha, "beta:", beta, "score", score, "pos:", pos)
		//}

		gobang.rollbackPossible()
		gobang.boardPannel[pos.y][pos.x] = component.GobangNil

		if score >= beta {
			return beta
		}

		if score > alpha {
			alpha = score
			if depth == gobang.boardDepth {
				gobang.boardResult = pos
			}
		}

		count++
		if count >= 9 {
			break
		}
	}

	return alpha
}

func (gobang *gobangImpl) Reset() {
	for y := 0; y < boardWidth; y++ {
		for x := 0; x < boardWidth; x++ {
			gobang.boardPannel[y][x] = 0
		}
	}
}

func (gobang *gobangImpl) Draw() {
	fmt.Println("draw chess gobang:")
	for y := 0; y < boardWidth; y++ {
		for x := 0; x < boardWidth; x++ {
			switch gobang.boardPannel[y][x] {
			case component.GobangWhite:
				fmt.Printf("%c ", '1')
			case component.GobangBlack:
				fmt.Printf("%c ", '2')
			default:
				fmt.Printf("%c ", '0')
			}
		}
		fmt.Print("\n")
	}
}
