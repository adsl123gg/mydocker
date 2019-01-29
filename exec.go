package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"mydocker/container"
	"io/ioutil"
	"encoding/json"
	"strings"
	"os/exec"
	"os"
	_ "mydocker/nsenter"
)

var (
	ENV_EXEC_PID = "mydocker_pid"
	ENV_EXEC_CMD = "mydocker_cmd"
)

func ExecContainer(containerName string, cmdArray []string) {
	pid, err := getContainerPidByName(containerName)
	if err != nil {
		log.Errorf("get container %s Pid in getContainerPidByName error: %v", containerName, err)
	}

	cmdStr := strings.Join(cmdArray, " ")
	log.Infof("container Pid is %s, command is %s", pid, cmdStr)

	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	os.Setenv(ENV_EXEC_PID, pid)
	os.Setenv(ENV_EXEC_CMD, cmdStr)

	if err := cmd.Run(); err != nil {
		log.Errorf("exec container %s error: %v", containerName, err)
	}

}

func getContainerPidByName(containerName string) (string, error) {
	configDir := fmt.Sprintf(container.DefaultInfoLocation, containerName)
        configFile := configDir + container.ConfigName

        content, err := ioutil.ReadFile(configFile)
        if err != nil {
                return "", err
        }
        var containerInfo container.ContainerInfo
        if err := json.Unmarshal(content, &containerInfo); err != nil {
                return "", err
        }

        return containerInfo.Pid, nil
}

