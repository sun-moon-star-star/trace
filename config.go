package trace

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Mysql Mysql `yaml:"mysql"`
}

type Mysql struct {
	Hostname string `yaml:"hostname"`
	Port     uint32 `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Network  string `yaml:"network"`
	Database string `yaml:"database"`
}

var GlobalConfig *Config

func loadConfig(filepath string) (*Config, error) {
	file, err := os.Open(filepath)

	if err != nil {
		return nil, err
	}

	defer file.Close()
	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(bytes, config)

	if err != nil {
		return nil, err
	}

	return config, nil
}

func init() {
	var err error
	GlobalConfig, err = loadConfig("./conf/server.yml")
	if err != nil {
		panic(err)
	}
}
