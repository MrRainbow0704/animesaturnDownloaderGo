package helper

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/MrRainbow0704/animesaturnDownloaderGo/internal/cache"
	log "github.com/MrRainbow0704/animesaturnDownloaderGo/internal/logger"
	"github.com/PuerkitoBio/goquery"

	"github.com/dlclark/regexp2"
)

func GetEpisodeLinks(c *http.Client, u string) ([]string, error) {
	var links []string
	cKey := cache.Key(u)
	if err := cKey.Get(&links); err == nil {
		log.Info("Usando la cache")
		return links, nil
	}

	req, _ := http.NewRequest("GET", u, nil)
	res, err := c.Do(req)
	if err != nil || res.StatusCode != 200 {
		log.Errorf("La richiesta HTTP ha prodotto un errore: %s. Response code: %d\n", err, res.StatusCode)
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Errorf("Errore durante il parsing della pagina: %s\n", err)
		return nil, err
	}
	ep0 := false
	doc.Find("a.bottone-ep").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		links = append(links, href)
		if strings.Contains(href, "ep-0") {
			ep0 = true
		}
	})
	if !ep0 {
		return append([]string{"NO EP 0"}, links...), nil
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

	req, _ := http.NewRequest("GET", u, nil)
	res, err := c.Do(req)
	if err != nil || res.StatusCode != 200 {
		log.Errorf("La richiesta HTTP ha prodotto un errore: %s. Response code: %d\n", err, res.StatusCode)
		return IndexedUrl{}, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Errorf("Errore durante il parsing della pagina: %s\n", err)
		return IndexedUrl{}, err
	}
	var link string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if strings.Contains(href, "/watch?") {
			link = href
		}
	})

	iurl = IndexedUrl{i, link}
	cKey.Set(iurl)
	return iurl, nil
}

func GetVideoLink(c *http.Client, u string, i int) (IndexedUrl, error) {
	var iurl IndexedUrl
	cKey := cache.Key(u, i)
	if err := cKey.Get(&iurl); err == nil {
		log.Info("Usando la cache")
		return iurl, nil
	}

	req, _ := http.NewRequest("GET", u, nil)
	res, err := c.Do(req)
	if err != nil || res.StatusCode != 200 {
		log.Errorf("La richiesta HTTP ha prodotto un errore: %s. Response code: %d\n", err, res.StatusCode)
		return IndexedUrl{}, err
	}
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

	u := fmt.Sprintf(BaseURL+"/animelist?search=%s&page=%d", url.PathEscape(s), p)
	req, _ := http.NewRequest("GET", u, nil)
	res, err := c.Do(req)
	if err != nil || res.StatusCode != 200 {
		log.Errorf("La richiesta HTTP ha prodotto un errore: %s. Response code: %d\n", err, res.StatusCode)
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Errorf("Errore durante il parsing della pagina: %s\n", err)
		return nil, err
	}
	doc.Find(".item-archivio").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(".info-archivio>h3>a").Text())
		href, ok := s.Find(".info-archivio>h3>a").Attr("href")
		if !ok {
			log.Error("Errore durante il parsing del link.\n")
			return
		}
		poster, ok := s.Find(".locandina-archivio").Attr("src")
		if !ok {
			log.Errorf("Errore durante il parsing del poster.\n")
			return
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

	req, _ := http.NewRequest("GET", u, nil)
	res, err := c.Do(req)
	if err != nil || res.StatusCode != 200 {
		log.Errorf("La richiesta HTTP ha prodotto un errore: %s. Response code: %d\n", err, res.StatusCode)
		return AnimeInfo{}, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Errorf("Errore durante il parsing della pagina: %s\n", err)
		return AnimeInfo{}, err
	}
	studioReg := regexp2.MustCompile(`(?<=Studio: )(.*)(?=\n)`, 0)
	statoReg := regexp2.MustCompile(`(?<=Stato: )(.*)(?=\n)`, 0)
	episodiReg := regexp2.MustCompile(`(?<=Episodi: )([\d?]*)(?=.*\n)`, 0)

	stats := strings.TrimSpace(doc.Find(".margin-anime-page:nth-child(2)>:nth-child(2)").Text())
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
		strings.TrimSpace(doc.Find(".margin-anime-page:nth-child(2)>:nth-last-child(3)").Text()),
		"\n",
	)
	hentai := false
	if doc.Find(".margin-anime-page:nth-child(2)>div").Length() == 6 {
		hentai = true
	}
	for i, t := range tags {
		tags[i] = strings.TrimSpace(t)
	}
	plot := strings.TrimSpace(doc.Find("#full-trama").Text())

	var epList []string
	doc.Find("#resultsxd>div>div>div>a").Each(func(i int, s *goquery.Selection) {
		epList = append(epList,
			strings.TrimSpace(strings.Split(strings.TrimSpace(s.Text()), " ")[1]),
		)
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

	req, _ := http.NewRequest("GET", BaseURL, nil)
	res, err := c.Do(req)
	if err != nil || res.StatusCode != 200 {
		log.Errorf("La richiesta HTTP ha prodotto un errore: %s. Response code: %d\n", err, res.StatusCode)
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Errorf("Errore durante il parsing della pagina: %s\n", err)
		return nil, err
	}
	doc.Find(".carousel-caption>a").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		href, ok := s.Attr("href")
		if !ok {
			log.Error("Errore durante il parsing del link.\n")
			return
		}
		href = BaseURL + href
		req1, _ := http.NewRequest("GET", href, nil)
		res1, err := c.Do(req1)
		if err != nil || res1.StatusCode != 200 {
			log.Errorf("La richiesta HTTP ha prodotto un errore: %s. Response code: %d\n", err, res.StatusCode)
			return
		}
		doc1, err := goquery.NewDocumentFromReader(res1.Body)
		if err != nil {
			log.Errorf("Errore durante il parsing della pagina: %s\n", err)
			return
		}
		poster, ok := doc1.Find(".cover-anime").Attr("src")
		if !ok {
			log.Errorf("Errore durante il parsing del poster.\n")
			return
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

	u := fmt.Sprintf(BaseURL+"/animelist?search=%s", url.PathEscape(s))
	req, _ := http.NewRequest("GET", u, nil)
	res, err := c.Do(req)
	if err != nil || res.StatusCode != 200 {
		log.Errorf("La richiesta HTTP ha prodotto un errore: %s. Response code: %d\n", err, res.StatusCode)
		return 0, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Errorf("Errore durante il parsing della pagina: %s\n", err)
		return 0, err
	}
	var pagine int
	pageRegexp := regexp2.MustCompile(`(?<=totalPages: )(\d*)(?=,)`, 0)
	code := doc.Find("body>div>script").Text()
	if code == "" {
		log.Info("No pagine.")
		return 1, nil
	}
	text, err := pageRegexp.FindStringMatch(code)
	if err != nil {
		log.Errorf("Errore durante il parsing dello script: %s\n", err)
		return 0, err
	}
	pagine, err = strconv.Atoi(text.String())
	if err != nil {
		log.Errorf("Errore durante la conversione delle pagine: %s", err)
		return 1, nil
	}

	cKey.Set(uint(pagine))
	return uint(pagine), nil
}
