package user

import (
	"SimpleCom/server"
	"net"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	Conn   net.Conn
	server *server.Server
}

// 创建一个User对象
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		Conn: conn,
	}

	// 当User被创建的时候就要进行监听消息
	go user.listenMessag()
	return user
}

// 监听消息
func (u *User) listenMessag() {
	for {
		msg := <-u.C
		u.Conn.Write([]byte(msg + "\n"))
	}
}

// 用户上线消息
func (u *User) Online() {
	u.server.MapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.MapLock.Unlock()

	u.server.BroadCast(u, "已上线")
}

// 用户下线消息
func (u *User) OffLine() {
	u.server.MapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.MapLock.Unlock()

	u.server.BroadCast(u, "已下线")
}

// 处理消息
func (u *User) DoMessage(msg string) {
	u.server.BroadCast(u, msg)
}
