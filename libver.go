package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

const (
	trigger = "compile"
)

type Target struct {
	GroupID    string
	ArtifactID string
	Version    string
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("ファイルパスを指定してください。")
		return
	}

	filepath := os.Args[1]
	dependencies, err := parseFile(filepath)
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}

	if len(dependencies) == 0 {
		fmt.Println("dependency libraryが見つかりませんでした。")
		return
	}

	c, err := setupClient()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, target := range dependencies {

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
		return nil, errors.New("環境変数 BINTRAY_API_USERNAME が設定されていません。")
	}
	password := os.Getenv("BINTRAY_API_PASSWORD")
	if username == "" {
		return nil, errors.New("環境変数 BINTRAY_API_PASSWORD が設定されていません。")
	}
	return NewClient(username, password), nil
}

// 例えば`compile 'com.google.code.gson:gson:2.7'` このようになっていたら com.google.code.gson、gson、2.7に分けたリストを返す。
func parseFile(filePath string) ([]Target, error) {
	targets := make([]Target, 0)

	fp, err := os.Open(filePath)
	if err != nil {
		return []Target{}, errors.Wrap(err, "open failed")
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		target := scanner.Text()
		if strings.Contains(target, trigger) {
			trimmedTarget := strings.Trim(target, " ")
			splitTargets := strings.Split(trimmedTarget, " ")
			if len(splitTargets) > 1 {
				if !strings.EqualFold(splitTargets[0], trigger) {
					continue
				}
				dependency := splitTargets[1]
				target := splitDependency(dependency)
				targets = append(targets, target)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return []Target{}, errors.Wrap(err, "scanner failed")
	}
	return targets, nil

}

func splitDependency(s string) Target {
	t := Target{}
	split := strings.Split(s, ":")
	if len(split) < 3 {
		return t
	}
	t.GroupID = split[0][1:] // "com.google.firebase などのようにダブル・シングルクォートがあるので取り除く
	t.ArtifactID = split[1]
	t.Version = split[2][:len(split[2])-1] // 0.5" のようにダブル・シングルクォートがあるので取り除く
	return t
}

// Bintrayからパッケージ情報を取得する
func fetchPackage(client *Client, target Target, c chan []ResponseMavenPackageSearch, e chan error) {

	results, err := client.SearchMavenPackage(target.GroupID, target.ArtifactID)
	//fmt.Println(target.GroupID, target.ArtifactID)
	if err != nil {
		e <- errors.Wrap(err, "SearchMavenPackage failed")
		return
	}
	c <- *results
}
