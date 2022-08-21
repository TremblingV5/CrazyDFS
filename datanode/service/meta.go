package DNService

import "github.com/TremblingV5/CrazyDFS/config/items"

func InitMeta(dn *DataNode, config items.DN) {
	dn.ReplicaName = ReplicaName(config.ReplicaName)
}
