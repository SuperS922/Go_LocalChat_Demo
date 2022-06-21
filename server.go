package main

import (
	"fmt"
	"net"
	"sync"
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

	user := NewUser(connfd)

	// 用户上线，将用户加入 OnlineMap 中
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	// 广播当前用户上线消息
	this.BroadCast(user, "is online.")

	// 当前handler阻塞
	select {}
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
