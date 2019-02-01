package main

import (
	"fmt"
	"syscall"
	"strconv"
	"encoding/json"
	"mydocker/container"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
)

func StopContainer(containerName string) {
	pid, err := getContainerPidByName(containerName)
        if err != nil {
                log.Errorf("get container %s Pid in getContainerPidByName error: %v", containerName, err)
        }

	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		log.Errorf("convert pid to int error: %v", err)
		return
	}

	if err := syscall.Kill(pidInt, syscall.SIGTERM); err != nil {
                log.Errorf("kill container %s error: %v", containerName, err)
        }

	containerInfo, err := getContainerInfoByName(containerName)
	if err != nil {
                log.Errorf("getContainerInfoByName in StopContainer %s error: %v", getContainerInfoByName, err)
                return
        }

	containerInfo.Status = container.STOP
	containerInfo.Pid = ""
	newContentBytes, err := json.Marshal(containerInfo)
	if err != nil {
                log.Errorf("json marshal in StopContainer error: %v", err)
                return
        }
	dirUrl := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	fileName := dirUrl + "/" + container.ConfigName
	if err := ioutil.WriteFile(fileName, newContentBytes, 0622); err != nil {
		log.Errorf("write config file %s error: %v", containerName, err)
	}

}

func getContainerInfoByName(containerName string) (*container.ContainerInfo, error) {
        configDir := fmt.Sprintf(container.DefaultInfoLocation, containerName)
        configFile := configDir + container.ConfigName

        content, err := ioutil.ReadFile(configFile)
        if err != nil {
                log.Errorf("getContainerInfoByName %s error: %v", configFile, err)
                return nil, err
        }
        var containerInfo container.ContainerInfo
        if err := json.Unmarshal(content, &containerInfo); err != nil {
                log.Errorf("getContainerInfoByName unmarshal error %v", err)
                return nil, err
        }

        return &containerInfo, nil
}

