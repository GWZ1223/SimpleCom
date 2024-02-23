package server

import (
	"SimpleCom/user"
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*user.User //在线用户对象是一个map集合userName:user
	mapLock   sync.RWMutex          //锁
	Message   chan string           //广播channel
}

// 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*user.User),
		Message:   make(chan string),
	}
	return server
}

func (s *Server) Handle(coon net.Conn) {
	//fmt.Println("连接成功")
	user := user.NewUser(coon)
	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()
	// 广播消息
	s.BroadCast(user, "已上线")
	select {}
}

// 广播消息
func (s *Server) BroadCast(user *user.User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + msg
	s.Message <- sendMsg
}

// 监听Message广播channel的goroutine 一旦有消息就发送给全部在线的User
func (s *Server) ListMessage() {
	for {
		msg := <-s.Message
		s.mapLock.Lock()
		for _, cli := range s.OnlineMap {
			cli.C <- msg
		}
		s.mapLock.Unlock()
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
