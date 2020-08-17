package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Shark     int    `yaml:"shark"`
	Port      int    `yaml:"port"`
	Addr      string `yaml:"addr"`
	MaxMemory string `yaml:"maxmemory"`
}

func (c *Config) GetConfig(filepath string) {
	yamlFile, err := ioutil.ReadFile(filepath)

	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
		panic("")
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		panic("")
	}
}
