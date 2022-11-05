package server

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// NewServer 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}

	return server
}

func (this *Server) Hander(conn net.Conn) {
	// 当前连接的业务
	fmt.Println("连接建立成功")
}

// Start 启动服务器的接口
func (this *Server) Start() {
	// 1 socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	// 1.1 close listen socket
	defer listener.Close()

	// 2 accept
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		// 3 do handler
		go this.Hander(conn)
	}
}
