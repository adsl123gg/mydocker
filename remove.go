package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"mydocker/container"
)

func RemoveContainer(containerName string) {
	containerInfo, err := getContainerInfoByName(containerName)
	if err != nil {
                log.Errorf("getContainerInfoByName in RemoveContainer %s error: %v", getContainerInfoByName, err)
                return
        }

	if containerInfo.Status == container.STOP {
		dirUrl := fmt.Sprintf(container.DefaultInfoLocation, containerName)
		if err := os.RemoveAll(dirUrl); err != nil {
			log.Errorf("remove %s in RemoveContainer error: %v", dirUrl, err)
		}
	} else {
		log.Errorf("can't remove running container %s", containerName)
	}

}

