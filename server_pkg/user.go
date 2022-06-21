package main

import (
	"net"
	"strings"
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
	// 用户下线
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	// 广播当前用户下线消息
	this.server.BroadCast(this, "is offline.")
}

// 给当前User对应的客户端发送消息
func (this *User) SendMsg(msg string) {
	this.Connfd.Write([]byte(msg))
}

func (this *User) rename(msg string) {
	newName := strings.Split(msg, "|")[1]
	// 判断 newName 是否已经存在
	_, ok := this.server.OnlineMap[newName]
	if ok {
		this.SendMsg("用户名已被占用!\n")
	} else {
		this.server.mapLock.Lock()
		delete(this.server.OnlineMap, this.Name)
		this.server.OnlineMap[newName] = this
		this.server.mapLock.Unlock()

		this.Name = newName
		this.SendMsg("更改成功：" + this.Name + "\n")
	}
}

func (this *User) sendPrivateMsg(msg string) {
	// 1. 获取对方的用户名
	sliceMsg := strings.Split(msg, "|")
	remoteName := sliceMsg[1]
	if remoteName == "" {
		this.SendMsg("消息格式不正确，请使用\"to|张三|消息内容\"格式。\n")
		return
	}

	// 2. 根据用户名获得 User 对象
	remoteUser, ok := this.server.OnlineMap[remoteName]
	if !ok {
		this.SendMsg(remoteName + "is offline.\n")
		return
	}

	// 3. 获取消息内容并且发送
	content := sliceMsg[2]
	if content == "" {
		this.SendMsg("无消息内容\n")
		return
	}

	remoteUser.SendMsg("[" + this.Name + "]:" + content + "\n")
}

// 用户处理消息的业务
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		// 查询当前在线的用户d
		this.server.mapLock.RLock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Name + "]" + ":" + "is online.\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.RUnlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		this.rename(msg)
	} else if len(msg) > 4 && msg[:3] == "to|" {
		this.sendPrivateMsg(msg)
	} else {
		this.server.BroadCast(this, msg)
	}

}

// 监听当前User channel 的方法，一旦有消息，就直接发送给对端客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		this.Connfd.Write([]byte(msg + "\n"))
	}
}
