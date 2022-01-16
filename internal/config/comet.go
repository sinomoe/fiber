package config

func GetComet(path string) *Comet {
	c, _ := Get(path)
	return c.(*Comet)
}

type Comet struct {
	WebsocketPort int      `yaml:"websocket_port"`
	Rpc           CometRpc `yaml:"rpc"`
	LogicUrl      string   `yaml:"logic_url"`
}

type CometRpc struct {
	Network string `yaml:"network"`
	Port    int    `yaml:"port"`
}
