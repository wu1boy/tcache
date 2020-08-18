package engine

import (
	"log"
	"tcache/conf"
)

func New(config *conf.Config) Cache {
	var c Cache

	if config.Engine == 1 {
		//内存存储
		c = MemoryEngineNew(config)
	}

	if c == nil {
		panic("unkown cache engine !")
	}

	log.Println("ready to server")

	return c

}
