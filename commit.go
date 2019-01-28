package main

import (
	"os/exec"
        "fmt"
        log "github.com/Sirupsen/logrus"
)

func commitContainer(imgName string) {
	mntURL := "/root/image/mnt"
	imgTar := "/root/image/imgs/" + imgName + ".tar"
	fmt.Println("imgName tar is %s ", imgTar)
	if _, err := exec.Command("tar", "-czf", imgTar, "-C", mntURL, ".").CombinedOutput(); err != nil {
		log.Errorf("tar %s error: %v", mntURL, err)
	}
}
