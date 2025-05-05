package helper

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	log "github.com/MrRainbow0704/animesaturnDownloaderGo/internal/logger"
	"github.com/PuerkitoBio/goquery"

	"github.com/dlclark/regexp2"
)

var BASEURL = "https://www.animesaturn.cx"

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
	Poster       string
}

type Anime struct {
	Info  AnimeInfo
	Title string
	Url   string
}

func GetEpisodeLinks(c *http.Client, u string) ([]string, error) {
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
	links := []string{}
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
	return links, nil
}

func GetStreamLink(c *http.Client, u string, i int) (IndexedUrl, error) {
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
	return IndexedUrl{i, link}, nil
}

func GetVideoLink(c *http.Client, u string, i int) (IndexedUrl, error) {
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
	return IndexedUrl{i, link}, nil
}

func GetSearchResults(c *http.Client, s string) ([]Anime, error) {
	u := fmt.Sprintf(BASEURL+"/animelist?search=%s", s)
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
	anime := []Anime{}
	doc.Find(".item-archivio").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(".info-archivio>h3>a").Text())
		href, ok := s.Find(".info-archivio>h3>a").Attr("href")
		if !ok {
			log.Error("Errore durante il parsing del link.\n")
			return
		}
		info, err := GetAnimeInfo(c, href)
		if err != nil {
			log.Error("Errore durante il parsing delle informazioni.\n")
			return
		}
		a := Anime{Title: title, Url: href, Info: info}
		anime = append(anime, a)
	})
	return anime, nil
}

func GetAnimeInfo(c *http.Client, u string) (AnimeInfo, error) {
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
	tags := strings.Split(strings.TrimSpace(doc.Find(".margin-anime-page:nth-child(2)>:nth-child(3)").Text()), "\n")
	for i, t := range tags {
		tags[i] = strings.TrimSpace(t)
	}
	plot := strings.TrimSpace(doc.Find("#full-trama").Text())
	poster, ok := doc.Find(".cover-anime").Attr("src")
	if !ok {
		log.Errorf("Errore durante il parsing del poster.\n")
		return AnimeInfo{}, errors.New("inpossibile trovare il poster dell'anime")
	}
	return AnimeInfo{
		EpisodeCount: nEps,
		Tags:         tags,
		Studio:       studio.Capture.String(),
		Status:       status.Capture.String(),
		Plot:         plot,
		Poster:       poster,
	}, nil
}

func GetDefaultAnime(c *http.Client) ([]Anime, error) {
	req, _ := http.NewRequest("GET", BASEURL, nil)
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
	anime := []Anime{}
	doc.Find(".carousel-caption>a").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		href, ok := s.Attr("href")
		if !ok {
			log.Error("Errore durante il parsing del link.\n")
			return
		}
		href = BASEURL + href
		info, err := GetAnimeInfo(c, href)
		if err != nil {
			log.Error("Errore durante il parsing delle informazioni.\n")
			return
		}

		a := Anime{Title: title, Url: href, Info: info}
		anime = append(anime, a)
	})
	return anime, nil
}
