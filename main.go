package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"net/http"
	"strings"
	"time"
	"path/filepath"

	"github.com/mmcdole/gofeed"
)

// Release 结构体，表示每个发布的版本信息
type Release struct {
	Tag     string    // 版本号
	Content string    // 发布说明
	URL     string    // 链接
	Date    time.Time // 发布时间
}

// Commit 结构体，表示每个提交的信息
type Commit struct {
	Title   string    // 提交标题
	URL     string    // 提交链接
	Date    time.Time // 提交时间
}

// GetReleases 函数解析 Atom feed 并返回发布的版本信息
func GetReleases(feedURL string) ([]Release, error) {
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 忽略 HTTPS 证书验证
	}
	client := &http.Client{Transport: customTransport}

	fp := gofeed.NewParser()
	fp.Client = client

	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	releases := make([]Release, 0)
	for _, item := range feed.Items {
		if item.PublishedParsed == nil {
			continue
		}

		releases = append(releases, Release{
			Tag:     item.Title,
			Content: item.Content,
			URL:     item.Link,
			Date:    *item.PublishedParsed,
		})
	}

	return releases, nil
}

// GetCommits 函数解析 Atom feed 并返回提交的版本信息
func GetCommits(feedURL string) ([]Commit, error) {
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 忽略 HTTPS 证书验证
	}
	client := &http.Client{Transport: customTransport}

	fp := gofeed.NewParser()
	fp.Client = client

	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	commits := make([]Commit, 0)
	for _, item := range feed.Items {
		if item.PublishedParsed == nil {
			continue
		}

		commits = append(commits, Commit{
			Title: item.Title,
			URL:   item.Link,
			Date:  *item.PublishedParsed,
		})
	}

	return commits, nil
}

func printHelp() {
	fmt.Println("使用说明：")
	fmt.Println("  -u    指定 GitHub 仓库链接（必填项）")
	fmt.Println("  -r    解析版本记录 releases.atom")
	fmt.Println("  -c    解析提交记录 commits.atom")
	fmt.Println("  -all  同时解析 releases.atom 和 commits.atom")
	fmt.Println("  -n    指定检查的数量，默认为 1")
	fmt.Println("  -o    指定输出文件的目录路径（发布版本为releases.atom 提交记录为commits.atom）")
}

func main() {
	// 定义命令行参数
	urlFlag := flag.String("u", "", "指定 GitHub 仓库链接（必填项）")
	releasesFlag := flag.Bool("r", false, "解析版本记录 releases.atom")
	commitsFlag := flag.Bool("c", false, "解析提交记录 commits.atom")
	allFlag := flag.Bool("all", false, "同时解析 releases.atom 和 commits.atom")
	numFlag := flag.Int("n", 1, "指定检查的数量")
	outputDir := flag.String("o", "", "指定输出文件的目录路径（文件夹）")

	flag.Parse()

	// 检查是否提供了链接
	if *urlFlag == "" {
		printHelp()
		return
	}

	feedURL := *urlFlag

	// 自动补全缺少的后缀
	if !strings.HasSuffix(feedURL, "/releases.atom") && !strings.HasSuffix(feedURL, "/commits.atom") {
		if *releasesFlag || *allFlag {
			feedURL = strings.TrimSuffix(feedURL, "/") + "/releases.atom"
		} else if *commitsFlag {
			feedURL = strings.TrimSuffix(feedURL, "/") + "/commits.atom"
		}
	}
	
	// 检查是否有指定的解析方式，如果没有则默认解析 releases.atom
	if !*releasesFlag && !*commitsFlag && !*allFlag {
		*releasesFlag = true
	}

	// 解析 releases.atom
	if *releasesFlag || *allFlag {
		releases, err := GetReleases(feedURL)
		if err != nil {
			fmt.Println("解析发布版本时出错:", feedURL, err)
			return
		}

		fmt.Println("=== 发布版本信息 ===")
		filePath := filepath.Join(*outputDir, "releases.atom")

    		// 清空文件
    		file, err := os.Create(filePath) // 覆盖模式，清空文件
    		if err != nil {
        		fmt.Printf("创建或清空发布版本文件时出错: %v\n", err)
        		return
    		}
    		defer file.Close() // 确保在函数结束时关闭文件
		for i := 0; i < *numFlag && i < len(releases); i++ {
			release := releases[i]
			localTime := release.Date.In(time.FixedZone("UTC+8", 8*3600))
			output := fmt.Sprintf("发布版本：%s\n发布时间：%s\n发布说明：%s\n版本链接：%s\n\n",
				release.Tag,
				localTime.Format("2006-01-02 15:04:05"),
				release.Content,
				release.URL)
			// 输出到控制台
			fmt.Print(output)
			// 将内容追加写入文件
        		if _, err := file.WriteString(output); err != nil {
            			fmt.Printf("写入发布版本文件时出错: %v\n", err)
            			return
        		}
		}
	}

	// 解析 commits.atom
	if *commitsFlag || *allFlag {
		commitFeedURL := strings.Replace(feedURL, "releases.atom", "commits.atom", 1)
		commits, err := GetCommits(commitFeedURL)
		if err != nil {
			fmt.Println("解析提交记录时出错:", feedURL, err)
			return
		}

		fmt.Println("=== 提交记录信息 ===")
		commitFilePath := filepath.Join(*outputDir, "commits.atom")

    		// 清空文件
    		commitFile, err := os.Create(commitFilePath) // 覆盖模式，清空文件
    		if err != nil {
        		fmt.Printf("创建或清空提交记录文件时出错: %v\n", err)
        		return
    		}
    		defer commitFile.Close() // 确保在函数结束时关闭文件
		for i := 0; i < *numFlag && i < len(commits); i++ {
			commit := commits[i]
			localTime := commit.Date.In(time.FixedZone("UTC+8", 8*3600))
			output := fmt.Sprintf("提交说明：%s\n提交时间：%s\n详细链接：%s\n\n",
				commit.Title,
				localTime.Format("2006-01-02 15:04:05"),
				commit.URL)
			// 输出到控制台
			fmt.Print(output)
			// 将内容追加写入文件
        		if _, err := commitFile.WriteString(output); err != nil {
            			fmt.Printf("写入提交记录文件时出错: %v\n", err)
            			return
        		}
		}
	}
}
