package server

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	Conn   net.Conn
	server *Server
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
	if msg == "who" {
		// 查询当前在线用户
		u.server.MapLock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线.\n"
			u.sendMsg(onlineMsg)
		}
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// 修改用户名: 用户消息格式 "rename:用户名"
		newName := strings.Split(msg, ":")[1]
		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.sendMsg("用户名已存在,请重新输入.\n")
		} else {
			u.server.MapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.Name = newName
			u.server.OnlineMap[newName] = u
			u.server.MapLock.Unlock()
			u.sendMsg("修改成功" + u.Name + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			u.sendMsg("消息格式不正确")
			return
		}
		remoteUser, ok := u.server.OnlineMap[remoteName]
		if !ok {
			u.sendMsg("当前用户名不存在")
			return
		}
		content := strings.Split(msg, "|")[2]
		if content == "" {
			u.sendMsg("内容不能为空")
			return
		}
		remoteUser.sendMsg(u.Name + "对你说：" + content)
	} else {
		u.server.BroadCast(u, msg)
	}
}

func (u *User) sendMsg(msg string) {
	_, err := u.Conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("操作失败", err.Error())
		return
	}
}
