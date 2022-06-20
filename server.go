package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// 创建一个 server 的接口
func NewServer(ip string, port int) *Server {
	server := Server{Ip: ip, Port: port}
	return &server
}

func (this *Server) Handler(connfd net.Conn) {
	fmt.Println("连接建立成功")
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
