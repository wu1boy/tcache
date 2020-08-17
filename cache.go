package main

import (
	"log"
	"tcache/engine"
)

func New(engine int16) engine.Cache {
	var c engine.Cache

	if engine == 1 {
		//内存存储
		c = newInMemoryCache()
	}



	if c == nil {
		panic("unkown cache engine !")
	}

	log.Println("ready to server")

	return  c
}
