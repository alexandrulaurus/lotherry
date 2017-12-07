package main

import (
    "fmt"
    "net/http"

    "github.com/PuerkitoBio/fetchbot"
    "github.com/PuerkitoBio/goquery"
	"log"
)

func main() {
    f := fetchbot.New(fetchbot.HandlerFunc(handler))
    queue := f.Start()
    queue.SendStringGet("http://ethersecret.com/random")
    queue.Close()
}

func handler(ctx *fetchbot.Context, res *http.Response, err error) {
    if err != nil {
        fmt.Printf("error: %s\n", err)
        return
    }
    defer res.Body.Close()
    fmt.Printf("[%d] %s %s\n", res.StatusCode, ctx.Cmd.Method(), ctx.Cmd.URL())
    if res.StatusCode == http.StatusOK {
		doc, docerror := goquery.NewDocumentFromResponse(res)
		if docerror != nil {
			log.Fatal(docerror)
		}
		doc.Find("tr").Each(func(i int, row *goquery.Selection) {
			row.Find("td").Each(func(i int, column *goquery.Selection) {
				fmt.Printf("%d = %s\n", i, column.Text())
			})
			fmt.Printf("======\n")
		})
	}
}
