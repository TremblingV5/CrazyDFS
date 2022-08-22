package NNService

import (
	"io/ioutil"
	"os"
	"strconv"

	"github.com/TremblingV5/CrazyDFS/utils"
	"gopkg.in/yaml.v3"
)

func (id MetaId) ToString() string {
	str := strconv.FormatInt(int64(id), 10)
	res := ""

	for i := 1; i <= 32-len(str); i++ {
		res += "0"
	}

	res += str
	return res
}

func (id MetaId) Next() MetaId {
	id += 1
	return id
}

func ReadYamlMeta[T FileMeta | DirMeta](meta *T, path string) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		utils.WriteLog(
			"error", "Read meta yaml file defeat",
		)
	}

	if err = yaml.Unmarshal(file, meta); err != nil {
		utils.WriteLog(
			"error", "Unmarshal meta yaml file defeat",
		)
	}
}

func WriteYamlMeta[T FileMeta | DirMeta](meta *T, path string) {
	bytes, err := yaml.Marshal(meta)
	if err != nil {
		utils.WriteLog(
			"error", "Marshal meta yaml file defeat",
		)
	}

	if err := ioutil.WriteFile(path, bytes, 0600); err != nil {
		utils.WriteLog(
			"error", "Write meta yaml file defeat",
		)
	}

}

func GenFSImage(tree *DirTree, path string) {
	temp := make(map[string]*DirTree)
	items, _ := ioutil.ReadDir(path)
	for _, item := range items {
		if item.IsDir() {
			// 是一个目录，向下遍历
			dirMeta := DirMeta{}
			ReadYamlMeta(&dirMeta, path+"/"+"$META$")
			nextTree := DirTree{
				Next:         make(map[string]*DirTree),
				Single:       item.Name(),
				Path:         path,
				IsDir:        true,
				DirMetaInfo:  dirMeta,
				FileMetaInfo: FileMeta{},
			}
			GenFSImage(
				&nextTree,
				path+"/"+item.Name(),
			)
			temp[item.Name()] = &nextTree
		} else {
			// 是一个文件，登记信息，结束
			fileMeta := FileMeta{}
			ReadYamlMeta(&fileMeta, path+"/"+item.Name())
			nextTree := DirTree{
				Next:         nil,
				Single:       item.Name(),
				Path:         path,
				IsDir:        false,
				DirMetaInfo:  DirMeta{},
				FileMetaInfo: fileMeta,
			}
			temp[item.Name()] = &nextTree
		}
	}
	tree.Next = temp
}

func DownFSImage(tree *DirTree, path string) {
	if tree.Next == nil || len(tree.Next) == 0 {
		return
	} else {
		for _, item := range tree.Next {
			nextPath := path + "/" + item.Single
			if item.IsDir {
				// 是目录，先创建下一级目录然后写如目录的meta
				os.Mkdir(nextPath, 0777)
				WriteYamlMeta(&tree.DirMetaInfo, nextPath+"/"+"$META$")
			} else {
				// 不是目录，直接写入meta
				WriteYamlMeta(&tree.FileMetaInfo, nextPath)
			}
		}
	}
}
