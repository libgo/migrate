package source

import (
	"fmt"
)

var readers = make(map[string]Reader)

type Reader interface {
	Open(string) (Reader, error)
	Goto(Module, int) error
	Next(Module) (string, int, error)
	List() []Module
}

func Open(scheme string, uri string) (Reader, error) {
	r, ok := readers[scheme]
	if !ok {
		return nil, fmt.Errorf("invalid reader")
	}

	return r.Open(uri)
}

func Register(name string, r Reader) {
	readers[name] = r
}
