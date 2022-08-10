package items

type Total struct {
	Total struct {
		Version  string `yaml:"version"`
		Name     string `yaml:"name"`
		Usage    string `yaml:"usage"`
		HelpName string `yaml:"help_name"`
	} `yaml:"total"`
	Log struct {
		Level      string `yaml:"level"`
		Filename   string `yaml:"filename"`
		MaxSize    int    `yaml:"max_size"`
		MaxAge     int    `yaml:"max_age"`
		MaxBackups int    `yaml:"max_backups"`
	}
}
