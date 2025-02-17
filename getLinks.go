package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/dlclark/regexp2"
)

func getEpisodeLinks(c *http.Client, u string) ([][]byte, error) {
	linkRegexp := regexp.MustCompile("(?i)<a[^>]*class=[\"'][^>]*bottone-ep[^>]*[\"'][^>]*[^>]*>")
	hrefRegexp := regexp2.MustCompile("(?i)(?<=href=[\"']).*?(?=[\"'])", 0)
	zeroEpRegexp := regexp.MustCompile("(?i)[^0-9A-Za-z]*ep-0[^0-9A-Za-z]*")

	req, _ := http.NewRequest("GET", u, nil)
	res, err := c.Do(req)
	if err != nil || res.StatusCode != 200 {
		fmt.Printf("Errore: %s\n", err)
		fmt.Printf("Status: %s\n", res.Status)
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Errore: %s\n", err)
		return nil, err
	}
	content := strings.Replace(string(body), " ", "", -1)
	links := linkRegexp.FindAll([]byte(content), -1)
	linksList := [][]byte{}
	for i := 0; i < len(links); i++ {
		if match, err := hrefRegexp.FindStringMatch(string(links[i])); err == nil {
			linksList = append(linksList, []byte(match.String()))
		}
	}
	if zeroEpRegexp.FindAll(linksList[0], -1) != nil {
		return linksList, nil
	}
	return append([][]byte{[]byte("NO EP 0")}, linksList...), nil
}

func getStreamLink(c *http.Client, u string, i int) (indexedUrl, error) {
	linkRegexp := regexp.MustCompile("(?i)<a[^>]*href=[\"'][^>]*watch\\?[^>]*[\"'][^>]*>")
	hrefRegexp := regexp2.MustCompile("(?i)(?<=href=[\"']).*?(?=[\"'])", 0)

	req, _ := http.NewRequest("GET", u, nil)
	res, err := c.Do(req)
	if err != nil || res.StatusCode != 200 {
		fmt.Printf("Errore: %s\n", err)
		fmt.Printf("Status: %s\n", res.Status)
		return indexedUrl{}, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Errore: %s\n", err)
		return indexedUrl{}, err
	}
	content := strings.Replace(string(body), " ", "", -1)
	link := linkRegexp.FindAll([]byte(content), -1)[0]
	var streamLink string
	if match, err := hrefRegexp.FindStringMatch(string(link)); err == nil {
		streamLink = match.String()
	}

	return indexedUrl{i, []byte(streamLink)}, nil
}

func getVideoLink(c *http.Client, u string, i int) (indexedUrl, error) {
	mp4Regexp := regexp2.MustCompile("https:\\/\\/.*?(?=\\.mp4)", 0)
	m3u8Regexp := regexp2.MustCompile("https:\\/\\/.*?(?=\\.m3u8)", 0)
	req, _ := http.NewRequest("GET", u, nil)
	res, err := c.Do(req)
	if err != nil || res.StatusCode != 200 {
		fmt.Printf("Errore: %s\n", err)
		fmt.Printf("Status: %s\n", res.Status)
		return indexedUrl{}, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Errore: %s\n", err)
		return indexedUrl{}, err
	}
	content := strings.Replace(string(body), " ", "", -1)
	mp4link, err := mp4Regexp.FindStringMatch(content)
	if err != nil {
		fmt.Printf("Errore: %s\n", err)
		return indexedUrl{}, err
	}
	if mp4link.Length == 0 {
		m3u8link, err := m3u8Regexp.FindStringMatch(content)
		if err != nil {
			fmt.Printf("Errore: %s\n", err)
			return indexedUrl{}, err
		}
		return indexedUrl{i, []byte(m3u8link.String() + ".m3u8")}, nil
	}
	return indexedUrl{i, []byte(mp4link.String() + ".mp4")}, nil
}
