package nodejs

import (
	"os"
	"strconv"
	"strings"
	"syscall"

	"aiexec-sandbox/internal/core/lib"
	"aiexec-sandbox/internal/static/nodejs_syscall"
)

//var allow_syscalls = []int{}

func InitSeccomp(uid int, gid int, enable_network bool) error {
	err := syscall.Chroot(".")
	if err != nil {
		return err
	}
	err = syscall.Chdir("/")
	if err != nil {
		return err
	}

	lib.SetNoNewPrivs()

	allowed_syscalls := []int{}
	allowed_not_kill_syscalls := []int{}

	allowed_syscall := os.Getenv("ALLOWED_SYSCALLS")
	if allowed_syscall != "" {
		nums := strings.Split(allowed_syscall, ",")
		for num := range nums {
			syscall, err := strconv.Atoi(nums[num])
			if err != nil {
				continue
			}
			allowed_syscalls = append(allowed_syscalls, syscall)
		}
	} else {
		allowed_syscalls = append(allowed_syscalls, nodejs_syscall.ALLOW_SYSCALLS...)
		allowed_not_kill_syscalls = append(allowed_not_kill_syscalls, nodejs_syscall.ALLOW_ERROR_SYSCALLS...)

		if enable_network {
			allowed_syscalls = append(allowed_syscalls, nodejs_syscall.ALLOW_NETWORK_SYSCALLS...)
		}
	}

	err = lib.Seccomp(allowed_syscalls, allowed_not_kill_syscalls)
	if err != nil {
		return err
	}

	// setuid
	err = syscall.Setuid(uid)
	if err != nil {
		return err
	}

	// setgid
	err = syscall.Setgid(gid)
	if err != nil {
		return err
	}

	return nil
}
