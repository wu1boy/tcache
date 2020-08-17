package main

import (
	"log"
	"runtime"
)

var lastTotalFreed uint64

func main()  {
	//configPath := flag.String("configPath","./tcache.yaml","配置文件路径")
	//flag.Parse()

	//conf.Config.GetConfig(*configPath)
	New()
	c := cache.New(*typ)

}


func printMemStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Printf("Alloc = %v TotalAlloc = %v  Just Freed = %v Sys = %v NumGC = %v\n",
		m.Alloc/1024, m.TotalAlloc/1024, ((m.TotalAlloc-m.Alloc)-lastTotalFreed)/1024, m.Sys/1024, m.NumGC)

	lastTotalFreed = m.TotalAlloc - m.Alloc
}