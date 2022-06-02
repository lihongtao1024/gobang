package component

type GobangColor int
type GobagWinner int

const (
	GobagWinnerNil = GobagWinner(iota)
	GobagWinnerPlayer
	GobagWinnerComputer
)

const (
	GobangNil = GobangColor(iota)
	GobangWhite
	GobangBlack
)

type Gobang interface {
	GetGameWinner() GobagWinner
	Move(x, y int8) (rx, ry int8)
	Draw()
}
