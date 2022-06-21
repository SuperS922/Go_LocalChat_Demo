package main

import (
	"fmt"
	"io"
	"net"
	"runtime"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	// 在线表 addr: *user
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播的channel
	Message chan string
}

// 创建一个 server 的接口
func NewServer(ip string, port int) *Server {
	server := Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return &server
}

// 监听Message广播消息channel的goroutine，一旦有消息就发送给全部的在线User
// 应当在Server启动时，就创建出来的goroutine
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message

		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// 广播消息的方法
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

func (this *Server) Handler(connfd net.Conn) {
	fmt.Println("连接建立成功")

	user := NewUser(connfd, this)

	user.Online()

	// 监听用户是否活跃的channel
	isLive := make(chan bool)

	// 接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)

		for {
			n, err := connfd.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("connfd.Read error.")
				return
			}

			// 提取用户的消息，去除 \n
			msg := string(buf[:n-1])

			// 用户针对 msg 进行消息处理
			user.DoMessage(msg)

			// 任意的操作都代表该用户是活跃的
			isLive <- true
		}
	}()

	// 当前handler阻塞
	for {
		select {
		case <-isLive:

		case <-time.After(time.Second * 30):
			user.SendMsg("长时间未活跃，你已被强制下线。\n")
			user.Offline()
			// 销毁资源
			close(user.C)
			connfd.Close()

			runtime.Goexit()
		}
	}
}

// 启动服务器的接口
func (this *Server) Start() {
	// socket listen
	listenfd, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err", err)
		return
	}
	// close listen socket
	defer listenfd.Close()

	// 启动监听Message的Goroutine
	go this.ListenMessager()

	for {
		// accept
		connfd, err := listenfd.Accept()
		if err != nil {
			fmt.Println("listenfd.Accept err", err)
		}

		// do handler
		go this.Handler(connfd)
	}
}
