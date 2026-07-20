package helper

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dlclark/regexp2"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"

	"github.com/MrRainbow0704/animesaturnDownloaderGo/internal/cache"
	log "github.com/MrRainbow0704/animesaturnDownloaderGo/internal/logger"
	"github.com/MrRainbow0704/animesaturnDownloaderGo/internal/version"
)

func GetEpisodeLinks(c *http.Client, u string) ([]string, error) {
	var links []string
	cKey := cache.Key(u)
	if err := cKey.Get(&links); err == nil {
		log.Info("Usando la cache")
		return links, nil
	}

	res, err := SendRequest(c, "GET", u)
	if err != nil {
		log.Errorf("La richiesta HTTP ha prodotto un errore: %s. Response code: %d\n", err, res.StatusCode)
		return nil, err
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Errorf("Errore durante il parsing della pagina: %s\n", err)
		return nil, err
	}
	ep0 := false
	doc.Find("a.ep-tile").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if !strings.HasPrefix("http", href) {
			href = BaseURL + href
		}
		links = append(links, href)
		if strings.Contains(href, "ep-0") {
			ep0 = true
		}
	})
	if !ep0 {
		links = append([]string{"NO EP 0"}, links...)
	}

	cKey.Set(links)
	return links, nil
}

func GetStreamLink(c *http.Client, u string, i int) (IndexedUrl, error) {
	var iurl IndexedUrl
	cKey := cache.Key(u, i)
	if err := cKey.Get(&iurl); err == nil {
		log.Info("Usando la cache")
		return iurl, nil
	}

	res, err := SendRequest(c, "GET", u)
	if err != nil {
		log.Errorf("La richiesta HTTP ha prodotto un errore: %s. Response code: %d\n", err, res.StatusCode)
		return IndexedUrl{}, err
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Errorf("Errore durante il parsing della pagina: %s\n", err)
		return IndexedUrl{}, err
	}
	link, ok := doc.Find("a.ept-btn--play").First().Attr("href")
	if !ok {
		log.Errorf("Errore durante il parsing del link\n")
		return IndexedUrl{}, errors.New("errore durante il parsing del link")
	}
	if !strings.HasPrefix("http", link) {
		link = BaseURL + link
	}

	iurl = IndexedUrl{i, link}
	cKey.Set(iurl)
	return iurl, nil
}

// Funzione di decodifica per il link del video. Estratta e convertita in go
// dalla funzione dec() offuscata nel player del sito.
func decode(b string, k string) string {
	/*
		Originale in JavaScript:

		function dec(b, k) {
			if (!b) return "";
			var s = atob(b);
			var o = "";
			var	i;
			k = k || "as";
				for (i = 0; i < s.length; i++) {
					o += String.fromCharCode(s.charCodeAt(i) ^ k.charCodeAt(i % k.length));
				}
			return o;
		}
	*/
	if b == "" {
		return ""
	}
	s, err := base64.StdEncoding.DecodeString(b)
	if err != nil {
		log.Errorf("Errore durante la decodifica base64: %s\n", err)
		return ""
	}
	if k == "" {
		k = "as"
	}
	var o strings.Builder
	for i := range len(s) {
		o.WriteString(string(s[i] ^ k[i%len(k)]))
	}
	return o.String()
}

