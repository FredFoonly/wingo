package main

import (
	//"bufio"
	//"fmt"
	//"io"
	//"net"
	//"os"
	//"path"
	"strings"

	//"github.com/BurntSushi/xgbutil"

	//"github.com/FredFoonly/wingo/commands"
	//"github.com/FredFoonly/wingo/logger"
)

func httpAddress() string {
	// Take the http addr from the command line if possible
	if len(flagHttpAddr) > 0 {
		return strings.TrimSpace(flagHttpAddr)
	}

	// We weren't handed a path on a plate, so have to synthesize it as best we can
	return ":8080"
}
