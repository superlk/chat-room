# chat-room
## 基于go的并发聊天室，

# 启动server端
## go run chat-room/src/main.go

# 测试
## 使用telent测试

bogon:~ superlk$ telnet 127.0.0.1 8000

Trying 127.0.0.1...

Connected to localhost.

Escape character is '^]'.

[127.0.0.1:55251]127.0.0.1:55251:i am here

[127.0.0.1:55251]127.0.0.1:55251:login

ni hao

[127.0.0.1:55251]127.0.0.1:55251:ni hao


[127.0.0.1:55254]127.0.0.1:55254:login

[127.0.0.1:55254]127.0.0.1:55254:hi

[127.0.0.1:55254]127.0.0.1:55254:ni hao

ni hao

[127.0.0.1:55251]127.0.0.1:55251:ni hao

[127.0.0.1:55254]127.0.0.1:55254:88

88

[127.0.0.1:55251]127.0.0.1:55251:88

[127.0.0.1:55254]127.0.0.1:55254:logout

# 百度贴吧并发爬虫.go


# 段子并发爬虫并保存到文件.go


请输入起始页 （>=1）

1

请输入终止页 （>=起始页）

10

准备爬取第1页到10页的网址

正在爬取第10个网页：https://www.pengfu.com/xiaohua_10.html

正在爬取第5个网页：https://www.pengfu.com/xiaohua_5.html

正在爬取第1个网页：https://www.pengfu.com/xiaohua_1.html

正在爬取第3个网页：https://www.pengfu.com/xiaohua_3.html

正在爬取第4个网页：https://www.pengfu.com/xiaohua_4.html

正在爬取第6个网页：https://www.pengfu.com/xiaohua_6.html

正在爬取第2个网页：https://www.pengfu.com/xiaohua_2.html

正在爬取第9个网页：https://www.pengfu.com/xiaohua_9.html

正在爬取第7个网页：https://www.pengfu.com/xiaohua_7.html

正在爬取第8个网页：https://www.pengfu.com/xiaohua_8.html

第10个页面爬完第8个页面爬完第1个页面爬完第6个页面爬完第9个页面爬完第4个页面爬完第7个页面爬完第2个页面爬完第3个页面爬完第5个页面爬完

Process finished with exit code 0



