package main

import (
	"mydocker/container"
	"fmt"
	"encoding/json"
	"os"
	"io/ioutil"
	"text/tabwriter"
	log "github.com/Sirupsen/logrus"
)

func ListContainer() {
	dirUrl := fmt.Sprintf(container.DefaultInfoLocation, "")
	dirUrl = dirUrl[:len(dirUrl) -1]

	files, err := ioutil.ReadDir(dirUrl)
	if err != nil {
		log.Errorf("read dir %s error: %v", dirUrl, err)
	}

	var containers []*container.ContainerInfo
	for _, file := range files {
		container, err := getContainerInfo(file)
		if err != nil {
			log.Errorf("get container %s info error: %v", container, err)
			continue
		}
		containers = append(containers, container)
	}

	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	for _, item := range containers {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			item.Id,
			item.Name,
			item.Pid,
			item.Status,
			item.Command,
			item.CreatedTime)
	}
	if err := w.Flush(); err != nil {
		log.Errorf("Flush error %v", err)
		return
	}
}

func getContainerInfo(file os.FileInfo) (*container.ContainerInfo, error) {
	containerName := file.Name()
	configDir := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	configFile := configDir + container.ConfigName

	content, err := ioutil.ReadFile(configFile)
	if err != nil {
                log.Errorf("get container %s info error: %v", configFile, err)
                return nil, err
	}
	var containerInfo container.ContainerInfo
	if err := json.Unmarshal(content, &containerInfo); err != nil {
		log.Errorf("Json unmarshal error %v", err)
		return nil, err
	}

	return &containerInfo, nil
}


