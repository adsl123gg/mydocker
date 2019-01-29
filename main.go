package main

import(
	"log"
	"github.com/urfave/cli"
	"os"
)

const usage = "mydocker is a sample, created by hqc"

func main() {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = usage
	
	app.Commands = []cli.Command{
		initCommand,
		runCommand,
		commitCommand,
		listCommand,
		logCommand,
		execCommand,
		stopCommand,
		removeCommand,
	}

	if err:= app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

