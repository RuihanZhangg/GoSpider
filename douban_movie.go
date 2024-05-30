/*
爬取豆瓣top250电影的信息并存入mysql数据库
静态数据爬取
使用协程来并行，每一页（共10页）使用一个goroutine爬取，从1.6s提升到0.6s。可以根据数据库表中电影的顺序来判断是否是并行的。
但是如果在每一页中，每个电影都使用一个goroutine来爬取，会导致goroutine数量过多，反而会变慢为1.1s。这是应为goroutine的创建和销毁也是需要时间的，而这些 goroutine 的运行时间都很短，反而降低了效率
*/
package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	USERNAME = "root"
	PASSWORD = "your_password"
	HOST     = "127.0.0.1"
	PORT     = "3306"
	DBNAME   = "douban_movie"
)

var DB *sql.DB

type Movie struct {
	Title    string
	Img      string
	Director string
	Actor    string
	Year     string
	Score    string
	Quote    string
}

func main() {
	if InitDB() {
		ch := make(chan int, 10)
		for i := 0; i < 10; i++ {
			go Spider(ch, strconv.Itoa(i*25), i)
		}
		for i := 0; i < cap(ch); i++ {
			v := <-ch
			fmt.Println("爬取完成..........", v)
		}
	}
}

func Spider(ch chan int, p string, i int) {
	fmt.Println("爬取豆瓣TOP250电影中...")
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://movie.douban.com/top250?start="+p, nil)
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
	
		director, actor, year := InfoSpite(info)
		var movie = Movie{title, img, director, actor, year, score, quote}
		_ = movie
		if InsertDB(movie){
			fmt.Println("插入成功..........")
		}
	})
	ch <- i
}


func InfoSpite(info string) (director string, actor string, year string) {
	directorReg, _ := regexp.Compile(`导演: (.*)主演: `)
	director = string(directorReg.Find([]byte(info)))
	actorReg, _ := regexp.Compile(`主演: (.*)`)
	actor = string(actorReg.Find([]byte(info)))
	yearReg, _ := regexp.Compile(`(\d+)`)
	year = string(yearReg.Find([]byte(info)))
	return
}

func InitDB() bool {
	var err error
	DB, err = sql.Open("mysql", USERNAME+":"+PASSWORD+"@tcp("+HOST+":"+PORT+")/"+DBNAME)
	if err != nil {
		fmt.Println("打开数据库失败..........", err)
		return false
	}
	DB.SetConnMaxLifetime(10)
	DB.SetMaxIdleConns(5)
	err = DB.Ping()
	if err != nil {
		fmt.Println("连接数据库失败..........", err)
		return false
	}
	fmt.Println("连接数据库成功..........")
	return true
}

func InsertDB(movie Movie) bool {
	stmt, err := DB.Prepare("insert into movie(title, img, director, actor, year, score, quote) values(?,?,?,?,?,?,?)")
	if err != nil {
		fmt.Println("插入数据失败..........", err)
		return false
	}
	_, err = stmt.Exec(movie.Title, movie.Img, movie.Director, movie.Actor, movie.Year, movie.Score, movie.Quote)
	if err != nil {
		fmt.Println("插入数据失败..........", err)
		return false
	}
	// fmt.Println("插入数据成功..........")
	return true
}
