package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	trigger = "@link"
)

func main() {
	// `// @link url` このようになっていたら、`url`を取り出す
	urls := make([]string, 0)

	fp, err := os.Open("app/build.gradle")
	if err != nil {
		fmt.Println(err)
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		target := scanner.Text()
		if strings.Contains(scanner.Text(), trigger) {
			trimmedTarget := strings.Trim(target, " ")
			splitTargets := strings.Split(trimmedTarget, " ")
			if len(splitTargets) > 2 {
				url := splitTargets[2]
				urls = append(urls, url)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	// url一覧からバージョンを取得する
	c := make(chan string, len(urls))
	e := make(chan error, len(urls))
	for _, u := range urls {
		go fetchReleaseVersion(u+"/releases", c, e)
	}

wait:
	for {
		select {
		case a := <-c:
			fmt.Println(a)
		case errrr := <-e:
			fmt.Println(errrr)
		case <-time.After(10 * time.Second):
			break wait
		}
	}

	return

}

func fetchReleaseVersion(url string, c chan string, e chan error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Print("url scarapping failed")
		e <- err
	}

	// FIXME: https://github.com/google/gson/releases このようにlatestが無いパターンがある

	doc.Find("div.release-meta").EachWithBreak(func(_ int, s *goquery.Selection) bool {

		latest := s.FilterFunction(func(_ int, s *goquery.Selection) bool {
			return s.Find("span").HasClass("latest")
		})

		if latest == nil {
			return true
		}

		t := latest.Find("ul").First().Find("span").Text()
		//fmt.Println(url)
		ver := fmt.Sprintf("%s:%s", url, t)
		c <- ver
		return false

	})
}
