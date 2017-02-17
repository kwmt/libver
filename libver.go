package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

var triggers = []string{"compile", "classpath"}

type Target struct {
	GroupID    string
	ArtifactID string
	Version    string
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please specify a build.gradle file path.")
		return
	}
	run(os.Args[1])

}

func run(filepath string) {
	dependencies, err := parseFile(filepath)
	if err != nil {
		fmt.Printf("Parse error. %+v", err)
		return
	}

	if len(dependencies) == 0 {
		fmt.Println("Not found dependency library.")
		return
	}

	c, err := setupClient()
	if err != nil {
		fmt.Println("Failed to setup.", err)
		return
	}

	isFetchedSupportLibraryVersion := false

	for _, target := range dependencies {

		switch target.GroupID {
		case "com.android.support":
			if isFetchedSupportLibraryVersion {
				continue
			}
			mavenPackage, err := fetchLatestSupportLibraryVersion()
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("%s:%s\n", mavenPackage.Name, mavenPackage.LatestVesrsion)
			isFetchedSupportLibraryVersion = true
			continue

		// TODO: Localの$ANDROID_HOME/extrasにあるバージョンと、
		// https://dl.google.com/android/repository/addon2-1.xmlから
		// 取得できるバージョンを取得してチェックしたい
		case "com.android.support.test",
			"com.android.databinding",
			"com.google.android",
			"com.google.firebase",
			"com.google.android.gms",
			"com.google.android.support",
			"com.google.android.wearable",
			"com.android.support.constraint":
			continue
		}

		results, err := c.SearchMavenPackage(target.GroupID, target.ArtifactID)
		if err != nil {
			fmt.Println(target.GroupID, err)
			continue
		}
		r := *results
		if len(r) > 0 {
			fmt.Printf("%s:%s\n", r[0].Name, r[0].LatestVesrsion)
		}
	}
}

func setupClient() (*Client, error) {
	username := os.Getenv("BINTRAY_API_USERNAME")
	if username == "" {
		return nil, errors.New("Environment variabl BINTRAY_API_USERNAME is not set.")
	}
	password := os.Getenv("BINTRAY_API_PASSWORD")
	if username == "" {
		return nil, errors.New("Environment variabl BINTRAY_API_PASSWORD is not set.")
	}
	return NewClient(username, password), nil
}

// 例えば`compile 'com.google.code.gson:gson:2.7'` このようになっていたら com.google.code.gson、gson、2.7に分けたリストを返す。
func parseFile(filePath string) ([]Target, error) {
	fp, err := os.Open(filePath)
	if err != nil {
		return []Target{}, errors.Wrap(err, "Failed to open file.")
	}
	defer fp.Close()

	targets := make([]Target, 0)

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		target := strings.ToLower(scanner.Text())
		if containTriggers(target, triggers) {
			trimmedTarget := strings.Trim(target, " ")
			splitTargets := strings.Split(trimmedTarget, " ")
			if len(splitTargets) > 1 {
				// complie or androidTestCompile etc (exclude compileSdkVersion etc)
				if !hasSuffixTriggers(splitTargets[0], triggers) {
					continue
				}
				target := splitDependency(splitTargets[1])
				// target.GroupIDに以下が入っていたら無視する
				// com.android.support
				if target == nil {
					continue
				}
				targets = append(targets, *target)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return []Target{}, errors.Wrap(err, "Failed to scanner.")
	}
	return targets, nil
}

func containTriggers(target string, triggers []string) bool {
	for _, trigger := range triggers {
		if strings.Contains(target, trigger) {
			return true
		}
	}
	return false
}

func hasSuffixTriggers(target string, triggers []string) bool {
	for _, trigger := range triggers {
		if strings.HasSuffix(target, trigger) {
			return true
		}
	}
	return false
}

// For example, split 'com.google.code.gson:gson:2.7' by ":".
func splitDependency(s string) *Target {
	t := new(Target)
	split := strings.Split(s, ":")
	if len(split) < 3 {
		return nil
	}
	t.GroupID = split[0][1:] // "com.google.firebase などのようにダブル・シングルクォートがあるので取り除く
	t.ArtifactID = split[1]
	t.Version = split[2][:len(split[2])-1] // 0.5" のようにダブル・シングルクォートがあるので取り除く
	return t
}

// unused. Fetch maven package info with goroutine.
func fetchPackage(client *Client, target Target, c chan []ResponseMavenPackageSearch, e chan error) {

	results, err := client.SearchMavenPackage(target.GroupID, target.ArtifactID)
	//fmt.Println(target.GroupID, target.ArtifactID)
	if err != nil {
		e <- errors.Wrap(err, "Failed to search maven package")
		return
	}
	c <- *results
}

// Fetch latest version from android developr Recent Support Library Revisions page https://developer.android.com/topic/libraries/support-library/revisions.html?hl=en
func fetchLatestSupportLibraryVersion() (*MavenPackage, error) {
	url := "https://developer.android.com/topic/libraries/support-library/revisions.html?hl=en"
	doc, err := goquery.NewDocument(url)

	if err != nil {
		return nil, errors.Wrap(err, "url scarapping failed")
	}

	mavenPackage := new(MavenPackage)
	mavenPackage.Name = "com.android.support"

	doc.Find("#body-content > div.jd-descr > h2").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		//val, _ := s.Attr("id")
		revText := s.First().Text() // Revision 25.1.1
		splitRevText := strings.Split(revText, " ")

		if len(splitRevText) < 2 {
			fmt.Println(splitRevText)
			return false
		}

		var validVersion = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)
		if validVersion.MatchString(splitRevText[1]) {
			mavenPackage.LatestVesrsion = splitRevText[1]
			return false
		}
		return false

	})
	return mavenPackage, nil
}
