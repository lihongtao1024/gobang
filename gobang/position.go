package gobang

import (
	"gobang/component"
	"sort"
)

var boardOffset = []*struct {
	x int8
	y int8
}{
	{-1, 1},
	{0, 1},
	{1, 1},
	{1, 0},
	{1, -1},
	{0, -1},
	{-1, -1},
	{-1, 0},
}

type position struct {
	x int8  `desc:"column"`
	y int8  `desc:"row"`
	s int32 `desc:"score"`
}

type positionHistory struct {
	addedPositions  []position
	removedPosition position
}

type possiblePositions struct {
	allHistory      []*positionHistory
	possibleMapping map[uint32]position
}

func newPossiblePositions() *possiblePositions {
	return &possiblePositions{
		allHistory:      make([]*positionHistory, 0),
		possibleMapping: make(map[uint32]position),
	}
}

func makePossibleKey(x, y int8) uint32 {
	return uint32(x)<<16 | uint32(y)
}

func (possibles *possiblePositions) addPriority(x, y int8) bool {
	key := makePossibleKey(x, y)
	if _, ok := possibles.possibleMapping[key]; ok {
		return false
	}

	possibles.possibleMapping[key] = position{x, y, 0}
	return true
}

func (possibles *possiblePositions) delPriority(x, y int8) bool {
	key := makePossibleKey(x, y)
	if _, ok := possibles.possibleMapping[key]; !ok {
		return false
	}

	delete(possibles.possibleMapping, key)
	return true
}

func (possibles *possiblePositions) addPossible(gobang *gobangImpl, x, y int8) {
	addeds := make([]position, 0)

	for _, offset := range boardOffset {
		offsetx, offsety := x+offset.x, y+offset.y
		if !gobang.isInBoard(offsetx, offsety) {
			continue
		}

		if gobang.boardPannel[offsety][offsetx] != component.GobangNil {
			continue
		}

		if possibles.addPriority(offsetx, offsety) {
			addeds = append(addeds, position{offsetx, offsety, 0})
		}
	}

	history := &positionHistory{}
	history.addedPositions = addeds

	if possibles.delPriority(x, y) {
		history.removedPosition.x = x
		history.removedPosition.y = y
	} else {
		history.removedPosition.x = -1
		history.removedPosition.y = -1
	}

	possibles.allHistory = append(possibles.allHistory, history)
}

func (possibles *possiblePositions) rollbackPossible() {
	l := len(possibles.allHistory)
	if l == 0 {
		return
	}

	history := possibles.allHistory[l-1]
	possibles.allHistory = possibles.allHistory[:l-1]

	for i := 0; i < len(history.addedPositions); i++ {
		pos := &history.addedPositions[i]
		possibles.delPriority(pos.x, pos.y)
	}

	pos := &history.removedPosition
	if pos.x != -1 && pos.y != -1 {
		possibles.addPriority(pos.x, pos.y)
	}
}

func (possibles *possiblePositions) getPossible(evaluation func(*position),
) []position {
	list := make([]position, 0, len(possibles.possibleMapping))
	for _, pos := range possibles.possibleMapping {
		evaluation(&pos)
		list = append(list, pos)
	}

	sort.Slice(list, func(i, j int) bool {
		pi, pj := list[i], list[j]
		if pi.s != pj.s {
			return pi.s > pj.s
		}

		if pi.x != pj.x {
			return pi.x < pj.x
		}

		return pi.y < pj.y
	})

	return list
}
