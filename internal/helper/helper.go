package helper

import (
	"io"
	"net/http"
	"net/url"

	log "github.com/MrRainbow0704/animesaturnDownloaderGo/internal/logger"
	"github.com/etherlabsio/go-m3u8/m3u8"
)

var (
	BaseURL  = "https://www.animesaturn.cx"
	Progress float64
	Total    float64
)

type IndexedUrl struct {
	Index int
	Url   string
}

type Anime struct {
	// Info   AnimeInfo
	Title  string
	Url    string
	Poster string
}

type AnimeInfo struct {
	EpisodeCount int
	Is18plus     bool
	Tags         []string
	Studio       string
	Status       string
	Plot         string
	EpisodesList []string
}

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

type segment struct {
	Url  string
	Size float64
}

func (sg *segment) Close() {
	Progress -= sg.Size
	Total -= sg.Size
}

type passThru struct {
	io.Reader
	Progress float64
	Size     float64
}

func (pt *passThru) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	if err == nil {
		pt.Progress += float64(n)
		Progress += float64(n)
	}
	return n, err
}

func (pt *passThru) Close() {
	Progress -= pt.Progress
	Total -= pt.Size
}

func ProgressStart_mp4(c *http.Client, us []IndexedUrl) {
	for _, u := range us {
		req, _ := http.NewRequest("HEAD", u.Url, nil)
		res, err := c.Do(req)
		if err != nil || res.StatusCode != 200 {
			log.Errorf("La richiesta HTTP ha prodotto un errore: %s. Response code: %d\n", err, res.StatusCode)
			return
		}
		Total += float64(res.ContentLength)
	}
}

func ProgressStart_m3u8(c *http.Client, us []IndexedUrl) {
	for _, u := range us {
		
		req, _ := http.NewRequest("GET", u.Url, nil)
		res, err := c.Do(req)
		if err != nil {
			log.Errorf("La richiesta HTTP ha prodotto un errore: %s. Response code: %d\n", err, res.StatusCode)
			return
		}
		playlist, err := m3u8.Read(res.Body)
		if err != nil {
			log.Errorf("Errore nella lettura della playlist: %s\n", err)
			return
		}
		defer res.Body.Close()
		Total += playlist.Duration()
	}
}
