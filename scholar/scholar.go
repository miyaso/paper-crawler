package scholar

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
	"time"
	"strconv"
)

type Article struct {
	Title      string
	Url        string
	PdfUrl     string
	Year       string
	Author     []string
	BookTitle  string
	Journal    string
	CitedBy    []string
	CitedNum   int
	Versions   []string
	VersionNum int
	ClusterId  string
	BibTex     string
}
type ArticleSlice struct {
	Articles []Article
}

var (
	baseUrl        = "https://scholar.google.co.jp"
	reTitle, _     = regexp.Compile("title={(.+)},")
	reAuthor, _    = regexp.Compile("author={(.+)},")
	reBookTitle, _ = regexp.Compile("booktitle={(.+)},")
	reYear, _      = regexp.Compile("year={(.+)},")
	reJournal, _   = regexp.Compile("journal={(.+)},")
	reClusterId, _ = regexp.Compile("cluster=([0-9]+)")
)

func ExtractStringForReg(str string, reg *regexp.Regexp) string {
	findList := reg.FindStringSubmatch(str)
	if len(findList) < 2 {
		return ""
	}
	return findList[1]
}

func ScrapCiteIds(citedUrl string, reqCiteNum int, maxRequestCitedNum int) []string {
	citeIds := []string{}
	doc, _ := goquery.NewDocument(baseUrl + citedUrl + "&num=20&start=" + string(reqCiteNum))
	doc.Find("div.gs_r").Each(func(_ int, artSelect *goquery.Selection) {
		if reqCiteNum >= maxRequestCitedNum {
			return
		}
		// Citations List
		s := artSelect.Find("div.gs_fl>a")
		versionUrl, _ := s.Eq(2).Attr("href")
		clusterId := ExtractStringForReg(versionUrl, reClusterId)
		citeIds = append(citeIds, clusterId)
		reqCiteNum++
		return
	})
	return citeIds
}

func ScrapArticle(url string, reqNum int, maxRequestNum int, maxRequestCitedNum int, crawlInterval int) []Article {
	articles := []Article{}
	doc, _ := goquery.NewDocument(url + "&start=" + string(reqNum))
	doc.Find("div.gs_r").Each(func(_ int, artSelect *goquery.Selection) {
		if reqNum >= maxRequestNum {
			return
		}
		art := Article{}

		art.PdfUrl, _ = artSelect.Find("div.gs_md_wp>a").Attr("href")
		art.Url, _ = artSelect.Find("h3.gs_rt>a").First().Attr("href")

		// Citations List
		s := artSelect.Find("div.gs_fl>a")

		if citedStrSplit := strings.Split(s.Eq(0).Text(), " "); len(citedStrSplit) > 2 {
			if citedNum, err := strconv.Atoi(citedStrSplit[2]); err == nil {
				art.CitedNum = citedNum
			}
		}
		citedUrl, _ := s.Eq(0).Attr("href")

		reqCiteNum := 0
		if art.CitedNum > maxRequestCitedNum {
			maxRequestCitedNum = art.CitedNum
		}
		for {
			time.Sleep(time.Duration(crawlInterval) * time.Second)
			citeIds := ScrapCiteIds(citedUrl, reqCiteNum, maxRequestCitedNum)
			reqCiteNum += len(citeIds)
			art.CitedBy = append(art.CitedBy, citeIds...)
			if len(citeIds) < 20 || reqCiteNum >= maxRequestCitedNum {
				break
			}
		}
		relatedUrl, _ := s.Eq(1).Attr("href")
		info := strings.Split(relatedUrl, ":")[1]
		versionUrl, _ := s.Eq(2).Attr("href")
		if versionStrSplit := strings.Split(s.Eq(2).Text(), " "); len(versionStrSplit) > 2 {
			if versionNum, err := strconv.Atoi(versionStrSplit[2]); err == nil {
				art.VersionNum = versionNum
			}
		}
		art.ClusterId = ExtractStringForReg(versionUrl, reClusterId)

		time.Sleep(time.Duration(crawlInterval) * time.Second)
		cite, _ := goquery.NewDocument(baseUrl + "/scholar?q=info:" + info + ":scholar.google.com/&output=cite&scirp=0&hl=en")
		bibtexUrl, _ := cite.Find("a.gs_citi").First().Attr("href")

		fmt.Println(bibtexUrl)
		time.Sleep(time.Duration(crawlInterval) * time.Second)
		bibTex, _ := goquery.NewDocument(baseUrl + bibtexUrl)

		art.BibTex = bibTex.Text()
		art.Title = ExtractStringForReg(bibTex.Text(), reTitle)
		art.Author = strings.Split(ExtractStringForReg(bibTex.Text(), reAuthor), ",")
		art.BookTitle = ExtractStringForReg(bibTex.Text(), reBookTitle)
		art.Journal = ExtractStringForReg(bibTex.Text(), reJournal)
		art.Year = ExtractStringForReg(bibTex.Text(), reYear)

		fmt.Println(url)
		fmt.Println()

		articles = append(articles, art)
		reqNum++
		return
	})
	return articles
}

func CrawlScholar(query string, maxRequestNum int, maxRequestCitedNum int, crawlInterval int) string {
	url := "https://scholar.google.co.jp/scholar?q="+query+"&btnG=&hl=en&as_sdt=0%2C5&num=20"
	fmt.Println(url)
	fmt.Println(maxRequestNum)
	articleSlice := ArticleSlice{}

	reqNum := 0
	for {
		articles := ScrapArticle(url, reqNum, maxRequestNum, maxRequestCitedNum, crawlInterval)
		reqNum += len(articles)
		articleSlice.Articles = append(articleSlice.Articles, articles...)
		if len(articles) < 20 || reqNum >= maxRequestNum {
			break
		}
		time.Sleep(time.Duration(crawlInterval) * time.Second)
	}
	// export to json
	b, err := json.Marshal(articleSlice)
	if err != nil {
		fmt.Println("json err:", err)
	}
	//fmt.Println(string(b))
	return string(b)
}

/*
func main() {
	query := "LDA"
	url := "https://scholar.google.co.jp/scholar?q=" + query + "&btnG=&hl=en&as_sdt=0%2C5&num=20"
	maxRequestNum := 10
	maxRequestCitedNum := 10
	crawlInterval := 3
	CrawlScholar(url, maxRequestNum, maxRequestCitedNum, crawlInterval)
}
*/
