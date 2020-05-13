# 调试

* 本地部署mysql服务，新建一个mysql数据库`spiders`
	- MySQL部署及相关笔记[DB_MySQL.mds](https://github.com/xiaodongQ/devNoteBackup/blob/master/%E5%90%84%E8%AF%AD%E8%A8%80%E8%AE%B0%E5%BD%95/DB_MySQL.md)
* 修改model.go文件中，数据库相关的参数
* NewDocument里面，http请求时添加http头信息
	-  否则屏蔽爬虫会返回418状态码 &{418  418 HTTP/1.1

```golang
	client := &http.Client{}
	reqest, _ := http.NewRequest("GET", url, nil)
	reqest.Header.Add("Cookie", "bid=BFkVq2-IJhY; Expires=Thu, 13-May-21 09:45:08 GMT; Domain=.douban.com; Path=/")
	reqest.Header.Add("User-Agent", "Mozilla/5.0")
	res, e := client.Do(reqest)
```


# 爬取豆瓣电影 Top250

爬虫是标配了，第一个就从最最最简单的爬虫开始写起吧

## 目标

我们的目标站点是 [豆瓣电影 Top250](https://movie.douban.com/top250)，估计大家都很眼熟了

本次爬取8个字段，用于简单的概括分析。具体的字段如下：

![image](https://i.loli.net/2018/03/20/5ab11596b8810.png)

简单的分析一下目标源
- 一页共25条
- 含分页（共10页）且分页规则是正常的
- 每一项的数据字段排序都是规则且不变

## 开始

由于量不大，我们的爬取步骤如下
- 分析页面，获取所有的分页
- 分析页面，循环爬取所有页面的电影信息
- 爬取的电影信息入库

### 安装
```
$ go get -u github.com/PuerkitoBio/goquery
```

### 运行
```
$ go run main.go
```

### 代码片段

#### 1、获取所有分页
```
func ParsePages(doc *goquery.Document) (pages []Page) {
	pages = append(pages, Page{Page: 1, Url: ""})
	doc.Find("#content > div > div.article > div.paginator > a").Each(func(i int, s *goquery.Selection) {
		page, _ := strconv.Atoi(s.Text())
		url, _ := s.Attr("href")

		pages = append(pages, Page{
			Page: page,
			Url:  url,
		})
	})

	return pages
}
```

#### 2、分析豆瓣电影信息
```
func ParseMovies(doc *goquery.Document) (movies []Movie) {
	doc.Find("#content > div > div.article > ol > li").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".hd a span").Eq(0).Text()

		...

		movieDesc := strings.Split(DescInfo[1], "/")
		year := strings.TrimSpace(movieDesc[0])
		area := strings.TrimSpace(movieDesc[1])
		tag := strings.TrimSpace(movieDesc[2])

		star := s.Find(".bd .star .rating_num").Text()

		comment := strings.TrimSpace(s.Find(".bd .star span").Eq(3).Text())
		compile := regexp.MustCompile("[0-9]")
		comment = strings.Join(compile.FindAllString(comment, -1), "")

		quote := s.Find(".quote .inq").Text()

		...

		log.Printf("i: %d, movie: %v", i, movie)

		movies = append(movies, movie)
	})

	return movies
}
```


### 数据
![image](https://i.loli.net/2018/03/21/5ab1309594741.png)

![image](https://i.loli.net/2018/03/21/5ab131ca582f8.png)

![image](https://i.loli.net/2018/03/21/5ab130d3a00d9.png)