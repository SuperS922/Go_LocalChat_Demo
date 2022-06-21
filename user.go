package main

import (
	"net"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	Connfd net.Conn
}

func NewUser(connfd net.Conn) *User {
	userAddr := connfd.RemoteAddr().String()

	user := User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		Connfd: connfd,
	}

	// 启动监听当前user channel消息的goroutine
	go user.ListenMessage()

	return &user
}

// 监听当前User channel 的方法，一旦有消息，就直接发送给对端客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		this.Connfd.Write([]byte(msg + "\n"))
	}
}
