package client

import (
	"context"
	"fmt"
	"time"

	"github.com/TremblingV5/CrazyDFS/config/items"
	"github.com/TremblingV5/CrazyDFS/proto"
	"github.com/TremblingV5/CrazyDFS/utils"
	"github.com/TremblingV5/CrazyDFS/values"
	"google.golang.org/grpc"
)

var config, _ = utils.InitNodeConfig(items.Client{}, values.ClientConfigPath)

func PutHandle(local string, remote string) {
	conn, err := grpc.Dial(
		config.NNAddr,
		grpc.WithInsecure(),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := proto.NewClient2NNClient(conn)
	ctx, cancel := context.WithTimeout(
		context.Background(), time.Second,
	)
	defer cancel()

	status, err := c.PutFile(ctx, &proto.PathArgs{
		PathName: remote,
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(status)
}

func GetHandle(local string, remote string) {

}

func DeleteHandle(remote string) {

}

func StatHandle(remote string) {

}

func RenameHandle(src string, target string) {

}

func MkdirHandle(remote string) {

}

func ListHandle(remote string) {

}
