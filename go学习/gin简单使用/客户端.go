package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func helpRead(resp *http.Response) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR2!: ", err)
	}
	fmt.Println(string(body))
}

//上传文件
func Uplaodflie() {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, _ := w.CreateFormFile("uploadFile", "demo.txt")
	fd, _ := os.Open("1.txt")
	fmt.Println("fd ==", fd)

	defer fd.Close()
	io.Copy(fw, fd)
	w.Close()
	resp, _ := http.Post("http://0.0.0.0:8888/upload", w.FormDataContentType(), buf)
	helpRead(resp)

}

func main() {
	//调用get，并获取返回值
	resp, _ := http.Get("http://0.0.0.0:8888/test1")
	helpRead(resp)

	resp, _ = http.Post("http://0.0.0.0:8888/test2", "", strings.NewReader(""))
	helpRead(resp)

	// GET传参数,使用gin的Param解析格式: /test3/:name/:passwd
	resp, _ = http.Get("http://0.0.0.0:8888/test3/name=TAO/passwd=123")
	helpRead(resp)

	resp, _ = http.Post("http://0.0.0.0:8888/test4/name=TAO/passwd=123", "", strings.NewReader(""))
	helpRead(resp)

	// 注意Param中':'和'*'的区别
	resp, _ = http.Get("http://0.0.0.0:8888/test5/name=TAO/")
	helpRead(resp)

	resp, _ = http.Get("http://0.0.0.0:8888/test6?name=TAO&passwd=123")
	helpRead(resp)

	resp, _ = http.Post("http://0.0.0.0:8888/test7?name=TAO&passwd=123", "", strings.NewReader(""))
	helpRead(resp)

	resp, _ = http.Post("http://0.0.0.0:8888/test8", "application/x-www-form-urlencoded", strings.NewReader("message=8888888&extra=999999"))
	helpRead(resp)

	//上传文件
	Uplaodflie()

	// 下面测试bind
	resp, _ = http.Post("http://0.0.0.0:8888/bindJSON", "application/json", strings.NewReader("{\"user\":\"TAO\", \"password\": \"123\"}"))
	helpRead(resp)

	// 下面测试bind FORM数据
	resp, _ = http.Post("http://0.0.0.0:8888/bindForm", "application/x-www-form-urlencoded", strings.NewReader("user=TAO&password=123"))
	helpRead(resp)

	// 下面测试接收JSON和XML数据
	resp, _ = http.Get("http://0.0.0.0:8888/someJSON")
	helpRead(resp)
	resp, _ = http.Get("http://0.0.0.0:8888/moreJSON")
	helpRead(resp)
	resp, _ = http.Get("http://0.0.0.0:8888/someXML")
	helpRead(resp)

	// 下面测试router 的GROUP
	resp, _ = http.Get("http://0.0.0.0:8888/g1/read1")
	helpRead(resp)
	resp, _ = http.Get("http://0.0.0.0:8888/g1/read2")
	helpRead(resp)
	resp, _ = http.Post("http://0.0.0.0:8888/g2/write1", "", strings.NewReader(""))
	helpRead(resp)
	resp, _ = http.Post("http://0.0.0.0:8888/g2/write2", "", strings.NewReader(""))
	helpRead(resp)

	// 测试加载HTML模板
	//resp,_ = http.Get("http://0.0.0.0:8888/index")
	//helpRead(resp)
}
