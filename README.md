#golang-websocket-redis
##使用方法：
> 1. 将js置于网页```<body>```标签内；
> 2. 配置```Redis```地址及密码；
> 3. 配置```Gin```路由；
> 4. 启动```golang```即可；

##Document
> 1. 客户端链接上websocket服务，服务端启动WebsocketHandler服务，5秒后将发送本机ID给到服务端；
> 2. 服务端接受到信息后，将存入Redis有序集合中，集合名为用户所浏览网页的url；
> 3. 服务端接受到信息后的延迟2秒后，开始给客户端发送10秒内的相应用户；
> 4. 客户端接受到信息后，将在body底部创建一个div节点，用以展现实时UV数据；
> 5. 当客户端关闭以后，服务端关闭websocket链接，结束WebsocketHandler；
> 6. 至此，结束；
> 7. 很简单的一套UV、PV监控服务，欢迎一起交流学习。
