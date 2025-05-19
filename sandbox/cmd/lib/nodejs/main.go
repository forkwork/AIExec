package main

import "aiexec-sandbox/internal/core/lib/nodejs"
import "C"

//export AiexecSeccomp
func AiexecSeccomp(uid int, gid int, enable_network bool) {
	nodejs.InitSeccomp(uid, gid, enable_network)
}

func main() {}
