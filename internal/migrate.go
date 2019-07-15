package internal

import (
	"fmt"

	"github.com/libgo/logx"
	"github.com/libgo/migrate/database"
	"github.com/libgo/migrate/source"
)

var (
	ErrDirty = fmt.Errorf("migration dirty")
)

type Migrate struct {
	source   source.Reader
	database database.Driver
}

func New(s source.Reader, d database.Driver) *Migrate {
	return &Migrate{source: s, database: d}
}

func (m *Migrate) Up(md source.Module) error {
	logx.Infof("Migrating module: %s", string(md))
	if err := m.database.Lock(md); err != nil {
		return err
	}
	defer m.database.Unlock(md)

	v, d, err := m.database.Version(md)
	if err != nil {
		return err
	}

	if d {
		return fmt.Errorf("module %s:%d is dirty, should clean it by yourself", string(md), v)
	}

	err = m.source.Goto(md, v)
	if err != nil {
		return fmt.Errorf("database version is beyond migration range, something panic")
	}

	for {
		sql, nv, err := m.source.Next(md)
		if err != nil && err == source.ErrTop {
			break
		}

		logx.Debugf("Migrating %s[%d]", string(md), nv)
		err = m.database.Exec(sql)

		if err != nil {
			m.database.SetVer(md, nv, true)
			return fmt.Errorf("module %s:%d exec failed, should clean it by yourself, err is %s", string(md), nv, err.Error())
		}

		m.database.SetVer(md, nv, false)
	}

	return nil
}

func (m *Migrate) UpAll() error {
	ms := m.source.List()

	for _, md := range ms {
		err := m.Up(md)
		if err != nil {
			return err
		}
	}
	return nil
}
