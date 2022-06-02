package ac

import "gobang/algorithm"

type acPattern struct {
	acStates []*acPatternState
}

type acPatternState struct {
	curKey       rune
	parentIndex  int
	selfIndex    int
	failIndex    int
	childIndexs  map[rune]int
	outputIndexs []int
}

type acSearcher struct {
	*acPattern
}

func NewPattern() algorithm.ACPattern {
	pattern := &acPattern{make([]*acPatternState, 1)}
	pattern.acStates[0] = newPatternState(-1, 0, 0)
	return pattern
}

func newPatternState(parent int, c rune, i int) *acPatternState {
	return &acPatternState{
		curKey:       c,
		parentIndex:  parent,
		selfIndex:    i,
		failIndex:    -1,
		childIndexs:  make(map[rune]int),
		outputIndexs: make([]int, 0),
	}
}

func (p *acPattern) addChildState(parent int, c rune, i int) int {
	if parent == -1 || i == -1 {
		panic("system error1.")
	}

	ps := p.acStates[parent]
	if _, ok := ps.childIndexs[c]; ok {
		panic("system error2.")
	}

	child := newPatternState(parent, c, i)
	p.acStates = append(p.acStates, child)

	ps.childIndexs[c] = i
	return i
}

func (p *acPattern) findChildState(parent int, c rune) int {
	if parent == -1 {
		panic("system error3.")
	}

	i, ok := p.acStates[parent].childIndexs[c]
	if !ok {
		return -1
	}

	return i
}

func (p *acPattern) AddPattern(id int, s string) {
	cur := 0
	for _, c := range s {
		next := p.findChildState(cur, c)
		if next == -1 {
			next = p.addChildState(cur, c, len(p.acStates))
		}

		cur = next
	}

	p.acStates[cur].outputIndexs = append(p.acStates[cur].outputIndexs, id)
}

func (p *acPattern) Done() {
	states := make([]int, 0)
	for _, child1 := range p.acStates[0].childIndexs {
		cs := p.acStates[child1]
		cs.failIndex = 0

		for _, child2 := range cs.childIndexs {
			states = append(states, child2)
		}
	}

	for len(states) > 0 {
		childstates := make([]int, 0)

		for _, cur := range states {
			curstate := p.acStates[cur]
			for _, child := range curstate.childIndexs {
				childstates = append(childstates, child)
			}

			curfail := p.acStates[curstate.parentIndex].failIndex
			for {
				fail := p.findChildState(curfail, curstate.curKey)
				if fail == -1 {
					curfail = p.acStates[curfail].failIndex
				} else {
					curstate.failIndex = fail
					if len(curstate.outputIndexs) > 0 {
						curfailstate := p.acStates[fail]
						curstate.outputIndexs = append(
							curstate.outputIndexs,
							curfailstate.outputIndexs...,
						)
					}
					break
				}

				if curfail == -1 {
					curstate.failIndex = 0
					break
				}
			}
		}

		states = childstates
	}
}

func NewSearcher(pattern algorithm.ACPattern) algorithm.ACSearcher {
	return &acSearcher{pattern.(*acPattern)}
}

func (s *acSearcher) Match(ss string) []int {
	cur := 0
	rs := []rune(ss)
	output := make([]int, 0)

	for i := 0; i < len(rs); {
		c := rs[i]
		next := s.findChildState(cur, c)
		if next == -1 {
			for {
				next := s.acStates[cur].failIndex
				if next == -1 {
					cur = 0
					i++
					break
				}

				cur = next
				next = s.findChildState(cur, c)
				if next == -1 {
					continue
				}

				cur = next
				i++
				break
			}
		} else {
			cur = next
			i++
		}

		curstate := s.acStates[cur]
		if len(curstate.outputIndexs) > 0 {
			output = append(output, curstate.outputIndexs...)
		}

		if curstate.failIndex != -1 {
			curfailstate := s.acStates[curstate.failIndex]
			if len(curfailstate.childIndexs) > 0 {
				output = append(output, curfailstate.outputIndexs...)
			}
		}
	}

	return output
}
