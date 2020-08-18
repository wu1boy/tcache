package tcp

import (
	"log"
	"net"
	"tcache/engine"
)

type Service struct {
	engine.Cache
}

func (s *Service) Listen(addr string) {

	l, err := net.Listen("tcp", addr)

	if err != nil {
		panic(err.Error())
	}

	for {
		client, err := l.Accept()

		if err != nil {
			log.Println(err.Error())
			continue
		}
		go s.process(client)
	}
}

func New(cache engine.Cache) *Service {
	return &Service{cache}
}
