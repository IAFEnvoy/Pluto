package util

import (
	"strings"
)

type NamedLock struct {
	Locked    map[string]struct{}
	Separator string
}

func NewNamedLock() *NamedLock {
	return &NamedLock{
		Locked:    make(map[string]struct{}),
		Separator: "/",
	}
}

func (t NamedLock) Lock(s ...string) {
	t.Locked[strings.Join(s, t.Separator)] = struct{}{}
}

func (t NamedLock) Unlock(s ...string) {
	delete(t.Locked, strings.Join(s, t.Separator))
}

func (t NamedLock) IsLocked(s ...string) bool {
	_, ok := t.Locked[strings.Join(s, t.Separator)]
	return ok
}
