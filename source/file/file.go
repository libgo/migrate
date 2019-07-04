package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/libgo/logx"
	"github.com/libgo/migrate/source"
)

func init() {
	source.Register("file", &File{})
}

type File struct {
	migrations source.Migrations
}

var (
	re = regexp.MustCompile(`^([0-9]+)_([0-9a-zA-Z\.]*)_?(.*)?\.sql$`)
)

func (f *File) Open(uri string) (source.Reader, error) {
	migrations := make(source.Migrations)

	//scan directory
	err := filepath.Walk(uri, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		module, ver, sql, err := read(path)
		if err != nil {
			return nil
		}

		if _, ok := migrations[module]; !ok {
			migrations[module] = &source.Migration{MaxVer: 0, Version: 0, Sql: map[int]string{}}
		}

		migrations[module].Sql[ver] = sql
		migrations[module].SetMax(ver)

		return nil
	})

	if err != nil {
		return nil, err
	}

	rr := &File{
		migrations: migrations,
	}
	return rr, nil
}

func read(file string) (source.Module, int, string, error) {
	logx.Debugf("reading file: '%s'", file)
	_, f := filepath.Split(file)

	m := re.FindStringSubmatch(f)
	if len(m) == 0 {
		return "", 0, "", fmt.Errorf("not a valid migration file")
	}

	v, err := strconv.Atoi(m[1])
	if err != nil {
		return "", 0, "", fmt.Errorf("not a valid migration file")
	}

	fi, err := os.Open(file)
	if err != nil {
		return "", 0, "", err
	}
	defer fi.Close()

	sql, err := ioutil.ReadAll(fi)
	if err != nil {
		return "", 0, "", err
	}

	return source.Module(m[2]), v, string(sql), nil
}

func (f *File) Goto(m source.Module, v int) error {
	return f.migrations.Goto(m, v)
}

func (f *File) Next(m source.Module) (string, int, error) {
	return f.migrations.Next(m)
}

func (f *File) List() []source.Module {
	return f.migrations.List()
}
