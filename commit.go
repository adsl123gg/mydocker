package main

import (
	"os/exec"
        "fmt"
	"mydocker/container"
        log "github.com/Sirupsen/logrus"
)

func commitContainer(containerName string, imgName string) {
	mntURL := fmt.Sprintf(container.MntLoc, containerName)
	imgTar := container.ImgURL + imgName + ".tar"
	fmt.Println("imgName tar is %s ", imgTar)
	if _, err := exec.Command("tar", "-czf", imgTar, "-C", mntURL, ".").CombinedOutput(); err != nil {
		log.Errorf("in commitContainer tar %s error: %v", mntURL, err)
	}
}
