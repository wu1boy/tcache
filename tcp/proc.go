package tcp

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

func (s *Service) readKey(r *bufio.Reader) (string,error) {

	l,err := readLen(r)

	if err != nil {
		return "", err
	}

	k := make([]byte,l)

	_,err = io.ReadFull(r,k)

	if err != nil {
		return "", err
	}

	return string(k), nil
}

func (s *Service) readKeyAndValue(r *bufio.Reader)(string,[]byte,error) {
	kl,err := readLen(r)

	if err != nil {
		return "", nil, err
	}

	vl,err := readLen(r)

	if err != nil {
		return "", nil, err
	}

	k := make([]byte,kl)
	v := make([]byte,vl)

	_,err = io.ReadFull(r,k)

	if err != nil {
		return "", nil, err
	}

	_,err = io.ReadFull(r,v)

	if err != nil {
		return "", nil, err
	}

	return string(k),v,nil

}

func readLen(r *bufio.Reader) (int,error) {
	s,err := r.ReadString(' ')

	if err != nil {
		return 0,err
	}

	fmt.Println(s)
	s = strings.TrimSpace(s)
	fmt.Println(s)

	l,err := strconv.Atoi(s)

	if err != nil {
		return 0,err
	}

	return l,nil
}

func sendResponse(value []byte,err error,conn net.Conn) error {
	if err != nil {
		serr := err.Error()
		tmp := fmt.Sprintf("-%d ",len(serr))+serr

		_,err = conn.Write([]byte(tmp))

		return err
	}

	vlen := fmt.Sprintf("%d ",len(value))
	_,err = conn.Write(append([]byte(vlen),value...))

	return err
}

func (s *Service) get(conn net.Conn,r *bufio.Reader) error {
	k,err := s.readKey(r)

	if err != nil {
		return err
	}

	v,err := s.Get(k)
	return sendResponse(v,err,conn)
}

func (s *Service) set(conn net.Conn,r *bufio.Reader) error {

	k,v,err := s.readKeyAndValue(r)

	if err != nil {
		return err
	}

	return sendResponse(nil,s.Set(k,v),conn)
}

func (s *Service) del(conn net.Conn,r *bufio.Reader) error {

	k,err := s.readKey(r)

	if err != nil {
		return err
	}

	return sendResponse(nil,s.Del(k),conn)
}

func (s *Service) process(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)

	for  {
		op,err := r.ReadByte()

		if err != nil {
			if err != io.EOF {
				log.Println("close connection read error!"+err.Error())
				return
			}
		}

		if op == 'S' {
			err = s.set(conn,r)
		} else if op == 'G' {
			err = s.get(conn,r)
		} else if op == 'D' {
			err = s.del(conn,r)
		} else {
			log.Println("invalid op close connection!")
			return
		}

		if err != nil {
			log.Println("op fail."+err.Error())
			return
		}
	}	
}
