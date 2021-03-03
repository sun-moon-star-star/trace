package trace

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server Server `yaml:"server"`
	Mysql  Mysql  `yaml:"mysql"`
}

type Server struct {
	DefaultTimeLayout string `yaml:"default_time_layout"`
	LogTimeLayout     string `yaml:"log_time_layout"`
	BaggageTimeLayout string `yaml:"baggage_time_layout"`
}

type Mysql struct {
	Hostname string `yaml:"hostname"`
	Port     uint32 `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Network  string `yaml:"network"`
	Database string `yaml:"database"`

	TraceTableName   string `yaml:"trace_table_name"`
	SpanTableName    string `yaml:"span_table_name"`
	TagTableName     string `yaml:"tag_table_name"`
	LogTableName     string `yaml:"log_table_name"`
	BaggageTableName string `yaml:"baggage_table_name"`

	ConnMaxLifeTime uint32 `yaml:"conn_max_life_time"`
	MaxIdleConns    uint32 `yaml:"max_idle_conns"`
	MaxOpenConns    uint32 `yaml:"max_open_conns"`
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

func setDefault() {
	if GlobalConfig.Server.DefaultTimeLayout == "" {
		GlobalConfig.Server.DefaultTimeLayout = "2006-01-02 15:04:05.000000"
	}
	if GlobalConfig.Server.LogTimeLayout == "" {
		GlobalConfig.Server.LogTimeLayout = GlobalConfig.Server.DefaultTimeLayout
	}
	if GlobalConfig.Server.BaggageTimeLayout == "" {
		GlobalConfig.Server.BaggageTimeLayout = GlobalConfig.Server.DefaultTimeLayout
	}
}

func init() {
	var err error
	GlobalConfig, err = loadConfig("./conf/server.yml")
	if err != nil {
		panic(err)
	}
	setDefault()
}
