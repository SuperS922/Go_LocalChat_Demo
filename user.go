package main

import (
	"net"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	Connfd net.Conn
	server *Server
}

func NewUser(connfd net.Conn, server *Server) *User {
	userAddr := connfd.RemoteAddr().String()

	user := User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		Connfd: connfd,
		server: server,
	}

	// 启动监听当前user channel消息的goroutine
	go user.ListenMessage()

	return &user
}

// 用户的上线业务
func (this *User) Online() {
	// 用户上线，将用户加入 OnlineMap 中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	// 广播当前用户上线消息
	this.server.BroadCast(this, "is online.")
}

// 用户的下线业务
func (this *User) Offline() {
	// 用户下线，将用户加入 OnlineMap 中
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	// 广播当前用户下线消息
	this.server.BroadCast(this, "is offline.")
}

// 用户处理消息的业务
func (this *User) DoMessage(msg string) {
	this.server.BroadCast(this, msg)
}

// 监听当前User channel 的方法，一旦有消息，就直接发送给对端客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		this.Connfd.Write([]byte(msg + "\n"))
	}
}
