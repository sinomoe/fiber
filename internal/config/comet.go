package config

func GetComet(path string) *Comet {
	c, _ := Get(path)
	return c.(*Comet)
}

type Comet struct {
	WebsocketPort int `yaml:"websocket_port"`
	Rpc           Rpc `yaml:"rpc"`
}

type Rpc struct {
	Network string `yaml:"network"`
	Port    int    `yaml:"port"`
}
