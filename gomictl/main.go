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

	fmt.Println(args)

	host := args["--host"].(string)
	db := args["<db>"].(string)

	if args["init"].(bool) {
		r, err := gomi.NewRepo(host, db)
		if err != nil {
			panic(err)
		}
		r.Init()
		return
	}

	//structure := args["<structure>"].(string)
	if args["structure"].(bool) {
		return
	}
	if args["migrate"].(bool) {
		return
	}
}
