package main

import (
	"fmt"
	"net/http"
	"github.com/PuerkitoBio/goquery"
)

func main() {
	Spider()
}

func Spider() {
	fmt.Println("爬取豆瓣电影中...")
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://movie.douban.com/top250", nil)
	if err != nil {
		fmt.Println("req err..........", err)
	}
	// req.Header.Set("Content-Type", "text/html; charset=utf-8")
	// req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	// req.Header.Set("Accept-Encoding", "gzip,deflate, br, zstd")
	// req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	// req.Header.Set("Cache-Control", "max-age=0")
	// req.Header.Set("Connection", "keep-alive")
	// req.Header.Set("Upgrade-Insecure-Requests", "1")
	// req.Header.Set("Sec-Fetch-Dest", "document")
	// req.Header.Set("Sec-Fetch-Mode", "navigate")
	// req.Header.Set("Sec-Fetch-Site", "none")
	// req.Header.Set("Sec-Fetch-User", "?1")
	// req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败..........", err)
	}
	defer resp.Body.Close()
	fmt.Println("请求成功..........")

	// 解析html
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("解析失败..........", err)
	}
	fmt.Println("解析成功..........")
	// #content > div > div.article > ol > li:nth-child(1) > div > div.info > div.hd > a > span:nth-child(1)
	doc.Find("#content > div > div.article > ol > li").Each(func(i int, s *goquery.Selection) {
		title := s.Find("div > div.info > div.hd > a > span:nth-child(1)").Text()
		img := s.Find("div > div.pic > a > img").AttrOr("src", "")
		info := s.Find("div > div.info > div.bd > p:nth-child(1)").Text()
		score := s.Find("div > div.info > div.bd > div > span.rating_num").Text()
		quote := s.Find("div > div.info > div.bd > p.quote > span").Text()

		fmt.Println("title:", title)
		fmt.Println("img:", img)
		fmt.Println("info:", info)
		fmt.Println("score:", score)
		fmt.Println("quote:", quote)
	})

	
}
