package container 

import(
	"os"
	"strings"
	"os/exec"
	"syscall"
	"mydocker/util"
	"fmt"
	log "github.com/Sirupsen/logrus"
)

var (
	RUNNING             string = "running"
	STOP                string = "stopped"
	Exit                string = "exited"
	DefaultInfoLocation string = util.GetContainerInfoLocation()
	ConfigName          string = "config.json"
	LogConfigName       string = "container.log"
	MntLoc              string =  util.GetMntURL()
        RootLoc             string =  util.GetRootURL()
        ImgURL              string =  util.GetImgURL()
)

type ContainerInfo struct {
	Pid         string `json:"pid"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Command     string `json:"command"`
	CreatedTime string `json:"createTime"`
	Status      string `json:"status"`
	ImageName   string `json:"ImageName"`
	IP	    string `json:"IP"`
	PortMapping []string `json:"portmapping"`
}

func NewParentProcess(tty bool, volume string, imageName string, containerName string, envSlice []string) (*exec.Cmd, *os.File) {
    readPipe, writePipe, err := NewPipe()
    if err != nil {
	log.Errorf("create pipe error %v", err)
	return nil, nil
    }

    cmd := exec.Command("/proc/self/exe", "init")
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
        syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	Unshareflags: syscall.CLONE_NEWNS,
    }

    if tty {
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
    } else {
	dirUrl := fmt.Sprintf(DefaultInfoLocation, containerName)
	if err := os.MkdirAll(dirUrl, 0622); err != nil {
                log.Errorf("Mkdir %s error %v", dirUrl, err)
                return nil, nil
        }

	logFileName := dirUrl + "/" + LogConfigName
        logFile, err := os.Create(logFileName)
	if err != nil {
		log.Errorf("create log file %s error %v", logFileName, err)
                return nil, nil
	}

	cmd.Stdout = logFile
	cmd.Stderr = logFile
    }
    cmd.ExtraFiles = []*os.File{readPipe}
    cmd.Env = append(os.Environ(), envSlice...)

    mntURL := NewWorkSpace(volume, imageName, containerName)
    cmd.Dir = mntURL
    println("NewParentProcess cmd.Dir is: " + cmd.Dir)
    return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil { return nil, nil, err}
	return read, write, err
}

func DeleteWorkSpace(containerName string, volume string) {
	rootURL := fmt.Sprintf(RootLoc, containerName)
        mntURL := fmt.Sprintf(MntLoc, containerName)

	if volume != "" {
		volumeURLs := volumeUrlExtract(volume)
                length := len(volumeURLs)
                if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
                        DeleteMountPointWithVolume(mntURL, volumeURLs)
                        log.Infof("unmount volume %s ", volumeURLs)
                } else {
                        log.Warningf("unmount volume parameter invalid")
                }
	}
	DeleteMountPoint(mntURL)
	//DeleteWriteLayer(rootURL)
	DeleteContainerDir(rootURL)
}

func DeleteMountPointWithVolume(mntURL string, volumeURLs []string) {
	containerUrl := mntURL + volumeURLs[1]
	cmd := exec.Command("umount", containerUrl)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        if err := cmd.Run(); err != nil {
                log.Errorf("DeleteMountPoint umount error: %v", err)
        }
}

func DeleteMountPoint(mntURL string) {
	cmd := exec.Command("umount", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("DeleteMountPoint umount error: %v", err)
		return
	}
	if err := os.RemoveAll(mntURL); err != nil {
                log.Errorf("DeleteMountPoint Remove error: %v", err)
        }
}

//func DeleteWriteLayer(rootURL string) {
func DeleteContainerDir(rootURL string) {
	//writeURL := rootURL + "writeLayer/"
        if PathExists(rootURL) {
                if err := os.RemoveAll(rootURL); err != nil {
                        log.Errorf("DeleteContainerDir remove %s error: %v", rootURL, err)
                }
        }
}

func NewWorkSpace(volume string, imageName string, containerName string) string {
	rootURL := fmt.Sprintf(RootLoc, containerName)
	mntURL := fmt.Sprintf(MntLoc, containerName)

	if !PathExists(rootURL) {
                if err := os.Mkdir(rootURL, 0777); err != nil {
                        log.Errorf("Mkdir %s of container %s root error %v", rootURL, containerName, err)
			return ""
                }
	}
        CreateReadOnlyLayer(rootURL, imageName)
        CreateWriteLayer(rootURL)
        CreateMountPoint(rootURL, mntURL, imageName)
	if volume != "" {
		volumeURLs := volumeUrlExtract(volume)
		length := len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			MountVolume(rootURL, mntURL, volumeURLs)
			log.Infof("mount volume %s ", volumeURLs)
		} else {
			log.Warningf("mount volume parameter invalid")
		}
	}
	return mntURL
}

func volumeUrlExtract(volume string) []string{
	var volumeURLs []string
	volumeURLs = strings.Split(volume, ":")
	return volumeURLs
}

func MountVolume(rootURL string, mntURL string, volumeURLs []string) {
	parentUrl := volumeURLs[0]
	if !PathExists(parentUrl) {
                if err := os.Mkdir(parentUrl, 0777); err != nil {
                        log.Errorf("Mkdir parentUrl %s error %v", parentUrl, err)
                }
	}
	containerUrl := volumeURLs[1]
	containerUrl = mntURL + containerUrl
        if !PathExists(containerUrl) {
                if err := os.Mkdir(containerUrl, 0777); err != nil {
                        log.Errorf("Mkdir containerUrl %s error %v", containerUrl, err)
                }
        }

	dirs := "dirs=" + parentUrl
        cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerUrl)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        if err := cmd.Run(); err != nil {
                log.Errorf("MountVolume error: %v", err)
        }
}

func CreateReadOnlyLayer(rootURL string, imageName string) {
        busyboxURL := rootURL + "busybox/"
        busyboxTarURL := ImgURL + imageName + ".tar"
        //busyboxTarURL := "./resource/busybox.tar"
        if !PathExists(busyboxURL) {
                if err := os.Mkdir(busyboxURL, 0777); err != nil {
                        log.Errorf("Mkdir %s error %v", busyboxURL, err)
                }
                if _, err := exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxURL).CombinedOutput(); err != nil {
                        log.Errorf("tar -xvf  %s error %v", busyboxURL, err)
                }
        }
}

func CreateWriteLayer(rootURL string) {
        writeURL := rootURL + "writeLayer/"
        if !PathExists(writeURL) {
                if err := os.Mkdir(writeURL, 0777); err != nil {
                        log.Errorf("Mkdir %s error %v", writeURL, err)
                }
        }
}

func CreateMountPoint(rootURL string, mntURL string, imageName string) {
        if !PathExists(mntURL) {
                if err := os.Mkdir(mntURL, 0777); err != nil {
                        log.Errorf("Mkdir %s error %v", mntURL, err)
                }
        }

	dirs := "dirs=" + rootURL + "writeLayer:" + rootURL + imageName
        cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        if err := cmd.Run(); err != nil {
                log.Errorf("%v", err)
        }
}

func PathExists(path string) bool {
        if _, err := os.Stat(path); os.IsNotExist(err) {
                return false
        }
        return true
}


