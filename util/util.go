package util 

import (
    "encoding/json"
    "os"
    "fmt"
)

type Configuration struct {
    init bool
    RootURL   string
    MntURL   string
    ImgURL   string
    ContainerInfoLocation	string
    IPAMPath	string
    NetworkPath	string
    EndpointPath	string
}

const configLoc = "/root/image/config/conf.json"
var configuration Configuration

func readConfig() {
	file, _ := os.Open(configLoc)
	defer file.Close()

	decoder := json.NewDecoder(file)
	configuration = Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("read config file error:", err)
	}

}

func GetRootURL() string {
	if configuration.init == false {
		readConfig()
	}
	return configuration.RootURL
}

func GetMntURL() string {
	if configuration.init == false {
                readConfig()
        }
	return configuration.MntURL
}

func GetImgURL() string {
        if configuration.init == false {
                readConfig()
        }
        return configuration.ImgURL
}

func GetContainerInfoLocation() string {
        if configuration.init == false {
                readConfig()
        }
        return configuration.ContainerInfoLocation
}


func GetIPAMPath() string {
        if configuration.init == false {
                readConfig()
        }
        return configuration.IPAMPath
}

func GetNetworkPath() string {
        if configuration.init == false {
                readConfig()
        }
        return configuration.NetworkPath
}

func GetEndPointPath() string {
	if configuration.init == false {
                readConfig()
        }
        return configuration.EndpointPath
}
