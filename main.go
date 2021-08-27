package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/anaskhan96/soup"
)

func main() {
	var start, end int
	fmt.Print("请输入起始页（>=1）：")
	fmt.Scan(&start)
	fmt.Print("请输入结束页（>=起始页）：")
	fmt.Scan(&end)

	//开始工作
	DoWork(start, end)
}

func DoWork(start, end int) {
	var fileName string
	fmt.Println("请输入爬取结果保存的文件名（例如：“result.txt”）")
	fmt.Scan(&fileName)
	fileName = "./" + fileName
	fileObj, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		panic(fmt.Errorf("open file %s failed, %v", fileName, err))
	}

	fmt.Printf("正在爬取页数范围是第%d页到第%d页。\n", start, end)
	page := make(chan int)
	for i := start; i <= end; i++ {
		//爬取主网页
		go SpiderPage(i, page, fileObj)
	}
	for i := start; i <= end; i++ {
		fmt.Printf("第%d页已经爬取完成。\n", <-page)
	}
	fileObj.Close() // 关闭文件
}

func SpiderPage(i int, page chan int, file *os.File) {
	// 明确爬取的url
	//http://www.downcc.com/soft/list_181_1.html
	requestUrl := "http://www.downcc.com/soft/list_181_" + strconv.Itoa(i) + ".html"
	// fmt.Printf("正在爬取第%d个网页：%s\n", i, requestUrl)

	// 开始爬取网页的内容
	doc := fetch(requestUrl)
	links := doc.Find("ul", "id", "li-change-color").FindAll("h3", "class", "soft-ht1")

	for _, link := range links {
		softUrl, _ := link.Find("a").Attrs()["href"]
		softName := link.Find("a", "class", "mg-r10")
		// fmt.Println("http://www.downcc.com"+softUrl, softName.Text())  终端输出爬取结果
		// 将结果写入到文件中
		// 多协程写是没问题的，因为go的标准库最终写文件的时候，会用读写锁
		// 一个goroutine独占文件句柄的 没有写完之前是不会让出文件句柄的, 所以不会错乱
		crawlResult := fmt.Sprintf("[%s]\t\t http://www.downcc.com/%s\n", softName.Text(), softUrl) 
		fmt.Fprintf(file, crawlResult) //往文件里写文件
	}

	page <- i
}

func fetch(url string) soup.Root {
	fmt.Println("Fetch Url", url)
	soup.Headers = map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36",
	}
	source, err := soup.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	doc := soup.HTMLParse(source)
	return doc
}
