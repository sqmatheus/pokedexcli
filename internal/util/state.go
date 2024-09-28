package util

import "github.com/sqmatheus/pokedexcli/internal/api"

type State[T any] struct {
	Next     string
	Previous string
}

func NewState[T any](next string, previous string) *State[T] {
	return &State[T]{Next: next, Previous: previous}
}

func (s *State[T]) Update(p api.Pagination[T]) {
	if p.Next != nil {
		s.Next = *p.Next
	} else {
		s.Next = ""
	}

	if p.Previous != nil {
		s.Previous = *p.Previous
	} else {
		s.Previous = ""
	}
}
