package helper

import (
	"net/http"
	"net/url"
)

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
