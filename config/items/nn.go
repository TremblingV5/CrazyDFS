package items

type NN struct {
	Port       string `yaml:"port"`
	Host       string `yaml:"host"`
	BlockSize  int64  `yaml:"blockSize"`
	ReplicaNum int64  `yaml:"replicaNum"`
	HBTimeout  int64  `yaml:"heartbeatTimeout"`
	DirTree    string `yaml:"dirTree"`
	Name       string `yaml:"name"`
	Path       string `yaml:"path"`
}
