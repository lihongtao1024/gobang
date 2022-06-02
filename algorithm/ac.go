package algorithm

type ACPattern interface {
	AddPattern(id int, s string)
	Done()
}

type ACSearcher interface {
	Match(ss string) []int
}
