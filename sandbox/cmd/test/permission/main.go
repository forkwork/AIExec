package main

import (
	"os/exec"

	"aiexec-sandbox/internal/core/lib/python"
)

func main() {
	python.InitSeccomp(0, 0, true)

	exec.Command("/bin/sh", "-c", "echo hello").Run()
}
