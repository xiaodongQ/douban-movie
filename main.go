// 爬取豆瓣电影 TOP250
package main

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/xiaodongQ/douban-movie/model"
	"github.com/xiaodongQ/douban-movie/parse"
)

var (
	BaseUrl = "https://movie.douban.com/top250"
)

// 新增数据
func Add(movies []parse.DoubanMovie) {
	for index, movie := range movies {
		if err := model.DB.Create(&movie).Error; err != nil {
			log.Printf("db.Create index: %d, err : %v", index, err)
		}
		log.Printf("db.Create index: %d ok", index)
	}
}

// 开始爬取
func Start() {
	var movies []parse.DoubanMovie

	pages := parse.GetPages(BaseUrl)
	i := 0
	for _, page := range pages {
		doc, err := goquery.NewDocument(strings.Join([]string{BaseUrl, page.Url}, ""))
		if err != nil {
			log.Println(err)
		}
		log.Printf("index:%d, get doc:%v", i, doc)
		i++

		movies = append(movies, parse.ParseMovies(doc)...)
	}

	log.Printf("len(movies):%d\n", len(movies))
	Add(movies)
}

func main() {
	Start()

	defer model.DB.Close()
}
