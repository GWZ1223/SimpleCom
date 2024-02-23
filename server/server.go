package server

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User //在线用户对象是一个map集合userName:user
	MapLock   sync.RWMutex     //读写锁
	Message   chan string      //广播channel
}

// 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// 连接处理
func (s *Server) Handle(coon net.Conn) {
	//fmt.Println("连接成功")
	user := NewUser(coon)
	// 用户上线
	user.Online()

	// 开启一个goroutine 用于处理客户端消息
	go func() {
		buff := make([]byte, 4096)
		for {
			n, err := coon.Read(buff)
			if n == 0 {
				// 用户下线
				user.OffLine()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read Err: ...", err)
				return
			}
			msg := string(buff[:n-1])
			user.DoMessage(msg)
		}
	}()
}

// 广播消息
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + msg
	s.Message <- sendMsg
}

// 监听Message广播channel的goroutine 一旦有消息就发送给全部在线的User
func (s *Server) ListMessage() {
	for {
		msg := <-s.Message
		s.MapLock.Lock()
		for _, cli := range s.OnlineMap {
			cli.C <- msg
		}
		s.MapLock.Unlock()
	}
}

func (s *Server) Start() {
	// 创建连接
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen is err")
		err.Error()
		return
	}

	// 一直开启监听Message
	go s.ListMessage()

	// 关闭连接
	defer func(listen net.Listener) {
		err := listen.Close()
		if err != nil {
			fmt.Println("关闭连接失败", err.Error())
		}
	}(listen)

	for {
		coon, err := listen.Accept()
		if err != nil {
			fmt.Println("listen is not accept")
			err.Error()
			continue
		}
		// 业务操作
		go s.Handle(coon)
	}
}
