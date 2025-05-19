package main

import (
	"aiexec-sandbox/internal/core/lib/python"
)
import "C"

//export AiexecSeccomp
func AiexecSeccomp(uid int, gid int, enable_network bool) {
	python.InitSeccomp(uid, gid, enable_network)
}

func main() {}
