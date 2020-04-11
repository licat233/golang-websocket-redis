package redis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

var (
	UserChan = make(chan []string, 5000)
	err      error
	RDconn   redis.Conn
)

//连接Redis数据库
func init() {
	RDconn, err = redis.Dial("tcp", "127.0.0.1:6379", redis.DialPassword(""))
	if err != nil {
		log.Fatal("Redis链接失败")
	}
	//fmt.Println("已连接Redis数据库")
}

//统计10秒内的在线人数,这个是非永驻进程
func GetOnlineNum(weburl string) (int64, error) {
	var nowtimestamp = int32(time.Now().Unix())
	return redis.Int64(RDconn.Do("ZCOUNT", weburl, nowtimestamp-10, nowtimestamp))
}

//将客户端发来的消息写入Redis,这个是永驻进程
func SetMsgLoop() {
	var data []string
	for {
		//放入请求队列
		select {
		case data = <-UserChan:
			//fmt.Println("Redis已经接受到客户端传来的data：", data)
			//键值有效期 10 秒
			if _, err := redis.Bool(RDconn.Do("ZADD", data[0],int32(time.Now().Unix()),data[1])); err != nil {
				fmt.Println("向", data[0], "集合添加成员失败:", err)
			}
		}
	}
}
