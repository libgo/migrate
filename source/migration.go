package source

import (
	"fmt"
)

type Module string

var (
	ErrNoModule error = fmt.Errorf("no this module")
	ErrTop      error = fmt.Errorf("beyond the up limit")
	ErrBottom   error = fmt.Errorf("beyond the bottom limit")
)

type Migration struct {
	MaxVer  int
	Version int
	Sql     map[int]string
}

func (mig *Migration) Goto(v int) error {
	if v > mig.MaxVer {
		return ErrTop
	}
	if v < 0 {
		return ErrBottom
	}

	mig.Version = v
	return nil
}

func (mig *Migration) SetMax(v int) {
	if v > mig.MaxVer {
		mig.MaxVer = v
	}
}

func (mig *Migration) Ver() int {
	return mig.Version
}

func (mig *Migration) Next() (string, int, error) {
	if mig.Version >= mig.MaxVer {
		return "", 0, ErrTop
	}
	mig.Version++
	return mig.Sql[mig.Version], mig.Version, nil
}

//
//
type Migrations map[Module]*Migration

func (migs Migrations) Goto(m Module, v int) error {
	if _, ok := migs[m]; !ok {
		return ErrNoModule
	}

	return migs[m].Goto(v)
}

func (migs Migrations) Ver(m Module) (int, error) {
	if _, ok := migs[m]; !ok {
		return 0, ErrNoModule
	}

	return migs[m].Ver(), nil
}

func (migs Migrations) Next(m Module) (string, int, error) {
	if _, ok := migs[m]; !ok {
		return "", 0, ErrNoModule
	}

	return migs[m].Next()
}

func (migs Migrations) List() []Module {
	ms := []Module{}
	for m, _ := range migs {
		ms = append(ms, m)
	}

	return ms
}
