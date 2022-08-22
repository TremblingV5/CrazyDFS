package datanode

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strconv"
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

var self = DNService.DataNode{}

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

	name := stream.BlockList.BlockL[0].BlockName
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

	res, err := c.HeartBeat(ctx, &proto.Heartbeat{
		Addr: utils.GetIP().String() + ":" + config.Port,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Heart beat: ", res)
	HeartBeat(config)
}

func BlockReport(config items.DN) {
	conn, err := grpc.Dial(config.NNHost+":"+config.NNPort, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := proto.NewDN2NNClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &proto.BlockList{
		BlockL: BlockList,
		Num:    int64(len(BlockList)),
		DiskUsage: func() int64 {
			num := 0
			for _, item := range BlockList {
				if item.ReplicaState == proto.BlockLocation_Using {
					num++
				}
			}
			return int64(num)
		}(),
		DNName:      config.Name,
		ReplicaName: config.ReplicaName,
	}
	res, err := c.BlockReport(ctx, req)
	if err != nil {
		panic(err)
	}

	fmt.Println("Block report: ", res)

	duration := time.Second * time.Duration(config.HeartBeatInterval)
	time.Sleep(duration * 5)

	BlockReport(config)
}

func Register(config items.DN) error {
	fmt.Println("Start to register data node to name node")

	conn, err := grpc.Dial(config.NNHost+":"+config.NNPort, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := proto.NewDN2NNClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	status, err := c.Register(ctx, &proto.RegisterDataNodeReq{
		New:  true,
		Addr: utils.GetIP().String() + ":" + config.Port,
		Name: config.Name,
		ReplicaName: config.ReplicaName,
		BlockNum: uint64(config.BlockNum),
		DiskUsage: 0,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Register: ", status)

	go HeartBeat(config)
	go BlockReport(config)

	return nil
}

func Init(config items.DN) {
	dataPath := config.Path + config.Name
	blockNum := config.BlockNum
	list, _ := ioutil.ReadDir(dataPath)
	if len(list) > 0 {
		// 已经初始化过
		BlockList = DNService.ReadInitedBlock(list)
		fmt.Println("Blocks has been initialized")
	} else {
		// 尚未初始化
		BlockList = DNService.InitBlock(dataPath, blockNum)
		fmt.Println("First time to init blocks")
	}

	DNService.InitMeta(&self, config)

	fmt.Println("Initialization successful with " + strconv.FormatInt(int64(len(BlockList)), 10) + " blocks")
}

func StartServer(path string) {
	config, err := utils.InitNodeConfig(items.DN{}, path)
	if err != nil {
		panic(err)
	}

	Init(config)

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
