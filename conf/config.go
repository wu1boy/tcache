package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Config struct {
	Shark          uint32      `yaml:"shark"`
	Port           string      `yaml:"port"`
	Addr           string      `yaml:"addr"`
	MaxMemory      interface{} `yaml:"maxmemory"`
	Engine         uint8       `yaml:"engine"`
	Eliminate      uint8       `yaml:"eliminate"`
	EliminatePolic uint8       `yaml:"eliminate-polic"`
	Hz             uint32      `yaml:"hz"`
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

func (c *Config) ParseMaxMemory() {
	if c.MaxMemory == "0" {
		c.MaxMemory = 0
		return
	}

	if len(c.MaxMemory.(string)) < 3 {
		log.Fatal("最大内存设置格式错误.请使用 KB MB GB 等")
	}

	maxmemory, err := strconv.ParseInt(strings.ToUpper(c.MaxMemory.(string)[:len(c.MaxMemory.(string))-2]), 10, 64)

	if err != nil {
		log.Println("最大内存设置格式错误.重置为默认值0. 如需修改,请使用 1KB 1MB 1GB 等类似格式,")
		return
	}

	unit := strings.ToUpper(c.MaxMemory.(string)[len(c.MaxMemory.(string))-2:])

	switch unit {
	case "KB":
		c.MaxMemory = maxmemory << 10
	case "MB":
		c.MaxMemory = maxmemory << 20
	case "GB":
		c.MaxMemory = maxmemory << 30
	default:
		log.Println("最大内存设置格式错误.重置为默认值0. 如需修改,请使用 1KB 1MB 1GB 等类似格式,")
	}
}
