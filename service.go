package main

import (
	"fmt"
	"tcache/timewheel"
	"time"
)

//var lastTotalFreed uint64

func main() {
	tw := timewheel.New(2 * time.Second,3,  func(data interface{}) {
		fmt.Println(11111)
	})
	tw.Start()
	tw.AddTimer(1 * time.Second,"k1","")
	tw.AddTimer(3 * time.Second,"k2","")

	select {

	}
	//configPath := flag.String("configPath","./tcache.yaml","配置文件路径")
	//flag.Parse()
	//
	//config := conf.GetConfig(*configPath)
	//cache := engine.New(config)
	//
	//var buffer strings.Builder
	//buffer.WriteString(config.Addr)
	//buffer.WriteString(":")
	//buffer.WriteString(config.Port)
	//
	//tcp.New(cache).Listen(buffer.String())

}





//func strtob(s string) []byte {
//	x := (*[2]uintptr)(unsafe.Pointer(&s))
//	h := [3]uintptr{x[0], x[1], x[1]}
//	d := *(*[]byte)(unsafe.Pointer(&h))
//
//	return d
//}
//
//func bytes2str(b []byte) string {
//	return *(*string)(unsafe.Pointer(&b))
//}

//func printMemStats() {
//	var m runtime.MemStats
//	runtime.ReadMemStats(&m)
//	log.Printf("Alloc = %v TotalAlloc = %v  Just Freed = %v Sys = %v NumGC = %v\n",
//		m.Alloc/1024, m.TotalAlloc/1024, ((m.TotalAlloc-m.Alloc)-lastTotalFreed)/1024, m.Sys/1024, m.NumGC)
//
//	lastTotalFreed = m.TotalAlloc - m.Alloc
//}
