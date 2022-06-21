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
