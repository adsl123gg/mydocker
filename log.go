package main

import (
	"fmt"
	"os"
	"mydocker/container"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
)

func LogContainer(containerName string) {
	configDir := fmt.Sprintf(container.DefaultInfoLocation, containerName)
        logConfigFile := configDir + container.LogConfigName

        content, err := ioutil.ReadFile(logConfigFile)
        if err != nil {
                log.Errorf("get container %s log error: %v", containerName, err)
                return
        }
	fmt.Fprint(os.Stdout, string(content))

}

