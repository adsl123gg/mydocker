package main

import (
	log "github.com/Sirupsen/logrus"
	"mydocker/container"
	"mydocker/network"
	"mydocker/subsystems"
	"fmt"
	"strings"
	"github.com/urfave/cli"
	"os"
)

var networkCommand =  cli.Command{
        Name: "network",
        Usage: `operate container network`,
	Subcommands: []cli.Command {
		{
			Name: "create",
			Usage: "create a container network",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "driver",
					Usage: "network driver",
				},
				cli.StringFlag{
					Name:  "subnet",
					Usage: "subnet cidr",
				},
			},
			Action:func(context *cli.Context) error {
				if len(context.Args()) < 1 {
					return fmt.Errorf("Missing network name")
				}
				network.Init()
				err := network.CreateNetwork(context.String("driver"), context.String("subnet"), context.Args()[0])
				if err != nil {
					return fmt.Errorf("create network error: %+v", err)
				}
				return nil
			},
		},
		{
			Name: "list",
			Usage: "list container network",
			Action:func(context *cli.Context) error {
				network.Init()
				network.ListNetwork()
				return nil
			},
		},
		{
			Name: "remove",
			Usage: "remove container network",
			Action:func(context *cli.Context) error {
				if len(context.Args()) < 1 {
					return fmt.Errorf("Missing network name")
				}
				network.Init()
				err := network.DeleteNetwork(context.Args()[0])
				if err != nil {
					return fmt.Errorf("remove network error: %+v", err)
				}
				return nil
			},
		},
	},
}

var removeCommand =  cli.Command{
        Name: "rm",
        Usage: `remove the stopped container`,
        Action: func(context *cli.Context) error {
                if len(context.Args()) < 1 {
                        return fmt.Errorf("Missing container command in stopCommand")
                }
                log.Infof("remove come on")
                containerName := context.Args().Get(0)
                RemoveContainer(containerName)
                return nil
        },
}

var stopCommand =  cli.Command{
        Name: "stop",
        Usage: `stop the container`,
        Action: func(context *cli.Context) error {
                if len(context.Args()) < 1 {
                        return fmt.Errorf("Missing container command in stopCommand")
                }
                log.Infof("stop come on")
                containerName := context.Args().Get(0)
                StopContainer(containerName)
                return nil
        },
}

var execCommand =  cli.Command{
        Name: "exec",
        Usage: `exec the cmd into container`,
        Action: func(context *cli.Context) error {
		if os.Getenv(ENV_EXEC_PID) != "" {
			log.Infof("pid callback pid %s", os.Getgid)
			return nil
		}
                if len(context.Args()) < 2 {
                        return fmt.Errorf("Missing container command in execCommand")
                }
                log.Infof("exec come on")
                containerName := context.Args().Get(0)
		var cmdArray []string
                for _, arg := range context.Args() {
                        cmdArray = append(cmdArray, arg)
                }
		cmdArray = cmdArray[1:]
		log.Infof("exec cmdArray is: %v", cmdArray)
                ExecContainer(containerName, cmdArray)
                return nil
        },
}

var logCommand =  cli.Command{
        Name: "log",
        Usage: `print the container log`,
        Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
                        return fmt.Errorf("Missing container command in logCommand")
                }
                log.Infof("log come on")
                containerName := context.Args().Get(0)
                LogContainer(containerName)
                return nil
        },
}

var listCommand =  cli.Command{
        Name: "ps",
        Usage: `list the containers`,
        Action: func(context *cli.Context) error {
                log.Infof("list come on")
                ListContainer()
                return nil
        },
}

var commitCommand =  cli.Command{
        Name: "commit",
        Usage: `commit a container`,
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 2 {
			return fmt.Errorf("Missing container command in commitCommand")
		}
	        log.Infof("commit come on")
        	containerName := context.Args().Get(0)
        	imgName := context.Args().Get(1)
        	commitContainer(containerName, imgName)
        	return nil 
    	},
}

var runCommand = cli.Command{
        Name: "run",
        Usage: `create a container`,
        Flags: []cli.Flag{
                cli.BoolFlag{Name:"ti",Usage:"enable tty"},
                cli.BoolFlag{Name:"d",Usage:"detach container"},
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.StringSliceFlag{
                        Name:  "e",
                        Usage: "enviroment variables",
                },
		cli.StringFlag{
			Name:  "net",
			Usage: "network",
		},
		cli.StringSliceFlag{
                        Name:  "p",
                        Usage: "port mapping",
                },
		cli.StringFlag{
			Name:  "v",
			Usage: "volume",
		},
		cli.StringFlag{
                        Name:  "name",
                        Usage: "container name",
                },
        },
        Action: func(context *cli.Context) error {
        	log.Infof("run come on")
                if len(context.Args()) < 2 {
			return fmt.Errorf("Missing container command in runCommand")
                }
		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}
		imageName := cmdArray[0]
		cmdArray = cmdArray[1:]

		log.Infof("run command %s", strings.Join(cmdArray, " "))
                tty := context.Bool("ti")
                detach := context.Bool("d")
		if tty && detach {
			return fmt.Errorf("ti and d can't exist at same time")
		}
                volume := context.String("v")
                network := context.String("net")
                envSlice := context.StringSlice("e")
                portMapping := context.StringSlice("p")
                containerName := context.String("name")
		res := &subsystems.ResourceConfig{MemoryLimit: context.String("m")}
                Run(tty, cmdArray, res, volume, imageName, containerName, envSlice, network, portMapping)
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
