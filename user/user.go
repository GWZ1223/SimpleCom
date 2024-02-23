package user

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	Conn net.Conn
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
