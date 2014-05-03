package main

import (
    "fmt"
    "github.com/fmd/gomi/gomi"
    "github.com/docopt/docopt-go"
)

func usage() string {
    return `gomictl.

Usage:
    gomictl init <db>
    gomictl structure <db> <structure>
    gomictl migrate <db> <structure>

Options:
    -h --host   MongoDB host string [default: "127.0.0.1"]. 
    -h --help   Show this screen.
    --version   Show version.`
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