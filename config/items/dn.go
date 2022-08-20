package items

type DN struct {
	Port              string `yaml:"port"`
	NNHost            string `yaml:"nnHost"`
	NNPort            string `yaml:"nnPort"`
	HeartBeatInterval int64  `yaml:"interval"`
	IOSize            int64  `yaml:"ioSize"`
	BlockSize         int64  `yaml:"blockSize"`
	BlockNum          int64  `yaml:"blockNum"`
	Name              string `yaml:"name"`
	ReplicaName       string `yaml:"replicaName"`
	Path              string `yaml:"path"`
}
