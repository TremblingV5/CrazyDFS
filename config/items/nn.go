package items

type NN struct {
	Port       string `yaml:"port"`
	Host       string `yaml:"host"`
	BlockSize  string `yaml:"blockSize"`
	ReplicaNum string `yaml:"replicaNum"`
	HBTimeout  int64  `yaml:"heartbeatTimeout"`
}
