package main

import (
	log "github.com/Sirupsen/logrus"
	"mydocker/container"
	"mydocker/subsystems"
	"fmt"
	"strings"
	"github.com/urfave/cli"
)

var commitCommand =  cli.Command{
        Name: "commit",
        Usage: `commit a container`,
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command in commitCommand")
		}
	        log.Infof("commit come on")
        	imgName := context.Args().Get(0)
        	commitContainer(imgName)
        	return nil 
    	},
}

var runCommand = cli.Command{
        Name: "run",
        Usage: `create a container`,
        Flags: []cli.Flag{
                cli.BoolFlag{Name:"ti",Usage:"enable tty"},
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "v",
			Usage: "volume",
		},
        },
        Action: func(context *cli.Context) error {
        	log.Infof("run come on")
                if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
                }
		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}

		log.Infof("run command %s", strings.Join(cmdArray, " "))
                tty := context.Bool("ti")
                volume := context.String("v")
		res := &subsystems.ResourceConfig{MemoryLimit: context.String("m")}
                Run(tty, cmdArray, res, volume)
                return nil
        },
}

var initCommand = cli.Command{
    Name:    "init",
    Usage:    "Init container process run user's process in container. Do not call it outside",
    Action:    func(context *cli.Context) error {
        log.Infof("init come on")
        cmd := context.Args().Get(0)
        log.Infof("init command %s", cmd)
        err := container.RunContainerInitProcess()
        return err
    },
}
