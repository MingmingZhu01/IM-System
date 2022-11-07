package server

import (
	"example.com/IM-System/user"
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	// 在线用户的列表
	OnlineMap map[string]*user.User
	mapLock   sync.RWMutex

	// 消息广播的channel
	Message chan string
}

// NewServer 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*user.User),
		Message:   make(chan string),
	}

	return server
}

// ListenMessager 监听Message广播消息channel的goroutine，一旦有消息就发送给全部在线User
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message

		// 将msg发送个全部的在线User
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// BroadCast 广播消息的方法
func (this *Server) BroadCast(user *user.User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	// fmt.Println(sendMsg)
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	// 当前连接的业务
	fmt.Println("连接建立成功")
	user := user.NewUser(conn)
	// 用户上线了，将用户加入到onlineMap中
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	// 广播当前用户上线消息
	this.BroadCast(user, "已上线")

	// 当前handler阻塞
	select {}
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

	// 启动监听Message的goroutine
	go this.ListenMessager()

	// 2 accept
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		// 3 do handler
		go this.Handler(conn)
	}
}
