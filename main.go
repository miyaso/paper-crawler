package main

import (
	"fmt"
	"os"
	"github.com/codegangsta/cli"
	"strings"
	"papper-crawler/scholar"
)

func main() {
	app := cli.NewApp()
	app.Name = "crawl_google_scholar_client"
	app.Usage = ""
	app.Version = "0.0.1"

	// グローバルオプション設定
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "dryrun, d", // 省略指定 => d
			Usage: "dryrun",
		},

	}

	app.Commands = []cli.Command{
		// コマンド設定
		{
			Name:    "get",
			Aliases: []string{"g"},
			Usage:   "Crawling google scholar.",
			Action:  getAction,
			Flags:  []cli.Flag{
				cli.IntFlag{
					Name: "maxRequestNum",
					Value: 10,
					Usage: "The maximum number of crawling artciles",
				},
				cli.IntFlag{
					Name: "maxRequestCitedNum",
					Value: 10,
					Usage: "The maximum number of articleId citing by crawling artciles ",
				},
				cli.IntFlag{
					Name: "crawlInterval",
					Value: 3,
					Usage: "Interval seconds for crawling. ",
				},
			},
		},
	}

	app.Before = func(c *cli.Context) error {
		// 開始前の処理をここに書く
		fmt.Println("処理開始")
		return nil // error を返すと処理全体が終了
	}

	app.After = func(c *cli.Context) error {
		// 終了時の処理をここに書く
		fmt.Println("処理終了")
		return nil
	}

	app.Run(os.Args)
}

func getAction(c *cli.Context) {

	// グローバルオプション
	var isDry = c.GlobalBool("dryrun")
	if isDry {
		fmt.Println("this is dry-run")
	}

	// パラメータ
	query := ""
	maxRequestNum := 10
	maxRequestCitedNum := 10
	crawlInterval := 3

	if len(c.Args()) ==  0{
		fmt.Println("[ERROR] Please input search query!")
		os.Exit(1)
	}
	query = strings.Join(c.Args(),"+") // c.Args()[0] と同じ意味

	fmt.Printf("Query : %s\n", query)

	if i := c.Int("maxRequestNum"); i > 0{
		maxRequestNum = i
	}
	fmt.Printf("maxRequestNum : %d\n", maxRequestNum)

	if i := c.Int("maxRequestCitedNum"); i > 0{
		maxRequestCitedNum = i
	}
	fmt.Printf("maxRequestCitedNum : %d\n", maxRequestCitedNum)
	if i := c.Int("crawlInterval"); i > 0{
		crawlInterval = i
	}
	fmt.Printf("crawlInterval : %d\n", crawlInterval)

	s := scholar.CrawlScholar(query, maxRequestNum, maxRequestCitedNum, crawlInterval)

	fmt.Println(s)
	os.Exit(0)
}