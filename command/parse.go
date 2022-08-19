package command

import (
	"fmt"
	"os"

	"github.com/TremblingV5/CrazyDFS/client"
	"github.com/TremblingV5/CrazyDFS/datanode"
	"github.com/TremblingV5/CrazyDFS/namenode"
	"github.com/TremblingV5/CrazyDFS/utils"
	"github.com/TremblingV5/CrazyDFS/values"
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
			Destination: &values.ShowVersion,
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
					Destination: &values.ClientConfigPath,
				},
			},
		},
		{
			Name:  "nn",
			Usage: "Start a namenode in command line",
			Action: func(c *cli.Context) {
				finalCommand = CMD_START_NN
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "config",
					Value:       "",
					Usage:       "Config file path for starting a namenode",
					Destination: &values.NameNodeConfigPath,
				},
			},
		},
		{
			Name:  "dn",
			Usage: "Start a datanode in command line",
			Action: func(c *cli.Context) {
				finalCommand = CMD_START_DN
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "config",
					Value:       "",
					Usage:       "Config file path for starting a datanode",
					Destination: &values.DataNodeConfigPath,
				},
			},
		},
		{
			Name:  "put",
			Usage: "Put data from local path to remote path",
			Action: func(c *cli.Context) {
				finalCommand = CMD_CLIENT_PUT
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "l",
					Value:       "",
					Usage:       "Local path",
					Destination: &values.LocalPath,
				},
				cli.StringFlag{
					Name:        "r",
					Value:       "",
					Usage:       "Remote path",
					Destination: &values.RemotePath,
				},
			},
		},
		{
			Name:  "get",
			Usage: "Get data from remote path to local path",
			Action: func(c *cli.Context) {
				finalCommand = CMD_CLIENT_GET
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "l",
					Value:       "",
					Usage:       "Local path",
					Destination: &values.LocalPath,
				},
				cli.StringFlag{
					Name:        "r",
					Value:       "",
					Usage:       "Remote path",
					Destination: &values.RemotePath,
				},
			},
		},
		{
			Name:  "rm",
			Usage: "Delete path",
			Action: func(c *cli.Context) {
				finalCommand = CMD_CLIENT_DELETE
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "r",
					Value:       "",
					Usage:       "Remote path",
					Destination: &values.RemotePath,
				},
			},
		},
		{
			Name:  "st",
			Usage: "List status for file in the provided path",
			Action: func(c *cli.Context) {
				finalCommand = CMD_CLIENT_STAT
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "r",
					Value:       "",
					Usage:       "Remote path",
					Destination: &values.RemotePath,
				},
			},
		},
		{
			Name:  "rn",
			Usage: "Rename a remote path to another",
			Action: func(c *cli.Context) {
				finalCommand = CMD_CLIENT_RENAME
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "s",
					Value:       "",
					Usage:       "Source path",
					Destination: &values.SrcPath,
				},
				cli.StringFlag{
					Name:        "t",
					Value:       "",
					Usage:       "Target path",
					Destination: &values.TargetPath,
				},
			},
		},
		{
			Name:  "md",
			Usage: "Make a directory",
			Action: func(c *cli.Context) {
				finalCommand = CMD_CLIENT_MKDIR
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "r",
					Value:       "",
					Usage:       "Remote path",
					Destination: &values.RemotePath,
				},
			},
		},
		{
			Name:  "ls",
			Usage: "List file list for provided path",
			Action: func(c *cli.Context) {
				finalCommand = CMD_CLIENT_LIST
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "r",
					Value:       "",
					Usage:       "Remote path",
					Destination: &values.RemotePath,
				},
			},
		},
	}

	appFlag.Action = func(c *cli.Context) error {
		if values.ShowVersion {
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
		fmt.Println("Start client by config path: " + values.ClientConfigPath)
	case CMD_START_NN:
		namenode.StartServer(values.NameNodeConfigPath)
	case CMD_START_DN:
		datanode.StartServer(values.DataNodeConfigPath)
	case CMD_CLIENT_PUT:
		client.PutHandle(values.LocalPath, values.RemotePath)
	case CMD_CLIENT_GET:
		client.GetHandle(values.LocalPath, values.RemotePath)
	case CMD_CLIENT_DELETE:
		client.DeleteHandle(values.RemotePath)
	case CMD_CLIENT_STAT:
		client.StatHandle(values.RemotePath)
	case CMD_CLIENT_RENAME:
		client.RenameHandle(values.SrcPath, values.TargetPath)
	case CMD_CLIENT_MKDIR:
		client.MkdirHandle(values.RemotePath)
	case CMD_CLIENT_LIST:
		client.ListHandle(values.RemotePath)
	}
}
