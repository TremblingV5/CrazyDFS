package NNService

import (
	"io/ioutil"
	"time"

	"github.com/TremblingV5/CrazyDFS/config/items"
	"github.com/TremblingV5/CrazyDFS/proto"
	"github.com/TremblingV5/CrazyDFS/utils"
	"github.com/TremblingV5/CrazyDFS/values"
	"gopkg.in/yaml.v3"
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

func (nn *NameNode) GetDN() DNMeta {
	return DNMeta{}
}

func (nn *NameNode) InitIdleQueue(args *proto.BlockList) {
	tempIdleQueue := InitQueue()

	for _, item := range args.BlockL {
		tempIdleQueue.EnQueue(item.BlockName)
	}

	if nn.IdleQueue[ReplicaName(args.ReplicaName)] == nil {
		nn.IdleQueue[ReplicaName(args.ReplicaName)] = make(map[string]*DNBlockQueue)
	}

	nn.IdleQueue[ReplicaName(args.ReplicaName)][args.DNName] = tempIdleQueue
}

func (nn *NameNode) InitFile2Block(
	metaId *MetaId,
	args *proto.BlockList,
	path string,
) {
	utils.CheckAndMkdir(path + "/" + args.ReplicaName)

	for _, item := range args.BlockL {
		item = item
		newBlockMeta := BlockMeta{
			ID:          metaId.ToString(),
			ReplicaInfo: []ReplicaMeta{},
		}

		nn.FileToBlock[NNBlockID(metaId.ToString())] = &newBlockMeta
		yamlObj := NNBlockMeta{
			Id:      metaId.ToString(),
			BlockId: make(map[string]DNBlockID),
		}
		bytes, _ := yaml.Marshal(yamlObj)
		if err := ioutil.WriteFile(path+"/"+args.ReplicaName+"/"+metaId.ToString(), bytes, 0777); err != nil {
			utils.WriteLog(
				"error", "Write file defeat",
			)
		}
		*metaId = metaId.Next()
	}
}

func (nn *NameNode) ReadFile2BlockAndReplicaList(path string) {
	list, _ := ioutil.ReadDir(path)

	for _, item := range list {
		nn.ReplicaList[item.Name()] = []string{}

		nnBlocksList, _ := ioutil.ReadDir(path + "/" + item.Name())
		for _, block := range nnBlocksList {
			var nnBlockMeta BlockMeta
			utils.ReadYaml(
				path+"/"+item.Name()+"/"+block.Name(),
				&nnBlockMeta,
			)
			nn.FileToBlock[NNBlockID(block.Name())] = &nnBlockMeta
		}
	}
}

func (nn *NameNode) ReadDN2NNBlockMap(path string) {
	list, _ := ioutil.ReadDir(path)

	for _, item := range list {
		var temp DN2NNBlockMap
		utils.ReadYaml(path + "/" + item.Name(), &temp)
		nn.DN2NNBlockMap[DNBlockID(item.Name())] = &temp
	}
}

func (nn *NameNode) ReadIdleQueue(path string) {

}
