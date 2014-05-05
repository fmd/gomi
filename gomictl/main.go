package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/fmd/gomi/gomi"
)

func usage() string {
	return `gomictl.

Usage:
    gomictl init <db> [--host=<hostname>]
    gomictl structure <db> <structure> [--host=<hostname>]
    gomictl migrate <db> [--host=<hostname>]
    gomictl --help [--host=<hostname>]
    gomictl --version [--host=<hostname>]

Options:
    -h | --host   MongoDB host string [default: 127.0.0.1].
    --help      Show this screen.
    --version   Show version.`
}

func main() {
	args, _ := docopt.Parse(usage(), nil, true, "gomictl v0.1.0", false)

	host := args["--host"].(string)
	db := args["<db>"].(string)

	r, err := gomi.NewRepo(host, db)
	if err != nil {
		panic(err)
	}

	if args["init"].(bool) {
		err = r.Init()
		if err != nil {
			panic(err)
		}

		return
	}

	if args["migrate"].(bool) {
		migrations, err := gomi.LoadMigrations()
		for _, migration := range migrations {
			err = r.Migrator.Apply(migration)
			if err != nil {
				fmt.Println(err)
			}
		}

		return
	}

	structure := args["<structure>"].(string)
	if args["structure"].(bool) {
		s, err := gomi.LoadStructure(structure)
		if err != nil {
			panic(err)
		}

		g := r.Migrator.NewMigration(s)
		err = g.Save()

		if err != nil {
			panic(err)
		}

		return
	}
}
