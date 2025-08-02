package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "child" {
		runChild()
		return
	}

	cmd := exec.Command("/proc/self/exe", "child")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWIPC,
	}

	if err := cmd.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func runChild() {

	if err := syscall.Sethostname([]byte("mycontainer")); err != nil {
		log.Fatalf("failed to set hostname %v", err)
	}

	if err := syscall.Chroot("./rootfs"); err != nil {
		log.Fatalf("Failed to chroot: %v", err)
	}
	if err := os.Chdir("/"); err != nil {
		log.Fatalf("Failed to change dir: %v", err)
	}
	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		log.Fatalf("Failed to mount /proc: %v", err)
	}
	env := append(os.Environ(), "PS1=\\u@\\h:\\w\\$ ")
	syscall.Exec("/bin/sh", []string{"/bin/sh"}, env)

	if err := syscall.Exec("/bin/sh", []string{"/bin/sh"}, env); err != nil {
		log.Fatalf("Failed to exec shell: %v", err)
	}

}
