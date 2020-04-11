package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"golang-websocket-redis/redis"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// http升级websocket协议的配置
var wsUpgrader = websocket.Upgrader{
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 客户端连接
type wsConnection struct {
	WebUrl     string
	isSend     bool
	wsSocket   *websocket.Conn // 底层websocket
	mutex      sync.Mutex      // 避免重复关闭管道
	isClosed   bool
	closeChan  chan byte // 关闭通知
}

func WebsocketHandler(resp http.ResponseWriter, req *http.Request) {
	// 应答客户端告知升级连接为websocket
	wsSocket, err := wsUpgrader.Upgrade(resp, req, nil)
	if err != nil {
		return
	}
	wsConn := &wsConnection{
		wsSocket:  wsSocket,
		closeChan: make(chan byte),
		isClosed:  false,
	}
	//fmt.Println("已经启用WebSocket协议")
	//读取客户端信息
	go wsConn.ReadMsgLoop()
	//给客户端发信息
	go wsConn.SendMsgLoop()
}

//循环读取客户端发来的消息，传入redis.UserChan
func (wsConn *wsConnection) ReadMsgLoop() {
	var byteData []byte
	var err error
	var stringData []string
	for {
		// 读一个message
		if _, byteData, err = wsConn.wsSocket.ReadMessage(); err != nil {
			goto ERROR
		}
		stringData = strings.Split(string(byteData), "*")
		//fmt.Println("已经读取到客户端发的信息：", stringData)
		// 放入请求队列
		select {
		case redis.UserChan <- stringData:
			if !wsConn.isSend {
				wsConn.WebUrl = stringData[0]
				wsConn.isSend = true
			}
		case <-wsConn.closeChan:
			break
		}
	}
ERROR:
	wsConn.wsClose()
}

//循环发送信息给客户端,2 秒钟发送一次
func (wsConn *wsConnection) SendMsgLoop() {
	var (
		count  int64
		err    error
		data   string
	)
	for {
		if wsConn.isSend {
			break
		}
	}
	time.Sleep(time.Second*2)
	for {
			if count, err = redis.GetOnlineNum(wsConn.WebUrl); err != nil {
				fmt.Println("获取", wsConn.WebUrl, "在线人数出错:", err)
				goto ERROR
			}
			//fmt.Println("已经获取到在线人数：", count)
		    data = strconv.Itoa(int(count))
			if err = wsConn.wsSocket.WriteMessage(websocket.TextMessage, []byte(data)); err != nil {
				fmt.Println("Send to", wsConn.WebUrl, " failed:", err)
				goto ERROR
			}
			//fmt.Println("已经把在线人数发送给客户端了")
			time.Sleep(time.Second*4)
	}
ERROR:
	wsConn.wsClose()
}

func (wsConn *wsConnection) wsClose() {
	_ = wsConn.wsSocket.Close()
	//保证wsClose只被执行一次，
	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if !wsConn.isClosed {
		wsConn.isClosed = true
		close(wsConn.closeChan)
	}
	//fmt.Println("WebSocket已经关闭")
}
