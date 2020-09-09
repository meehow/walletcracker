package main

import (
	"log"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gocolly/colly/v2"
)

func crawl() []common.Address {
	accounts := make([]string, 0, 10000)
	c := colly.NewCollector()
	c.CacheDir = "cache"
	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		if e.ChildAttr("i", "title") == "Contract" {
			return
		}
		texts := e.ChildTexts("td")
		if len(texts) != 6 {
			panic("unexpected number of columns")
		}
		accounts = append(accounts, texts[1])
	})
	c.OnHTML("a[aria-label=Next]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Println(err)
	})
	c.Visit("https://etherscan.io/accounts/1?ps=100")
	sort.Strings(accounts)
	addresses := make([]common.Address, len(accounts))
	for i := range accounts {
		addresses[i] = common.HexToAddress(accounts[i])
	}
	return addresses
}
