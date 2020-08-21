package engine

import (
	"log"
	"tcache/conf"
)

func New(config *conf.Config) Cache {
	var cache Cache

	if config.Engine == 1 {
		//内存存储
		cache = MemoryEngineNew(config)

	}

	if cache == nil {
		panic("unkown cache engine !")
	}

	log.Println("ready to server")

	return cache

}
