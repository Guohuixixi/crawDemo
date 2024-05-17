// github.com/Guohuixixi/crawDemo/craw
package craw

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	// 正则提取表达式
	geReImg = `http://i.52desktop.cn:81/upimg/allimg/[^"]+?(\.((jpg)|(png)|(jpeg)|(gif)|(bmp)))`

	// 存放图片链接的数据管道
	chanImageUrls chan string
	// 启动多个goroutine，sync.WaitGroup来实现goroutine的同步
	waitGroup sync.WaitGroup
	// 用于监控协程
	chanTask chan string
	// 监控协程总数
	chanTaskCount = 23
	// 监控当前的下载任务是否完成的总数
	count int
)

// 爬取图片链接，获取当前页的所有图片链接
func getImgs(url string) (urls []string) {
	pageStr := GetPageStr(url)
	re := regexp.MustCompile(geReImg)
	results := re.FindAllStringSubmatch(pageStr, -1) //  -1 表示搜索所有可能的匹配项。
	fmt.Printf("共找到%d条结果\n", len(results))
	for _, result := range results {
		url := result[0]
		if strings.Contains(result[0], "20191204") { // 过滤条件读取图片地址
			urls = append(urls, url)
			// fmt.Println(url)

		}
	}
	return urls

}

func getImgUrls(url string) {
	urls := getImgs(url)
	// 遍历切片里所有链接，存入数据管道
	for _, url := range urls {
		// 存放图片链接的数据管道，发送图片链接
		chanImageUrls <- url
	}
	// 标识当前协程完成
	// 每完成一个任务，写一条数据
	// 用于监控协程知道已经完成了几个任务
	chanTask <- url
	waitGroup.Done()

}

// 任务统计
func CheckOk() {
	for {
		// 从chanTask中接收值并赋值给变量,url
		url := <-chanTask
		fmt.Printf("%s 完成了爬虫任务\n", url)
		count++
		// 任务统计个协程是否全部完成,完成了久关闭通道
		// 因为的我循环起始地址的索引为2，所以这里减去2
		if count == chanTaskCount-2 {
			close(chanImageUrls)
			break
		}
	}
	waitGroup.Done()

}
func downloadFile(savePath string, url string) {

	v, err := http.Get(url)
	if err != nil {
		fmt.Println("http.Get error", err)
		return
	}
	defer v.Body.Close()
	fileName := savePath + path.Base(url) // 获取文件名
	bytes, err := ioutil.ReadAll(v.Body)
	if err != nil {
		fmt.Println("io.ReadAll error", err)
		return
	}
	// 直接写入内容，比Crawl1.go代码更简单
	err = ioutil.WriteFile(fileName, bytes, 0666)
	if err != nil {
		fmt.Println("ioutil.WriteFile error", err)
		return
	}
	fmt.Println(fileName, " download is success")

}

// 下载协程，从chanImageUrls管道中读取链接下载,读取的数据来自于这里函数getImgUrls，对chanImageUrls的写入
func DownloadImg() {
	for url := range chanImageUrls {
		downloadFile("D:\\school\\internShipTC\\day 2\\crawDemo\\img2\\", url)
	}
	waitGroup.Done()
}

func GetCrawl2Img() {
	// 初始化管道
	chanImageUrls = make(chan string, 100)
	chanTask = make(chan string, chanTaskCount)
	start := time.Now()

	for i := 2; i < chanTaskCount; i++ {
		waitGroup.Add(1)
		// 根据观察规则，发现下载图片地址网站
		go getImgUrls("http://www.52desktop.cn/html/DLZM/KPBZ/20191205/15898_" + strconv.Itoa(i) + ".html")

	}
	// 任务统计个协程是否全部完成
	waitGroup.Add(1)
	go CheckOk()
	// 创建一个新的协程来执行 Wait 后面的操作
	//go func() {
	//	// 等待所有协程完成
	//	waitGroup.Wait()
	//
	//	// 执行一些操作
	//
	//}()
	// 下载协程，从管道中读取链接下载
	for i := 0; i < 5; i++ {
		waitGroup.Add(1)
		go DownloadImg()

	}

	waitGroup.Wait()
	fmt.Println("All workers have completed.")
	elapse := time.Since(start)
	fmt.Println("elapsed time is : ", elapse, "s")

}
