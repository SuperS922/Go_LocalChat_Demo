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
			client.PublicChat()
			break
		case 2:
			// 私聊模式
			fmt.Println("私聊模式选择。。。")
			client.PrivateChat()
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
	/* for {
		buf := make([]byte, 4096)
		client.connfd.Read(buf)
		fmt.Printf(buf)
	} */
}

func (client *Client) PublicChat() {
	// 提升用户输入消息
	var chatMsg string

	fmt.Println("----------请输入聊天内容，exit退出")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		// 发送给服务器
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.connfd.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("connfd.Write err:", err)
				break
			}
		}

		fmt.Println("----------请输入聊天内容，exit退出")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) SelectUser() {
	sendMsg := "who\n"
	_, err := client.connfd.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn write err:", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.SelectUser()
	fmt.Println("-------请输入聊天对象[用户名]：")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println("----------请输入消息内容：")
		fmt.Scanln(chatMsg)

		for chatMsg != "exit" {
			// 发送给服务器
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
				_, err := client.connfd.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("connfd.Write err:", err)
					break
				}
			}

			fmt.Println("----------请输入聊天内容，exit退出")
			fmt.Scanln(&chatMsg)
		}

		client.SelectUser()
		fmt.Println("-------请输入聊天对象[用户名]：")
		fmt.Scanln(&remoteName)
	}
}
