package server

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

func (s *Server) Handle(coon net.Conn) {
	fmt.Println("连接成功")
}

func (s *Server) Start() {
	// 创建连接
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen is err")
		err.Error()
		return
	}

	// 关闭连接
	defer listen.Close()

	for {
		coon, err := listen.Accept()
		if err != nil {
			fmt.Println("listen is not accept")
			err.Error()
			continue
		}
		go s.Handle(coon)
	}
}
