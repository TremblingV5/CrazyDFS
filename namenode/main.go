package namenode

import (
	"context"
	"net"

	"github.com/TremblingV5/CrazyDFS/config/items"
	"github.com/TremblingV5/CrazyDFS/proto"
	"github.com/TremblingV5/CrazyDFS/utils"
	"github.com/TremblingV5/CrazyDFS/values"
	"google.golang.org/grpc"
)

var Config, _ = utils.InitNodeConfig(items.NN{}, values.DataNodeConfigPath)

type serverC2N struct {
	proto.UnimplementedClient2NNServer
}

type serverD2N struct {
	proto.UnimplementedDN2NNServer
}

// 接收DN的HeartBeat
func (s serverD2N) DNHeartBeat(ctx context.Context, heartbeat *proto.Heartbeat) (*proto.DatanodeOperation, error) {
	return nil, nil
}

func (s serverD2N) BlockReport(ctx context.Context, list *proto.BlockReplicaList) (*proto.OperationStatus, error) {
	return nil, nil
}

func (s serverD2N) RegisterDN(ctx context.Context, args *proto.RegisterDataNodeReq) (*proto.OperationStatus, error) {
	return nil, nil
}

func (s serverC2N) FileOperation(ctx context.Context, mode *proto.FileOperationArgs) (*proto.FileLocationResp, error) {
	return nil, nil
}

func (s serverC2N) CreateFile(ctx context.Context, mode *proto.FileOperationArgs) (*proto.FileLocationResp, error) {
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

func StartServer(path string) {
	listen, err := net.Listen("tcp", "0.0.0.0:"+Config.Port)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	proto.RegisterDN2NNServer(s, &serverD2N{})
	proto.RegisterClient2NNServer(s, &serverC2N{})

	if err := s.Serve(listen); err != nil {
		panic(err)
	}
}
