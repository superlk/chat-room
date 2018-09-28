package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func HttpGet(url string) (result string, err error) {
	resp, err1 := http.Get(url) //发送get请求
	if err1 != nil {
		err = err1
		return
	}
	defer resp.Body.Close()

	//读取网页内容
	buf := make([]byte, 4*1024)
	for {
		n, _ := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		result += string(buf[:n])
	}
	return
}

func SpiderPageJoy(url string) (title, content string, err error) {
	result, err1 := HttpGet(url)
	if err1 != nil {
		fmt.Println(" SpiderPageJoy httpget error =", err)
		err = err1
		return
	}
	//取关键信息
	//1， 取标题
	re1 := regexp.MustCompile(`<h1>(?s:(.*?))</h1>`)
	if re1 == nil {
		//fmt.Println(" re1  error ")
		err = fmt.Errorf("er1 error")
		return
	}
	//取关键信息
	tmpTitle := re1.FindAllStringSubmatch(result, 1) //最后一个参数为1 ，只过滤一个
	for _, data := range tmpTitle {
		title = data[1]
		title = strings.Replace(title, "\t", "", -1)
		break
	}

	//取内容<div class="content-txt pt10">段子内容<a id ="prev" href="
	re2 := regexp.MustCompile(`<div class="content-txt pt10">(?s:(.*?))<a id="prev" href="`)
	if re2 == nil {
		//fmt.Println(" re1  error ")
		err = fmt.Errorf("er2 error")
		return
	}
	tmpContent := re2.FindAllStringSubmatch(result, -1)
	for _, data := range tmpContent {
		content = data[1]
		content = strings.Replace(content, "\t", "", -1)
		content = strings.Replace(content, "\n", "", -1)
		content = strings.Replace(content, "<br/>", "", -1)
		break
	}
	return
}

func StoreJoyToFile(i int, fileTile, fileContent []string) {
	//新建文件
	f, err := os.Create(strconv.Itoa(i) + ".txt")
	if err != nil {
		fmt.Println("os.create error = ", err)
		return
	}

	defer f.Close()

	//写内容
	n := len(fileTile)
	for i := 0; i < n; i++ {
		//写标题
		f.WriteString(fileTile[i] + "\n")
		// 写内容
		f.WriteString(fileContent[i] + "\n")
		f.WriteString("\n============================\n")

	}

}

func SpiderPage(i int, page chan int) {
	// 明确要爬取的url
	url := "https://www.pengfu.com/xiaohua_" + strconv.Itoa(i) + ".html"
	fmt.Printf("正在爬取第%d个网页：%s\n", i, url)

	//开始爬取页面内容
	result, err := HttpGet(url)
	if err != nil {
		fmt.Println(" httpget error =", err)
		return
	}
	// 取<h1 class="dp-b"><a href=' 一个段子url链接"
	//正则表达式
	re := regexp.MustCompile(`<h1 class="dp-b"><a href="(?s:(.*?))"`)
	if re == nil {
		fmt.Println(" re  error ")
		return
	}
	//取关键信息
	joyUrls := re.FindAllStringSubmatch(result, -1)
	//fmt.Println(joyUrls)

	fileTile := make([]string, 0)
	fileContent := make([]string, 0)

	//取网址
	//第一个返回下标，第二个返回内容
	for _, data := range joyUrls {
		//fmt.Println("url =",data[1])
		//开始爬取每个笑话，每个段子
		title, content, err := SpiderPageJoy(data[1])
		if err != nil {
			fmt.Println("SpiderPageJoy error = ", err)
			continue
		}
		//fmt.Printf("title = #%v#\n",title)
		//fmt.Printf("content = #%v#\n",content)
		fileTile = append(fileTile, title)         //追加内容
		fileContent = append(fileContent, content) //追加内容
	}
	//fmt.Println("fileTile =", fileTile)
	//fmt.Println("fileContent =", fileContent)

	// 把内容写入到文件
	StoreJoyToFile(i, fileTile, fileContent)

	page <- i // 写内容，写第几页爬完

}

func DoWork(start, end int) {
	fmt.Printf("准备爬取第%d页到%d页的网址\n", start, end)

	page := make(chan int)

	for i := start; i <= end; i++ {
		// 定义一个函数，爬取主页面
		go SpiderPage(i, page)
	}

	for i := start; i <= end; i++ {
		fmt.Printf("第%d个页面爬完", <-page)
	}
}

func main() {
	var start, end int
	fmt.Println("请输入起始页 （>=1）")
	fmt.Scan(&start)
	fmt.Println("请输入终止页 （>=起始页）")
	fmt.Scan(&end)

	DoWork(start, end)
}