func GetVideoLink(c *http.Client, u string, i int) (IndexedUrl, error) {
	var iurl IndexedUrl
	cKey := cache.Key(u, i)
	if err := cKey.Get(&iurl); err == nil {
		log.Info("Usando la cache")
		return iurl, nil
	}

	res, err := SendRequest(c, "GET", u)
	if err != nil {
		log.Errorf("La richiesta HTTP ha prodotto un errore: %s. Response code: %d\n", err, res.StatusCode)
		return IndexedUrl{}, err
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Errorf("Errore durante il parsing della pagina: %s\n", err)
		return IndexedUrl{}, err
	}
	var link string
	doc.Find("video").Each(func(i int, s *goquery.Selection) {
		if c := s.ChildrenFiltered("source"); c.Length() > 0 {
			c.Each(func(i int, ss *goquery.Selection) {
				src, _ := ss.Attr("src")
				if strings.Contains(src, ".mp4") || strings.Contains(src, ".m3u8") {
					link = src
				}
			})
			return
		}
		src, _ := s.Attr("src")
		if strings.Contains(src, ".mp4") || strings.Contains(src, ".m3u8") {
			link = src
		}
	})
	if link == "" {
		// Il player è offuscato, usa un browser headless per ottenere i link
		playerURL, ok := doc.Find("#watch-iframe").First().Attr("src")
		if !ok {
			log.Errorf("Errore durante il parsing del link del player\n")
			return IndexedUrl{}, errors.New("errore durante il parsing del link del player")
		}
		p, _ := launcher.LookPath()
		u := launcher.New().Bin(p).Leakless(false).Headless(!version.IsDev()).MustLaunch()
		browser := rod.New().ControlURL(u)
		if err := browser.Connect(); err != nil {
			log.Errorf("Errore durante la connessione al browser headless: %s\n", err)
			return IndexedUrl{}, err
		}
		video, err := browser.MustPage(playerURL).MustWaitStable().Element("video")
		if err != nil {
			log.Errorf("Errore durante l'ottenimento dell'elemento video: %s\n", err)
			return IndexedUrl{}, err
		}
		linkp, err := video.Attribute("src")
		if err != nil {
			log.Errorf("Errore durante l'ottenimento dell'attributo src del video: %s\n", err)
			return IndexedUrl{}, err
		}
		link = *linkp
		browser.MustClose()
	}

	iurl = IndexedUrl{i, link}
	cKey.Set(iurl)
	return iurl, nil
}

func GetSearchResults(c *http.Client, s string, p uint) ([]Anime, error) {
	var anime []Anime
	cKey := cache.Key(s, p, BaseURL)
	if err := cKey.Get(&anime); err == nil {
		log.Info("Usando la cache")
		return anime, nil
	}

	if p == 0 {
		pagine, err := GetPageNumber(c, s)
		if err != nil {
			log.Errorf("Errore durante l'ottenimento delle pagine: %s", err)
			return nil, err
		}
		var anime []Anime
		for i := range pagine {
			a, err := GetSearchResults(c, s, uint(i+1))
			if err != nil {
				return nil, err
			}
			anime = append(anime, a...)
		}

		cKey.Set(anime)
		return anime, nil
	}

	u := fmt.Sprintf(BaseURL+"/filter/%d?key=%s", p, url.PathEscape(s))
	res, err := SendRequest(c, "GET", u)
	if err != nil {
		log.Errorf("La richiesta HTTP ha prodotto un errore: %s. Response code: %d\n", err, res.StatusCode)
		return nil, err
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Errorf("Errore durante il parsing della pagina: %s\n", err)
		return nil, err
	}
	doc.Find(".ac.group").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h3").Text())
		href, ok := s.Attr("href")
		if !ok {
			log.Error("Errore durante il parsing del link.\n")
			return
		}
		if !strings.HasPrefix("http", href) {
			href = BaseURL + href
		}

		poster, ok := s.Find("img").Attr("src")
		if !ok {
			log.Errorf("Errore durante il parsing del poster.\n")
			return
		}
		if !strings.HasPrefix("http", poster) {
			poster = BaseURL + poster
		}
		a := Anime{Title: title, Url: href, Poster: poster}
		anime = append(anime, a)
	})

	cKey.Set(anime)
	return anime, nil
}

