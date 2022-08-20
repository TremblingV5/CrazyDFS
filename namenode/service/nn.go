package NNService

import (
	"time"

	"github.com/TremblingV5/CrazyDFS/config/items"
	"github.com/TremblingV5/CrazyDFS/utils"
	"github.com/TremblingV5/CrazyDFS/values"
)

var config, _ = utils.InitNodeConfig(items.NN{}, values.DataNodeConfigPath)

func (nn *NameNode) RegisterDN(name string, DNAddr string, diskUsage int64) {
	meta := DNMeta{
		Ip:          DNAddr,
		DiskUsage:   diskUsage,
		HeartbeatTS: time.Now().Unix(),
		Status:      DNUp,
	}
	nn.DNList[name] = &meta
}

func (nn *NameNode) HeartbeatReceiver() {
	duration := time.Second * time.Duration(config.HBTimeout)
	time.Sleep(duration)

	for id, dn := range nn.DNList {
		if time.Since(time.Unix(dn.HeartbeatTS, 0)) > duration {
			nn.DNList[id].Status = DNDown
		}
	}
	nn.HeartbeatReceiver()
}

func (nn *NameNode) GetDN() DNMeta {
	return DNMeta{}
}
