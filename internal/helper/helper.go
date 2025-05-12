package helper

import (
	"io"
	"net/http"
	"net/url"
)

var BaseURL = "https://www.animesaturn.cx"

type cookieJar struct {
	jar map[string][]*http.Cookie
}

func (p *cookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	p.jar[u.Host] = cookies
}

func (p *cookieJar) Cookies(u *url.URL) []*http.Cookie {
	return p.jar[u.Host]
}

func NewClient() *http.Client {
	return &http.Client{Jar: &cookieJar{make(map[string][]*http.Cookie)}}
}

type PassThru struct {
	io.Reader
	Progress int64
}

func (pt *PassThru) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	if err == nil {
		pt.Progress += int64(n)
		Progress += int64(n)
	}

	return n, err
}

type IndexedUrl struct {
	Index int
	Url   string
}

type AnimeInfo struct {
	EpisodeCount int
	Tags         []string
	Studio       string
	Status       string
	Plot         string
	FirstEpisode int
	LastEpisode int
}

type Anime struct {
	Info   AnimeInfo
	Title  string
	Url    string
	Poster string
}
