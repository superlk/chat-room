package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type Clinet struct {
	C    chan string // 用户发送数据的管道
	Name string      // 用户名
	Addr string      //网络地址
}

// 保存在线用户，cliAddr ==== > Client
var onlineMap map[string]Clinet

var message = make(chan string)

//新开一个携程，转发消息，只要有消息来了，遍历map，给mao每个成员发送这个消息
func Manager() {
	// 给map 新建，分配空间
	onlineMap = make(map[string]Clinet)
	for {
		msg := <-message //没有消息前，这里会堵塞
		//遍历map，给mao每个成员发送这个消息
		for _, cli := range onlineMap {
			cli.C <- msg
		}
	}

}

func writeMagToClient(cli Clinet, conn net.Conn) {
	for msg := range cli.C { //给当前客户发送信息
		conn.Write([]byte(msg + "\n"))
	}

}

func MakeMsg(cli Clinet, msg string) (buf string) {
	buf = "[" + cli.Addr + "]" + cli.Name + ":" + msg
	return
}

func HandleConn(conn net.Conn) { // 处理用户链接
	//	获取客户端的网络地址
	cliAddr := conn.RemoteAddr().String()

	// 创建一个结构体,默认，用户名和网络地址一样
	cli := Clinet{make(chan string), cliAddr, cliAddr}
	// 把结构体添加到map
	onlineMap[cliAddr] = cli

	// 重新开一个携程，给当前客户端发送信息
	go writeMagToClient(cli, conn)
	// 广播某个在线
	//message <- "[" + cli.Addr + "]" + cli.Name + ": login"
	message <- MakeMsg(cli, "login")
	//提示我是谁
	cli.C <- MakeMsg(cli, "i am here")

	isQuit := make(chan bool)  // 对方是否主动退出
	hasData := make(chan bool) //对方是否有数据发送

	// 新建一个携程，接收用户发送过来的的数据
	go func() {
		buf := make([]byte, 2048)
		for {
			n, err := conn.Read(buf)
			if n == 0 { // 两种情况，对方断开，或出现问题
				isQuit <- true //主动退出
				fmt.Println("conn read error =", err)
				return
			}
			msg := string(buf[:n])

			if len(msg) == 3 && msg == "who" {
				// 变量map,发送当前用户发送所有当前成员
				conn.Write([]byte("user list:\n"))
				for _, tmp := range onlineMap {
					msg = tmp.Addr + ":" + tmp.Name + "\n"
					conn.Write([]byte(msg))
				}

			} else if len(msg) >= 8 && msg[:6] == "rename" {
				//renme|mike
				name := strings.Split(msg, "|")[1]
				cli.Name = name
				onlineMap[cliAddr] = cli
				conn.Write([]byte("rename ok"))

			} else {
				message <- MakeMsg(cli, msg)

			}

			hasData <- true //代表有数据
			// 转发消息
			//message <- MakeMsg(cli, msg)
		}
	}()
	for {
		//通过select检测channel的流动
		select {
		case <-isQuit:
			delete(onlineMap, cliAddr)         //从map里移除当前用户
			message <- MakeMsg(cli, "logout ") //广播谁下线
			return
		case <-hasData:
		case <-time.After(60 * time.Second): // 60s后移除
			delete(onlineMap, cliAddr)
			message <- MakeMsg(cli, "time out leave out")
			return
		}
	}
}

func main() {
	// 监听
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println(" listen error ==", err)
		return
	}

	defer listener.Close()

	//新开一个携程，转发消息，只要有消息来了，遍历map，给map每个成员发送这个消息
	go Manager()

	// 主携程，循环堵塞等待用户链接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(" listen accept error = ", err)
			continue
		}
		go HandleConn(conn) //处理用户链接
	}
}
