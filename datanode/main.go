package datanode

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/TremblingV5/CrazyDFS/config/items"
	DNService "github.com/TremblingV5/CrazyDFS/datanode/service"
	"github.com/TremblingV5/CrazyDFS/proto"
	"github.com/TremblingV5/CrazyDFS/utils"
	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedClient2DNServer
}

var BlockList = []*proto.BlockLocation{}

func (s server) GetBlock(mode *proto.FileOperationArgs, blockServer proto.Client2DN_GetBlockServer) error {
	b := DNService.GetBlock(mode.FileName, "r")

	for b.HasNextChunk() {
		chunk, n, err := b.GetNextBlock()
		if err != nil {
			return err
		}
		blockServer.Send(&proto.File{
			Content: (*chunk)[:n],
		})
	}

	b.Close()

	return nil
}

func (s server) PutBlock(blockServer proto.Client2DN_WriteBlockServer) error {
	stream, err := blockServer.Recv()

	if err == io.EOF {
		blockServer.SendAndClose(&proto.OperationStatus{Success: false})
	}

	name := stream.BlockReplicaList.BlockReplicaList[0].BlockName
	b := DNService.GetBlock(name, "w")
	file := make([]byte, 0)

	for {
		stream, err := blockServer.Recv()
		if err == io.EOF {
			b.Close()
			blockServer.SendAndClose(&proto.OperationStatus{Success: true})
			break
		}
		content := stream.File.Content
		if err = b.WriteChunk(content); err != nil {
			blockServer.SendAndClose(&proto.OperationStatus{Success: false})
		}

		file = append(file, content...)
	}

	BlockList = append(BlockList, &proto.BlockLocation{
		BlockName: name,
		IpAddr:    string(utils.GetIP()),
		BlockSize: b.GetFileSize(),
	})

	return nil
}

func HeartBeat(config items.DN) {
	duration := time.Second * time.Duration(config.HeartBeatInterval)
	time.Sleep(duration)

	conn, err := grpc.Dial(config.NNHost+":"+config.NNPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := proto.NewDN2NNClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.HeartBeat(ctx, &proto.Heartbeat{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Heart beat: ", res)
	HeartBeat(config)
}

func BlockReport(config items.DN) {
	duration := time.Second * time.Duration(config.HeartBeatInterval)
	time.Sleep(duration * 10)

	conn, err := grpc.Dial(config.NNHost+":"+config.NNPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := proto.NewDN2NNClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.BlockReport(ctx, &proto.BlockReplicaList{
		BlockReplicaList: BlockList,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Block report: ", res)
	BlockReport(config)
}

func Register(config items.DN) error {
	conn, err := grpc.Dial(config.NNHost+":"+config.NNPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := proto.NewDN2NNClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	status, err := c.Register(ctx, &proto.RegisterDataNodeReq{
		New: true,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Register: ", status)

	go HeartBeat(config)
	go BlockReport(config)

	return nil
}

func StartServer(path string) {
	config, err := utils.InitNodeConfig(items.DN{}, path)
	if err != nil {
		panic(err)
	}

	listen, err := net.Listen("tcp", "0.0.0.0:"+config.Port)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	if err = Register(config); err != nil {
		panic(err)
	}

	proto.RegisterClient2DNServer(s, &server{})
	if err = s.Serve(listen); err != nil {
		panic(err)
	}
}
