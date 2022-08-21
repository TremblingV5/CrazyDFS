package NNService

import (
	"sync"
)

type ReplicaState int
type DNStatus int

type BlockMeta struct {
	ID          string
	ReplicaInfo []ReplicaMeta
}

type ReplicaMeta struct {
	BlockName string
	FileSize  int64
	Ip        string
	DNName    string
	State     ReplicaState
	ReplicaID string
}

type DNMeta struct {
	Ip          string
	DiskUsage   int64
	HeartbeatTS int64
	Status      DNStatus
}

const (
	DNUp           = DNStatus(1)
	DNDown         = DNStatus(2)
	ReplicaPending = DNStatus(3)
	ReplicaCommit  = DNStatus(4)
)

type NNBlockID string
type DNBlockID string

type NNAddr string

type NNRole int

type ReplicaName string

const (
	ActiveNN  = NNRole(0)
	StandByNN = NNRole(1)
)

type NameNode struct {
	FileToBlock   map[NNBlockID]*BlockMeta // NNBlockID单独对应一个块，每个BlockMeta含有一个ReplicaList，List中的数据均为这个块的副本
	DN2NNBlockMap map[DNBlockID]NNBlockID  // DNBlock到NNBlock的映射，每个DNBlock都能反向寻找到一个NNBlock
	// BlockToLocation map[string][]*ReplicaMeta
	DNList    map[string]*DNMeta
	IdleQueue map[ReplicaName]map[string]*DNBlockQueue

	ReplicaList  map[string][]string // 当前NameNode中有那些ReplicaList
	NameNodeList map[string][]NNAddr // 当前NameNode集合

	Role NNRole

	BlockSize     int64
	ReplicaFactor int64
}

type DirString string
type FileString string

type MetaId int64

type DirTree struct {
	Next map[string]DirTree

	Single       string
	Path         string
	IsDir        bool
	DirMetaInfo  DirMeta
	FileMetaInfo FileMeta
}

type DirMeta struct {
	Path string `yaml:"path"`
}

type FileMeta struct {
	Name       string            `yaml:"name"`
	Blocks     map[string]string `yaml:"blocks"`
	CrateTime  string            `yaml:"crate"`
	UpdateTime string            `yaml:"update"`
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

type NNBlockMeta struct {
	Id      string               `yaml:"id"`
	BlockId map[string]DNBlockID `yaml:"DNBlocks"`
}

type ReplicaMetaYamlList struct {
	Blocks map[string]ReplicaMetaYaml `yaml:"blocks"`
}

type ReplicaMetaYaml struct {
	BlockName string       `yaml:"blockName"`
	FileSize  int64        `yaml:"fileSize"`
	Ip        string       `yaml:"ip"`
	DNName    string       `yaml:"dn"`
	State     ReplicaState `yaml:"state"`
	ReplicaID string       `yaml:"replicaId"`
}
