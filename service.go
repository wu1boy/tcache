package main

import (
	"flag"
	"fmt"
	"runtime"
	"tcache/conf"
)

//var lastTotalFreed uint64

func main() {

	configPath := flag.String("configPath","./tcache.yaml","配置文件路径")
	flag.Parse()

	config := conf.GetConfig(*configPath)
	//解析maxmemory 的值.
	config.ParseMaxMemory()


	//////////
	var ms *runtime.MemStats
	runtime.ReadMemStats(ms)
	fmt.Println(ms.Alloc)

	return
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