func GetAnimeInfo(c *http.Client, u string) (AnimeInfo, error) {
	var info AnimeInfo
	cKey := cache.Key(u)
	if err := cKey.Get(&info); err == nil {
		log.Info("Usando la cache")
		return info, nil
	}

	res, err := SendRequest(c, "GET", u)
	if err != nil {
		log.Errorf("La richiesta HTTP ha prodotto un errore: %s. Response code: %d\n", err, res.StatusCode)
		return AnimeInfo{}, err
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Errorf("Errore durante il parsing della pagina: %s\n", err)
		return AnimeInfo{}, err
	}
	studioReg := regexp2.MustCompile(`(?<=Studio: )(.*)(?=\n)`, 0)
	statoReg := regexp2.MustCompile(`(?<=Stato: )(.*)(?=\n)`, 0)
	episodiReg := regexp2.MustCompile(`(?<=Episodi: )([\d?]*)(?=.*\n)`, 0)

	stats := strings.TrimSpace(doc.Find("aside>:nth-child(4)").Text())
	studio, err := studioReg.FindStringMatch(stats)
	if err != nil {
		log.Errorf("Errore durante il parsing dello studio: %s\n", err)
		return AnimeInfo{}, err
	}
	status, err := statoReg.FindStringMatch(stats)
	if err != nil {
		log.Errorf("Errore durante il parsing dello status: %s\n", err)
		return AnimeInfo{}, err
	}
	eps, err := episodiReg.FindStringMatch(stats)
	if err != nil {
		log.Errorf("Errore durante il parsing del numero di episodi: %s\n", err)
		return AnimeInfo{}, err
	}
	nEspsStr := eps.Capture.String()
	nEps := 0
	if !strings.Contains(nEspsStr, "?") {
		nEps, err = strconv.Atoi(nEspsStr)
		if err != nil {
			log.Errorf("Errore durante la conversione del numero di episodi: %s\n", err)
			return AnimeInfo{}, err
		}
	}
	tags := strings.Split(
		strings.TrimSpace(doc.Find("header>.ag-genres").Text()),
		"\n",
	)
	hentai := false
	if doc.Find(".adult-gate__backdrop").Length() != 0 {
		hentai = true
	}
	for i, t := range tags {
		tags[i] = strings.TrimSpace(t)
	}
	plot := strings.TrimSpace(doc.Find("section.ag-story>div").Text())

	var epList []string
	doc.Find("a.ep-tile").Each(func(i int, s *goquery.Selection) {
		epList = append(epList, strings.TrimSpace(s.Text()))
	})

	info = AnimeInfo{
		EpisodeCount: nEps,
		Tags:         tags,
		Studio:       studio.Capture.String(),
		Status:       status.Capture.String(),
		Plot:         plot,
		EpisodesList: epList,
		Is18plus:     hentai,
	}
	cKey.Set(info)
	return info, nil
}

func GetDefaultAnime(c *http.Client) ([]Anime, error) {
	var anime []Anime
	cKey := cache.Key(BaseURL)
	if err := cKey.Get(&anime); err == nil {
		log.Info("Usando la cache")
		return anime, nil
	}

	res, err := SendRequest(c, "GET", BaseURL)
	if err != nil {
		log.Errorf("La richiesta HTTP ha prodotto un errore: %s. Response code: %d\n", err, res.StatusCode)
		return nil, err
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Errorf("Errore durante il parsing della pagina: %s\n", err)
		return nil, err
	}
	doc.Find(".swiper-slide>.hero-slide").Each(func(i int, s *goquery.Selection) {
		title := s.Find("h2").First().Text()
		href, ok := s.Find(".hero-actions>.hero-btn-info").Attr("href")
		if !ok {
			log.Error("Errore durante il parsing del link.\n")
			return
		}
		if !strings.HasPrefix("http", href) {
			href = BaseURL + href
		}
		res1, err := SendRequest(c, "GET", href)
		if err != nil {
			log.Errorf("La richiesta HTTP ha prodotto un errore: %s. Response code: %d\n", err, res.StatusCode)
			return
		}
		defer res1.Body.Close()
		doc1, err := goquery.NewDocumentFromReader(res1.Body)
		if err != nil {
			log.Errorf("Errore durante il parsing della pagina: %s\n", err)
			return
		}
		poster, ok := doc1.Find(".anime-poster-card>img").Attr("src")
		if !ok {
			log.Errorf("Errore durante il parsing del poster.\n")
			return
		}
		if !strings.HasPrefix("http", poster) {
			poster = BaseURL + poster
		}

		a := Anime{Title: title, Url: href, Poster: poster}
		anime = append(anime, a)
	})

	cKey.Set(anime)
	return anime, nil
}

func GetPageNumber(c *http.Client, s string) (uint, error) {
	var pages uint
	cKey := cache.Key(s, BaseURL)
	if err := cKey.Get(&pages); err == nil {
		log.Info("Usando la cache")
		return pages, nil
	}

	u := fmt.Sprintf(BaseURL+"/filter?key=%s", url.PathEscape(s))
	res, err := SendRequest(c, "GET", u)
	if err != nil {
		log.Errorf("La richiesta HTTP ha prodotto un errore: %s. Response code: %d\n", err, res.StatusCode)
		return 0, err
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Errorf("Errore durante il parsing della pagina: %s\n", err)
		return 0, err
	}
	var pagine int
	text := doc.Find("nav.mt-section a.page-num:nth-last-child(2)").Text()
	if text == "" {
		log.Info("No pagine.")
		return 1, nil
	}

	pagine, err = strconv.Atoi(text)
	if err != nil {
		log.Errorf("Errore durante la conversione delle pagine: %s", err)
		return 1, nil
	}

	cKey.Set(uint(pagine))
	return uint(pagine), nil
}
