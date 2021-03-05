package trace

import (
	"io/ioutil"
	"os"
	"trace/uuid"

	"gopkg.in/yaml.v2"
)

type ConfigT struct {
	Server Server `yaml:"server"`
	Mysql  Mysql  `yaml:"mysql"`
}

type TraceId struct {
	ProjectIdCheck    int8   `yaml:"project_id_check"`    // 强制设置为project_id
	ProjectIdStrategy int8   `yaml:"project_id_strategy"` // project_id选取策略key的方式: 1. 本机ip
	ProjectId         uint16 `yaml:"project_id"`          // 期望的project_id
}

type Server struct {
	ModuleName string  `yaml:"module_name"`
	TraceId    TraceId `yaml:"trace_id"`
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
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
}

var Config *ConfigT

func loadConfig(filepath string) (*ConfigT, error) {
	file, err := os.Open(filepath)

	if err != nil {
		return nil, err
	}

	defer file.Close()
	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, err
	}

	config := &ConfigT{}
	err = yaml.Unmarshal(bytes, config)

	if err != nil {
		return nil, err
	}

	return config, nil
}

func setDefault() {
	// projectId必须单实例唯一
	uuid.GlobalUUIDGenerator.ProjectId = Config.Server.TraceId.ProjectId
}

func init() {
	var err error
	Config, err = loadConfig("/Users/lurongming/sun-moon-star-star/trace/conf/server.yml")
	if err != nil {
		panic(err)
	}
	setDefault()
}
