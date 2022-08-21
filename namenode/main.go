package namenode

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"time"

	"github.com/TremblingV5/CrazyDFS/config/items"
	NNService "github.com/TremblingV5/CrazyDFS/namenode/service"
	"github.com/TremblingV5/CrazyDFS/proto"
	"github.com/TremblingV5/CrazyDFS/utils"
	"github.com/TremblingV5/CrazyDFS/values"
	"google.golang.org/grpc"
)

var Config, _ = utils.InitNodeConfig(items.NN{}, values.DataNodeConfigPath)
var CurrMetaId = NNService.MetaId(1)

var self = NNService.NameNode{
	FileToBlock:   make(map[NNService.NNBlockID]*NNService.BlockMeta),                 // 每次有新的slice时创建
	DN2NNBlockMap: make(map[NNService.DNBlockID]NNService.NNBlockID),                  // 向DN的Block写入数据后创建
	IdleQueue:     make(map[NNService.ReplicaName]map[string]*NNService.DNBlockQueue), // 有新的节点Block Report后创建，然后只维护
	DNList:        make(map[string]*NNService.DNMeta),                                 // DN列表，只存放DN的信息
	ReplicaList:   make(map[string][]string),                                          // Replica列表，第一个string是Replica的唯一Name，第二个string是DN的Name
	NameNodeList:  make(map[string][]NNService.NNAddr),                                // string为NN的名称，后续只记录地址
	Role:          NNService.ActiveNN,                                                 // 记录当前NN的角色是Active还是StandBy
	BlockSize:     Config.BlockSize,
	ReplicaFactor: Config.ReplicaNum,
}

type serverC2N struct {
	proto.UnimplementedClient2NNServer
}

type serverD2N struct {
	proto.UnimplementedDN2NNServer
}

// 接收DN的HeartBeat
func (s serverD2N) HeartBeat(ctx context.Context, args *proto.Heartbeat) (*proto.DatanodeOperation, error) {
	fmt.Println("Receive heartbeat: ", args)

	flag := false
	for _, item := range self.DNList {
		if item.Ip == args.Addr {
			flag = true
			item.HeartbeatTS = time.Now().Unix()
			item.DiskUsage = int64(args.DiskUsage)
		}
	}

	if flag {
		return &proto.DatanodeOperation{
			IpAddr: args.Addr,
		}, nil
	} else {
		return nil, nil
	}
}

func (s serverD2N) BlockReport(ctx context.Context, args *proto.BlockList) (*proto.OperationStatus, error) {
	fmt.Println("Receive block report from " + args.DNName)

	if _, ok := self.ReplicaList[args.ReplicaName]; ok {
		// 已经存在的副本集
		list := self.ReplicaList[args.ReplicaName]
		if utils.IsExist(list, args.DNName) {
			// 已经注册过的DN
			// 接收BlockReport并更新信息
		} else {
			// 尚未注册过的DN
			// 向IdleQueue添加未使用节点
			self.InitIdleQueue(args)
		}
	} else {
		// 尚未存在的副本集
		// 初始化FileToBlock
		// 向IdleQueue添加未使用节点
		self.InitFile2Block(
			&CurrMetaId,
			args,
			Config.Path,
		)
		self.InitIdleQueue(args)
		self.ReplicaList[args.ReplicaName] = append(self.ReplicaList[args.ReplicaName], args.ReplicaName)
	}

	return &proto.OperationStatus{
		Success: true,
	}, nil
}

func (s serverD2N) Register(ctx context.Context, args *proto.RegisterDataNodeReq) (*proto.OperationStatus, error) {
	newDN := NNService.DNMeta{
		Ip:          args.Addr,
		DiskUsage:   int64(args.DiskUsage),
		HeartbeatTS: time.Now().Unix(),
		Status:      NNService.DNUp,
	}
	self.DNList[args.Name] = &newDN

	self.ReplicaList[args.ReplicaName] = append(self.ReplicaList[args.ReplicaName], args.Name)

	res := &proto.OperationStatus{
		Success: true,
	}
	return res, nil
}

func (s serverC2N) FileOperation(ctx context.Context, mode *proto.FileOperationArgs) (*proto.FileLocationResp, error) {
	return nil, nil
}

func (s serverC2N) PutFile(ctx context.Context, args *proto.PathArgs) (*proto.FileLocationResp, error) {
	return nil, nil
}

func (s serverC2N) Meta(ctx context.Context, mode *proto.FileOperationArgs) (*proto.OperationStatus, error) {
	return nil, nil
}

func (s serverC2N) Rename(ctx context.Context, args *proto.RenameArgs) (*proto.OperationStatus, error) {
	return nil, nil
}

func (s serverC2N) GetFileMeta(ctx context.Context, args *proto.PathArgs) (*proto.FileMetaInfo, error) {
	return nil, nil
}

func (s serverC2N) GetDirMeta(ctx context.Context, args *proto.PathArgs) (*proto.DirMetaList, error) {
	return nil, nil
}

func (s serverC2N) PutComplete(ctx context.Context, args *proto.MetaStore) (*proto.OperationStatus, error) {
	return nil, nil
}

func (s serverC2N) RenewLock(ctx context.Context, args *proto.GetLeaseArgs) (*proto.OperationStatus, error) {
	return nil, nil
}

func Init(config items.NN) {
	NNService.InitBlockMetaList(&self, config)
	self.ReadFile2BlockAndReplicaList(config.Path)
	self.ReadDN2NNBlockMap(config.Path)
	self.ReadIdleQueue(config.Path)

	dataPath := config.Path + config.Name
	list, _ := ioutil.ReadDir(dataPath)

	max := NNService.MetaId(0)

	for _, item := range list {
		if item.Name() > max.ToString() {
			res, _ := strconv.Atoi(item.Name())
			max = NNService.MetaId(res)
		}
	}
}

func StartServer(path string) {
	config, _ := utils.InitNodeConfig(items.NN{}, path)
	Config = config

	Init(config)

	listen, err := net.Listen("tcp", "0.0.0.0:"+Config.Port)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	proto.RegisterDN2NNServer(s, &serverD2N{})
	proto.RegisterClient2NNServer(s, &serverC2N{})

	fmt.Println("Start name node server on " + utils.GetIP().String() + ":" + Config.Port + " successfully")

	if err := s.Serve(listen); err != nil {
		panic(err)
	}
}
