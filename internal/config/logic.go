package config

func GetLogic(path string) *Logic {
	c, _ := Get(path)
	return c.(*Logic)
}

type Logic struct {
	Queue      Queue  `yaml:"queue"`
	Port       int    `yaml:"port"`
	AuthSecret string `yaml:"auth_secret"`
}

type Queue struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	Stream   string `yaml:"stream"`
	Group    string `yaml:"group"`
	DB       int    `yaml:"db"`
}
