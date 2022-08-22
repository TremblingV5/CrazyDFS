package CService

import "github.com/TremblingV5/CrazyDFS/proto"

type serverC2N struct {
	proto.UnimplementedClient2NNServer
}

type serverC2D struct {
	proto.UnimplementedClient2DNServer
}
