package config

import "github.com/BurntSushi/toml"

var (
	global *Config
)

// LoadGlobalConfig 加载全局配置
func LoadGlobalConfig(fpath string) error {
	c, err := ParseConfig(fpath)
	if err != nil {
		return err
	}
	global = c
	return nil
}

// GetGlobalConfig 获取全局配置
func GetGlobalConfig() *Config {
	if global == nil {
		return &Config{}
	}
	return global
}

// ParseConfig 解析配置文件
func ParseConfig(fpath string) (*Config, error) {
	var c Config
	_, err := toml.DecodeFile(fpath, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

type Config struct {
	Http Http `toml:"http"`
	DB   Postgres `toml:"db"`
	Base Base `toml:"base"`
}


type Postgres struct {
	Addr string `toml:"addr"`
	User string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
	Init bool `toml:"init"`
	OpenListener bool `toml:"open_listener"`
}

type Http struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
}


type Base struct {
	InitSwagger bool `toml:"init_swagger"`
	OpenPprof   bool `toml:"open_pprof"`
	CacheLuaVm  bool `toml:"cache_lua_vm"`
	OpenGrpc bool `toml:"open_grpc"`
	GrpcPort string `toml:"grpc_port"`
}
