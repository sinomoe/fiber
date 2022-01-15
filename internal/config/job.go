package config

func GetJob(path string) *Job {
	c, _ := Get(path)
	return c.(*Job)
}

type Job struct {
	Rpc   JobRpc `yaml:"rpc"`
	Queue Queue  `yaml:"queue"`
}

type JobRpc struct {
	Network string `yaml:"network"`
	Port    int    `yaml:"port"`
	Retry   int    `yaml:"retry"`
}
