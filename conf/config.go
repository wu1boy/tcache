package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Shark     uint32 `yaml:"shark"`
	Port      string `yaml:"port"`
	Addr      string `yaml:"addr"`
	MaxMemory string `yaml:"maxmemory"`
	Engine    uint8  `yaml:"engine"`
}

func GetConfig(filepath string) *Config {
	yamlFile, err := ioutil.ReadFile(filepath)

	if err != nil {
		log.Fatalf("打开配置文件错误: #%v ", err)
	}

	config := &Config{}
	err = yaml.Unmarshal(yamlFile, config)

	if err != nil {
		log.Fatalf("yaml格式错误: %v", err)
	}

	return config
}
