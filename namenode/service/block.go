package NNService

import (
	"io/fs"
	"io/ioutil"

	"github.com/TremblingV5/CrazyDFS/config/items"
	"github.com/TremblingV5/CrazyDFS/utils"
	"gopkg.in/yaml.v3"
)

func InitBlockMetaList(nn *NameNode, config items.NN) {
	dataPath := config.Path + config.Name
	list, _ := ioutil.ReadDir(dataPath)

	if len(list) > 0 {
		ReadInitedBlockMeta(
			nn, config, list, dataPath,
		)
	} else {

	}
}

func ReadInitedBlockMeta(
	nn *NameNode,
	config items.NN,
	fileList []fs.FileInfo,
	path string,
) {
	blockMetaList := make(map[NNBlockID]*BlockMeta)
	tempDN2NNBlockMap := make(map[DNBlockID]NNBlockID)

	for _, item := range fileList {
		bytes, _ := ioutil.ReadFile(path + "\\" + item.Name())
		var blocks ReplicaMetaYamlList

		if err := yaml.Unmarshal(bytes, blocks); err != nil {
			utils.WriteLog(
				"error", "Unmarshal block meta file defeat",
			)
		}

		var replicas []ReplicaMeta

		for _, value := range blocks.Blocks {
			newReplicaMeta := ReplicaMeta{
				BlockName: value.BlockName,
				FileSize:  value.FileSize,
				Ip:        value.Ip,
				DNName:    value.DNName,
				State:     value.State,
				ReplicaID: value.ReplicaID,
			}
			replicas = append(replicas, newReplicaMeta)
			tempDN2NNBlockMap[DNBlockID(value.DNName+"."+value.BlockName)] = NNBlockID(item.Name())
		}

		newBlockMeta := BlockMeta{
			ID:          item.Name(),
			ReplicaInfo: replicas,
		}
		blockMetaList[NNBlockID(item.Name())] = &newBlockMeta
	}

	nn.FileToBlock = blockMetaList
}
