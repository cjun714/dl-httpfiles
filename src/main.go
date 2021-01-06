package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := os.Args[1]
	site := getSiteName(url)

	e := getList(site, url, "/")
	if e != nil {
		log.Fatal(e)
	}
}

func getList(site, url, dirName string) error {
	fmt.Println("Dir:", dirName)

	// if not a full path, then add site to url
	if !strings.Contains(url, "http://") {
		url = site + url
	}

	resp, e := http.Get(url)
	if e != nil {
		i := 0
		for i < 3 { // retry 3 times
			resp, e = http.Get(url)
			if e == nil {
				break
			}
			i++
		}
		if i == 3 {
			return e
		}
	}
	defer resp.Body.Close()

	doc, e := goquery.NewDocumentFromReader(resp.Body)
	if e != nil {
		return e
	}
	doc.Find("li").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Find("a").Attr("href")
		fName := s.Text()
		typ, exist := s.Attr("type")

		if exist && typ == "circle" && fName == "@eaDir" {
			return // skip '@eaDir'
		}

		if exist && typ == "circle" {
			subDirName := dirName + fName + "/"
			e := getList(site, link, subDirName)
			if e != nil {
				log.Println("get list failed: ", url, ". ", e)
			}
		} else {
			fmt.Println("File:", fName)
			fmt.Println(site + link[1:])
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
