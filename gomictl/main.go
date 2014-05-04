package gomictl

import (
	"github.com/docopt/docopt-go"
	"github.com/fmd/gomi/gomi"
)

func usage() string {
	return `gomictl.

Usage:
    gomictl init <db>
    gomictl structure <db> <structure>
    gomictl migrate <db>

Options:
    -h --host   MongoDB host string [default: "127.0.0.1"].
    -h --help   Show this screen.
    --version   Show version.`
}

func init(db string) {
	m := gomi.NewMigrator()
}

func main() {
	args, _ := docopt.Parse(usage(), nil, true, "gomictl v0.1.0", false)
	db := args["<db>"].(string)
	if args["init"].(bool) {
		return
	}

	structure := args["<structure>"].(string)
	if args["structure"].(bool) {
		return
	}
	if args["migrate"].(bool) {
		return
	}
}
