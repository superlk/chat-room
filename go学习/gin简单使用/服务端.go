package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io"
	"net/http"
	"os"
)

// 处理最基本的get
func func1(c *gin.Context) {
	// 回复一个200ok，在client的http-get的resp的body中获取数据
	c.String(http.StatusOK, "test1 ok")
}

//处理最基本的post
func func2(c *gin.Context) {
	// 回复一个200ok，在client的http-get的resp的body中获取数据
	c.String(http.StatusOK, "test2 ok")
}

func func3(c *gin.Context) {
	name := c.Param("name")
	passwd := c.Param("passwd")
	var result string
	result = name + passwd
	c.JSON(http.StatusOK, result)
}

func func4(c *gin.Context) {
	name := c.Param("name")
	passwd := c.Param("passwd")
	c.String(http.StatusOK, "参数:%s %s  test4 OK", name, passwd)
}

func func5(c *gin.Context) {
	name := c.Param("name")
	passwd := c.Param("passwd")
	c.String(http.StatusOK, "参数:%s %s  test4 OK", name, passwd)
}

func func6(c *gin.Context) {
	name := c.Query("name")
	passwd := c.Query("passwd")
	c.String(http.StatusOK, "参数:%s %s  test6 OK", name, passwd)
}

func func7(c *gin.Context) {
	name := c.Query("name")
	passwd := c.Query("passwd")
	var result string
	result = "name:" + name + " passwd:" + passwd
	c.JSON(http.StatusOK, result)
}

func func8(c *gin.Context) {
	message := c.PostForm("message")
	extra := c.PostForm("extra")
	nick := c.DefaultPostForm("nick", "anonymous")

	c.JSON(200, gin.H{
		"status":  "test8:posted",
		"message": message,
		"extra":   extra,
		"nick":    nick,
	})
}

// 接收client 上传文件，multipart/form-data格式的数据
// 从FormFile中获取相关的文件data!
// 然后写入本地文件
func func9(c *gin.Context) {
	file, header, err := c.Request.FormFile("uploadFile")
	filename := header.Filename
	fmt.Println(header.Filename)
	//创建临时接收文件
	out, err := os.Create("copy_" + filename)
	if err != nil {
		fmt.Println("copy error ", err)
	}
	defer out.Close()
	//copy数据
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Println("copy error 11 ", err)
	}
	c.String(http.StatusOK, "upload file success")
}

// Binding数据
// 注意:后面的form:user表示在form中这个字段是user,不是User, 同样json:user也是
// 注意:binding:"required"要求这个字段在client端发送的时候必须存在,否则报错!
type Login struct {
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// bind josn
func funcBindJSON(c *gin.Context) {
	var json Login
	// binding JSON,本质是将request中的Body中的数据按照JSON格式解析到json变量中
	if c.BindJSON(&json) == nil {
		if json.User == "TAO" && json.Password == "123" {
			c.JSON(http.StatusOK, gin.H{"JSON=== status": "you are logged in"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"JSON=== status": "unauthorized"})
		}

	} else {
		c.JSON(404, gin.H{"JSON=== status": "binding JSON error!"})
	}
}

//bind FORM
func funcBindForm(c *gin.Context) {
	var form Login
	//本质是将c中的request中的BODY数据解析到form中

	// 方法一: 对于FORM数据直接使用Bind函数, 默认使用使用form格式解析,if c.Bind(&form) == nil
	// 方法二: 使用BindWith函数,如果你明确知道数据的类型
	if c.BindWith(&form, binding.Form) == nil {
		if form.User == "TAO" && form.Password == "123" {
			c.JSON(http.StatusOK, gin.H{"FORM=== status": "you are logged in"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"FORM=== status": "unauthorized"})
		}
	} else {
		c.JSON(404, gin.H{"FORM=== status": "binding FORM error!"})
	}
}

//group test
func func10(c *gin.Context) {
	c.String(http.StatusOK, "test10 ok")
}

func func11(c *gin.Context) {
	c.String(http.StatusOK, "test11 ok")
}

func func12(c *gin.Context) {
	c.String(http.StatusOK, "test12 ok")
}

func func13(c *gin.Context) {
	c.String(http.StatusOK, "test13 ok")
}

func main() {
	// 注册一个默认的路由器
	router := gin.Default()

	// 最基本的用法
	router.GET("/test1", func1)
	router.POST("/test2", func2)

	//带参数传值，'：'必须要匹配，'*'选择匹配
	router.GET("/test3/:name/:passwd", func3)
	router.POST("/test4/:name/:passwd", func4)
	router.GET("/test5/:name/*passwd", func5)

	// 使用gin的query参数形式/test6?name=jane&passwd=666
	router.GET("/test6", func6)
	router.POST("/test7", func7)

	//参数是form中获得,即从Body中获得,忽略URL中的参数
	router.POST("/test8", func8)

	//接收上传的文件，
	router.POST("/upload", func9)

	// bind json数据
	router.POST("/bindJSON", funcBindJSON)

	//bind FORM数据
	router.POST("/bindForm", funcBindForm)

	//测试json,xml 等格式的rendering
	router.GET("/someJSON", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hey, budy", "status": http.StatusOK})
	})

	router.GET("/moreJSON", func(c *gin.Context) {
		// 注意:这里定义了tag指示在json中显示的是user不是User
		var msg struct {
			Name    string `json:"user"`
			Message string
			Number  int
		}
		msg.Name = "TAO"
		msg.Message = "hey, budy"
		msg.Number = 123
		// 下面的在client的显示是"user": "TAO",不是"User": "TAO"
		// 所以总体的显示是:{"user": "TAO", "Message": "hey, budy", "Number": 123
		c.JSON(http.StatusOK, msg)
	})

	//  测试发送XML数据
	router.GET("/someXML", func(c *gin.Context) {
		c.XML(http.StatusOK, gin.H{"name": "TAO", "message": "hey, budy", "status": http.StatusOK})
	})

	//group
	group1 := router.Group("/g1")
	group1.GET("/read1", func10)
	group1.GET("/read2", func11)

	group2 := router.Group("/g2")
	group2.POST("/write1", func12)
	group2.POST("/write2", func13)

	// 下面测试静态文件服务
	// 显示当前文件夹下的所有文件/或者指定文件
	router.StaticFS("/showDir", http.Dir("."))
	router.Static("/files", "/bin")
	router.StaticFile("/image", "demo.txt")

	//加载模版template
	//router.LoadHTMLGlob("templates/*")
	router.LoadHTMLFiles("templates/index.html", "templates/template2.html")
	// 或者使用这种方法加载也是OK的: router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "indxe.html", gin.H{
			"title": "GIN:test html template",
		})
	})

	//绑定端口是8888
	router.Run(":8888")
}
