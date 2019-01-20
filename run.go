package main

import (
	log "github.com/Sirupsen/logrus"
	"mydocker/container"
	"mydocker/subsystems"
	"os"
	"strings"
)

func Run(tty bool, cmdArray []string, res *subsystems.ResourceConfig) {
    parent, writePipe := container.NewParentProcess(tty)
    if parent == nil {
	log.Errorf("New parent process error")
	return
    }
    if err := parent.Start(); err != nil {
        log.Error(err)
    }

    cgroupMgr := subsystems.NewCgroupManager("mydocker-cgroup-hqc")
    defer cgroupMgr.Destory()

    cgroupMgr.Set(res)
    cgroupMgr.Apply(parent.Process.Pid)

    sendInitCommand(cmdArray, writePipe)

    parent.Wait()
    os.Exit(-1)
}

func sendInitCommand(cmdArray []string, writePipe *os.File) {
	command := strings.Join(cmdArray, " ")
	log.Infof("command all is %s", cmdArray)
	writePipe.WriteString(command)
	writePipe.Close()
}


