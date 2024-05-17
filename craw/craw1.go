// github.com/Guohuixixi/crawDemo/craw
package craw

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	// 正则查找图片的地址
	reImg = `http://i.52desktop.cn:81/upimg/allimg/[^"]+?(\.((jpg)|(png)|(jpeg)|(gif)|(bmp)))`
)

// 读取链接中所有的内容
func GetPageStr(url string) (pageStr string) {

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("http.Get error", err)
		return
	}
	defer resp.Body.Close()
	// 读取页面内容
	pagebytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll error", err)
		return
	}
	// 字节转换为字符串
	pageStr = string(pagebytes)
	return pageStr
}

// 下载图片，核心内容是获取图片的reader、write对象
func downloadImg(savePath string, url string) error {

	v, err := http.Get(url)
	if err != nil {
		fmt.Println("http.Get error", err)
		return err
	}
	defer v.Body.Close()
	fileName := path.Base(url) // 获取文件名
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(v.Body, 32*1024)
	// 创建单个文件保存下来
	file, err := os.Create(savePath + fileName)
	if err != nil {
		fmt.Println("ioutil.WriteFile error", err)
		return err
	}
	// // 获得文件的writer对象
	writer := bufio.NewWriter(file)

	written, err := io.Copy(writer, reader)
	if err != nil {
		fmt.Println("io.Copy error", err)
		return err
	}
	fmt.Println(fileName, "download is success, Total length:", written)
	return nil

}

// 给main.go进行调用的外部地址
func GetCrawl1Img() {
	start := time.Now()
	for i := 2; i <= 21; i++ {
		// 根据观察规则，发现下载图片地址网站
		url := "http://www.52desktop.cn/html/DLZM/KPBZ/20191205/15898_" + strconv.Itoa(i) + ".html"
		pageStr := GetPageStr(url)
		re := regexp.MustCompile(reImg)
		results := re.FindAllStringSubmatch(pageStr, -1) //  -1 表示搜索所有可能的匹配项。
		for _, result := range results {
			if strings.Contains(result[0], "20191204") {
				// 过滤条件读取图片地址
				downloadImg("D:\\school\\internShipTC\\day 2\\crawDemo\\img1\\", result[0])

			}

		}

	}
	elapse := time.Since(start)
	fmt.Println("elapsed time is : ", elapse, "s")

}
