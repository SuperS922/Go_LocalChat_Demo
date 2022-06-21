package main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	connfd     net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	client := Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}

	// 链接server
	connfd, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.ServerIp, client.ServerPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.connfd = connfd

	return &client
}

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println("链接服务器失败")
		return
	} else {
		fmt.Println("链接成功")
	}

	// 此处阻塞
	select {}
}
