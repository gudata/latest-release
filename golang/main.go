package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	repoURL := os.Args[1]

	re := regexp.MustCompile(`github.com/([^/]+)/([^/]+).*$`)
	match := re.FindStringSubmatch(repoURL)
	if match == nil {
		fmt.Println("Invalid GitHub repository URL")
		os.Exit(1)
	}
	owner, repo := match[1], match[2]

	// Construct the releases page URL
	url := fmt.Sprintf("https://github.com/%s/%s/releases", owner, repo)
	fmt.Printf("I got that we need to fetch latest release from %s\n", url)

	// Send a GET request to the releases page and get the response HTML
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Find the link to the latest release tarball using a regex pattern
	pattern := `nclude-fragment loading="lazy" src="([^"]+)"`
	re = regexp.MustCompile(pattern)
	latestReleaseLinks := re.FindAllStringSubmatch(string(html), -1)
	latestReleaseLink := latestReleaseLinks[0][1]

	// Send another GET request to the latest release page and get the response HTML
	resp, err = http.Get(latestReleaseLink)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	html, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Find the link to the latest release tarball using a regex pattern
	pattern = fmt.Sprintf(`a href="(/%s/%s/releases/download/[^"]+)" rel="nofollow"`, owner, repo)
	re = regexp.MustCompile(pattern)
	downloadLinks := re.FindAllStringSubmatch(string(html), -1)

	if len(downloadLinks) == 0 {
		log.Fatal("No download links found")
	}

	// Loop through the download links and find the one for Linux x86_64
	fmt.Println(downloadLinks)
	var downloadLink string
	for _, link := range downloadLinks {
		path := link[1]
		if (pathContains(path, "x86_64") || pathContains(path, "amd64") || pathContains(path, "musl")) && pathContains(path, "linux") {
			fmt.Printf("Matched string: %s\n", path)
			downloadLink = fmt.Sprintf("https://github.com%s", path)
			break
		}
	}
	if downloadLink == "" {
		fmt.Println("No match found")
		os.Exit(1)
	}

	// Download the latest release tarball
	resp, err = http.Get(downloadLink)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	filename := getFilenameFromURL(downloadLink)
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Downloaded %s\n", filename)
}

func pathContains(path string, substring string) bool {
	return regexp.MustCompile(substring).MatchString(path)
}

func getFilenameFromURL(url string) string {
	segments := strings.Split(url, "/")
	return segments[len(segments)-1]
}
