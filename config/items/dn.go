package items

type DN struct {
	Port              string `yaml:"port"`
	NNHost            string `yaml:"nnHost"`
	NNPort            string `yaml:"nnPort"`
	HeartBeatInterval int64  `yaml:"interval"`
	IOSize            int64  `yaml:"ioSize"`
	BlockSize         int64  `yaml:"blockSize"`
}
