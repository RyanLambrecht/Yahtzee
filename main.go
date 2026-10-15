package main

import (
	"flag"
	"yahtzee/ui/cli"
	"yahtzee/ui/web"
)

func main() {
	cmdMode := flag.Bool("c", false, "start in command line mode")
	flag.Parse()
	if *cmdMode {
		cli.Run()
		return
	}

	web.Run()
}
