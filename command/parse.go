package command

import (
	"fmt"
	"os"

	"github.com/TremblingV5/CrazyDFS/utils"
	"github.com/urfave/cli"
)

func Parse(arguments []string) {
	var finalCommand Command

	appFlag := cli.NewApp()
	appFlag.Version = utils.TotalConf.Total.Version
	appFlag.HideVersion = false
	appFlag.Name = utils.TotalConf.Total.Name
	appFlag.Usage = utils.TotalConf.Total.Usage
	appFlag.HelpName = utils.TotalConf.Total.HelpName

	appFlag.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "version v",
			Usage:       "Show version",
			Destination: &ShowVersion,
		},
	}

	appFlag.Commands = []cli.Command{
		{
			Name:  "client",
			Usage: "Start a client in command line",
			Action: func(c *cli.Context) {
				finalCommand = CMD_START_CLIENT
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "config",
					Value:       "",
					Usage:       "Config file path for starting a client",
					Destination: &ClientConfigPath,
				},
			},
		},
		{
			Name:  "namenode",
			Usage: "Start a namenode in command line",
			Action: func(c *cli.Context) {
				finalCommand = CMD_START_NN
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "config",
					Value:       "",
					Usage:       "Config file path for starting a namenode",
					Destination: &NameNodeConfigPath,
				},
			},
		},
		{
			Name:  "datanode",
			Usage: "Start a datanode in command line",
			Action: func(c *cli.Context) {
				finalCommand = CMD_START_DN
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "config",
					Value:       "",
					Usage:       "Config file path for starting a datanode",
					Destination: &DataNodeConfigPath,
				},
			},
		},
	}

	appFlag.Action = func(c *cli.Context) error {
		if ShowVersion {
			cli.ShowVersion(c)
			os.Exit(0)
			return nil
		} else {
			cli.ShowAppHelp(c)
			os.Exit(0)
			return nil
		}
	}

	if err := appFlag.Run(arguments); err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	} else {
		Call(finalCommand)
	}
}

func Call(cmd Command) {
	switch cmd {
	case CMD_START_CLIENT:
		fmt.Println("Start client by config path: " + ClientConfigPath)
	case CMD_START_NN:
		fmt.Println("Start namenode by config path: " + NameNodeConfigPath)
	case CMD_START_DN:
		fmt.Println("Start datanode by config path: " + DataNodeConfigPath)
	}
}
