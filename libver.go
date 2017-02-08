package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

const (
	trigger = "@link"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ファイルパスを指定してください。")
		return
	}


	filepath := os.Args[1]
	urls, err := parseFile(filepath)
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}

	if len(urls) == 0 {
		fmt.Println("urlが見つかりませんでした。")
		return
	}

	chStr := make(chan string, len(urls))
	chErr := make(chan error, len(urls))
	for _, u := range urls {
		go fetchReleaseVersion(u+"/releases", chStr, chErr)
	}

	func() {
		for {
			select {
			case str := <-chStr:
				fmt.Println(str)
			case e := <-chErr:
				fmt.Println(e)
			case <-time.After(10 * time.Second):
				return
			}
		}
	}()
}

// `// @link url` このようになっていたら、`url`を取り出す
func parseFile(filePath string) ([]string, error) {
	urls := make([]string, 0)

	fp, err := os.Open(filePath)
	if err != nil {
		return []string{}, errors.Wrap(err, "open failed")
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		target := scanner.Text()
		if strings.Contains(target, trigger) {
			trimmedTarget := strings.Trim(target, " ")
			splitTargets := strings.Split(trimmedTarget, " ")
			if len(splitTargets) > 2 {
				url := splitTargets[2]
				urls = append(urls, url)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return []string{}, errors.Wrap(err, "scanner failed")
	}
	return urls, nil

}

// url一覧からバージョンを取得する
func fetchReleaseVersion(url string, c chan string, e chan error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		e <- errors.Wrap(err, "url scarapping failed")
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
