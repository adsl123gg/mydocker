package container

import (
	"os"
	"os/exec"
	"syscall"
	"strings"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"fmt"
)

func RunContainerInitProcess(command string, args []string) error {
    cmdArray := readUserCommand()
    if cmdArray == nil || len(cmdArray) == 0 {
	return fmt.Errorf("Run container get user command error, cmdArray is nil")
    }


    defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
    syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
    path, err := exec.LookPath(cmdArray[0])
    if err != nil {
	log.Errorf("exec.LookPath error %v", err)
        return nil
    }
    log.Infof("Find path is %s", path)
    if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
        log.Errorf(err.Error())
    }
    return nil
}

func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Errorf("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}

