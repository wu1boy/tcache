package tcp

import (
	"log"
	"net"
	"tcache/engine"
)

type Service struct {
	engine.Cache
}

func (s *Service)Listen() {
	l,err := net.Listen("tcp","127.0.0.1:23456")

	if err != nil {
		panic(err.Error())
	}

	for {
		c,err := l.Accept()

		if err != nil {
			log.Println(err.Error())
			continue
		}

		go s.process(c)
	}
}

func New(c engine.Cache) *Service {
	return &Service{c}
}
