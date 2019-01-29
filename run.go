package main

import (
	log "github.com/Sirupsen/logrus"
	"mydocker/container"
	//"mydocker/util"
	"mydocker/subsystems"
	"os"
	"strings"
	"fmt"
	"time"
	"strconv"
	"math/rand"
	"encoding/json"
)

func Run(tty bool, cmdArray []string, res *subsystems.ResourceConfig, volume string, imageName string, containerName string, envSlice []string) {
    id, containerName := getContainerIdName(containerName)
    parent, writePipe := container.NewParentProcess(tty, volume, imageName, containerName, envSlice)
    if parent == nil {
	log.Errorf("New parent process error")
	return
    }
    if err := parent.Start(); err != nil {
        log.Error(err)
    }

    //record container info
    containerName, err := recordContainerInfo(parent.Process.Pid, cmdArray, id, imageName, containerName)
    if err != nil {
    	log.Errorf("Record container info error %v", err)
    	return
    }

    cgroupMgr := subsystems.NewCgroupManager("mydocker-cgroup-hqc")
    defer cgroupMgr.Destory()

    cgroupMgr.Set(res)
    cgroupMgr.Apply(parent.Process.Pid)

    sendInitCommand(cmdArray, writePipe)

    if tty {
	parent.Wait()
	deleteContainerInfo(containerName)
	container.DeleteWorkSpace(containerName, volume)
    }

    //mntURL := "/root/image/mnt/"
    //rootURL :=  "/root/image/"
    //mntURL := util.GetMntURL()
    //rootURL :=  util.GetRootURL()
    os.Exit(0)
}

func deleteContainerInfo(containerName string) {
	dirUrl := fmt.Sprintf(container.DefaultInfoLocation, containerName)
        if err := os.RemoveAll(dirUrl); err != nil {
                log.Errorf("deleteContainerInfo Remove dir %s error %v", dirUrl, err)
        }
}

func sendInitCommand(cmdArray []string, writePipe *os.File) {
	command := strings.Join(cmdArray, " ")
	log.Infof("command all is %s", cmdArray)
	writePipe.WriteString(command)
	writePipe.Close()
}

func recordContainerInfo(containerPID int, commandArray []string, id string, imageName string, containerName string) (string, error) {
	command := strings.Join(commandArray, "")
	createTime := time.Now().Format("2019-01-28 18:06:05")

	containerInfo := &container.ContainerInfo{
		Id:          id,
		Pid:         strconv.Itoa(containerPID),
		Command:     command,
		CreatedTime: createTime,
		Status:      container.RUNNING,
		Name:        containerName,
		ImageName:   imageName,
	}

	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return "", err
	}
	jsonStr := string(jsonBytes)

	dirUrl := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	if !container.PathExists(dirUrl) {
		if err := os.MkdirAll(dirUrl, 0622); err != nil {
			log.Errorf("Mkdir error %s error %v", dirUrl, err)
			return "", err
		}
	}
	fileName := dirUrl + "/" + container.ConfigName
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		log.Errorf("Create file %s error %v", fileName, err)
		return "", err
	}
	if _, err := file.WriteString(jsonStr); err != nil {
		log.Errorf("File write string error %v", err)
		return "", err
	}

	return containerName, nil
}

func getContainerIdName(containerName string) (string, string) {
	id := randContainerId(10)
        if containerName == "" {
                containerName = id
        }
	return id, containerName
}

func randContainerId(n int) string {
	letterBytes := "1234567890"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
