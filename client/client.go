package client

import (
	"flag"
	"fmt"
	"net"
)

var ServerIp string
var ServerPort int

func init() {
	flag.StringVar(&ServerIp, "ip", "127.0.0.1", "设置服务器的ip地址 默认为127.0.0.1")
	flag.IntVar(&ServerPort, "port", 8890, "设置服务器的端口，默认端口为8890")
}

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

// 创建连接函数
func NewClient(ServerIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   ServerIp,
		ServerPort: serverPort,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ServerIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial is Err")
	}
	client.conn = conn
	return client
}

func main() {
	//命令行解析
	flag.Parse()

	newClient := NewClient(ServerIp, ServerPort)
	if newClient == nil {
		fmt.Println(">>>>>>>>>>>> 连接服务器失败")
		return
	}
	fmt.Println(">>>>>>>>>>>> 连接服务器成功")
	select {}

}
