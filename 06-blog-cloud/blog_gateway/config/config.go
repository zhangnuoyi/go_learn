package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Config 网关配置结构体
type Config struct {
	Name string `yaml:"Name"`
	Host string `yaml:"Host"`
	Port int    `yaml:"Port"`
	Jwt  JwtConfig  `yaml:"Jwt"`
	Upstreams []Upstream `yaml:"Upstreams"`
}

// JwtConfig JWT配置结构体
type JwtConfig struct {
	Secret       string   `yaml:"Secret"`
	ExcludePaths []string `yaml:"ExcludePaths"`
}

// Upstream 上游服务配置结构体
type Upstream struct {
	Name     string       `yaml:"Name"`
	Http     HttpConfig   `yaml:"Http"`
	Mappings []MappingRule `yaml:"Mappings"`
}

// HttpConfig HTTP配置结构体
type HttpConfig struct {
	Target  string `yaml:"Target"`
	Timeout int    `yaml:"Timeout"`
}

// MappingRule 映射规则结构体
type MappingRule struct {
	Method string `yaml:"Method"`
	Path   string `yaml:"Path"`
}

// LoadConfig 加载配置文件
func LoadConfig(filePath string) (*Config, error) {
	// 读取文件内容
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("读取配置文件失败: %v", err)
		return nil, err
	}

	// 解析YAML配置
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Printf("解析配置文件失败: %v", err)
		return nil, err
	}

	return &config, nil
}
