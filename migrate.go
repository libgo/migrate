package main

import (
	"flag"

	"github.com/libgo/logx"
	"github.com/libgo/migrate/database"
	_ "github.com/libgo/migrate/database/mysql"
	"github.com/libgo/migrate/internal"
	"github.com/libgo/migrate/source"
	_ "github.com/libgo/migrate/source/file"
)

var (
	d  string
	p  string
	m  string
	dg bool
)

func init() {
	// Register dsn format -> [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	flag.StringVar(&d, "d", "root:ddg1208@tcp(192.168.10.191:3306)/dolphin", "database uri")
	// Path to migrate files
	flag.StringVar(&p, "p", "./migrate", "migration source file path")
	// Module to up, default is "all"
	flag.StringVar(&m, "m", "all", "module to up")
	// Debug flag
	flag.BoolVar(&dg, "D", false, "dubug")
}

func main() {
	flag.Parse()

	logx.SetGlobalLevel(logx.InfoLevel)
	if dg {
		logx.SetGlobalLevel(logx.DebugLevel)
	}

	db, err := database.Open("mysql", d)
	if err != nil {
		logx.Errorf("Open database error: %s", err.Error())
		return
	}

	r, err := source.Open("file", p)
	if err != nil {
		logx.Errorf("Read source migration error: %s", err.Error())
		return
	}

	mig := internal.New(r, db)
	if m == "all" {
		err = mig.UpAll()
	} else {
		err = mig.Up(source.Module(m))
	}

	if err != nil {
		logx.Errorf("Migrating error: %s", err.Error())
	}
}
