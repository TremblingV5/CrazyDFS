package NNService

import "sync"

type ReplicaState int
type DNStatus int

type BlockMeta struct {
	Name string
	Gs   int64
	ID   int64
}

type ReplicaMeta struct {
	BlockName string
	FileSize  int64
	Ip        string
	State     ReplicaState
	ReplicaID int64
}

type DNMeta struct {
	Ip         string
	DiskUsage  int64
	HearbeatTS int64
	Status     DNStatus
}

const (
	DNUp           = DNStatus(1)
	DNDown         = DNStatus(2)
	ReplicaPending = DNStatus(3)
	ReplicaCommit  = DNStatus(4)
)

type NameNode struct {
	FileToBlock     map[string][]BlockMeta
	BlockToLocation map[string][]ReplicaMeta
	DNList          []DNMeta

	BlockSize     int64
	ReplicaFactor int64
}

type Lease struct {
	Holder     string
	LastUpdate int64
	Path       *[]string
	// 允许抢占
	SoftLimit int64
	// 允许收回
	HardLimit int64
}

type LeaseMgr struct {
	FileToMetaMap map[string]Lease
	mu            sync.Mutex
}
