package main

import (
	log "github.com/Sirupsen/logrus"
	"mydocker/container"
	"mydocker/subsystems"
	"fmt"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
        Name: "run",
        Usage: `create a container`,
        Flags: []cli.Flag{
                cli.BoolFlag{Name:"ti",Usage:"enable tty"},
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
        },
        Action: func(context *cli.Context) error {
                if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
                }
		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}

                tty := context.Bool("ti")
		res := &subsystems.ResourceConfig{MemoryLimit: context.String("m")}
                Run(tty, cmdArray, res)
                return nil
        },
}

var initCommand = cli.Command{
    Name:    "init",
    Usage:    "Init container process run user's process in container. Do not call it outside",
    Action:    func(context *cli.Context) error {
        log.Infof("init come on")
        cmd := context.Args().Get(0)
        log.Infof("command %s", cmd)
        err := container.RunContainerInitProcess(cmd, nil)
        return err
    },
}
