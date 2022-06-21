package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	connfd     net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       -1,
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

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1. 公聊模式")
	fmt.Println("2. 私聊模式")
	fmt.Println("3. 更新用户名")
	fmt.Println("0. 退出")

	fmt.Scan(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("-----请输入合法范围的数字-----")
		return false
	}
}

func (client *Client) UpdateName() bool {
	fmt.Println("请输入用户名：")
	fmt.Scan(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.connfd.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("connfd.Write error", err)
		return false
	}

	return true
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}

		switch client.flag {
		case 1:
			// 公聊模式
			fmt.Println("公聊模式选择。。。")
			break
		case 2:
			// 私聊模式
			fmt.Println("私聊模式选择。。。")
			break
		case 3:
			// 更新用户名
			fmt.Println("更新用户名。。。")
			client.UpdateName()
			break
		case 0:
			// 退出
			break
		}
	}
}

func (client *Client) DealResponse() {
	// 与下面的 for 是相同的作用
	// 永久阻塞监听的 IO方式   同步IO
	io.Copy(os.Stdout, client.connfd) // 输出重定向吗？

	// for {
	// 	buf := make([]byte, 4096)
	// 	client.connfd.Read(buf)
	// 	fmt.Printf(buf)
	// }
}
