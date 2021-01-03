package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := os.Args[1]
	site := getSiteName(url)

	getList(site, url, "/")
}

func getList(site, url, dirName string) error {
	fmt.Println("Dir:", dirName)

	// if not a full path, then add site to url
	if strings.Index(url, "http://") == -1 {
		url = site + url
	}

	resp, e := http.Get(url)
	if e != nil {
		return nil
	}
	defer resp.Body.Close()

	doc, e := goquery.NewDocumentFromReader(resp.Body)
	doc.Find("li").Each(func(i int, s *goquery.Selection) {
		l, _ := s.Find("a").Attr("href")
		fName := s.Text()
		typ, exist := s.Attr("type")

		if exist && typ == "circle" && fName == "@eaDir" {
			return // skip '@eaDir'
		}

		if exist && typ == "circle" {
			subDirName := dirName + fName + "/"
			getList(site, l, subDirName)
		} else {
			fmt.Println("File:", fName)
			fmt.Println(site + l)
		}
	})

	fmt.Println("-----------------------------------------")

	return nil
}

func getSiteName(url string) string {
	tmpURL := url[7:]
	idx := strings.Index(tmpURL, "/")
	return "http://" + tmpURL[:idx] + "/"
}
