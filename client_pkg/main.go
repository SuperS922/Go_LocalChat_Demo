package main

import (
	"flag"
	"fmt"
)

var serverIp string
var serverPort int

// init 函数是在 main 函数之前执行的
//  ./client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认是127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器Port(默认为8888)")
}

func main() {
	// 命令行解析
	flag.Parse() // 会完成命令行的解析

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("链接服务器失败")
		return
	} else {
		fmt.Println("链接成功")
	}

	client.Run()
}
