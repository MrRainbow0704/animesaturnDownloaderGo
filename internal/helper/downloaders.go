package helper

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/MrRainbow0704/animesaturnDownloaderGo/internal/cache"
	log "github.com/MrRainbow0704/animesaturnDownloaderGo/internal/logger"
	"github.com/etherlabsio/go-m3u8/m3u8"
)

func downloadFile(c *http.Client, out *os.File, url string) error {
	// Get the data
	req, _ := http.NewRequest("GET", url, nil)
	res, err := c.Do(req)
	if err != nil || res.StatusCode != 200 {
		log.Errorf("La richiesta HTTP ha prodotto un errore: %s\n", err)
		return err
	}
	defer res.Body.Close()

	src := &passThru{Reader: res.Body, Size: float64(res.ContentLength)}
	defer src.Close()

	// Write the body to file
	log.Infof("Scrivendo il file `%s`...\n", out.Name())
	if _, err = io.Copy(out, src); err != nil {
		log.Errorf("La scrittura del file `%s` ha prodotto un errore: %s\n", out.Name(), err)
		return err
	}
	log.Infof("Terminata la scrittura del file `%s`.\n", out.Name())
	return nil
}

func downloadSegment(c *http.Client, out *os.File, seg *segment) error {
	// Get the data
	req, _ := http.NewRequest("GET", seg.Url, nil)
	res, err := c.Do(req)
	if err != nil || res.StatusCode != 200 {
		log.Errorf("La richiesta HTTP ha prodotto un errore: %s\n", err)
		return err
	}
	defer res.Body.Close()

	Progress += seg.Size
	defer seg.Close()

	// Write the body to file
	log.Infof("Scrivendo il file `%s`...\n", out.Name())
	if _, err = io.Copy(out, res.Body); err != nil {
		log.Errorf("La scrittura del file `%s` ha prodotto un errore: %s\n", out.Name(), err)
		return err
	}
	log.Infof("Terminata la scrittura del file `%s`.\n", out.Name())
	return nil
}

func Downloader_mp4(c *http.Client, path string, filename string, jobs <-chan IndexedUrl) {
	for j := range jobs {
		name := filepath.Join(path, filename+strconv.Itoa(j.Index)+".mp4")
		out, err := os.Create(name)
		if err != nil {
			log.Errorf("La creazione del file `%s` ha prodotto un errore: %s\n", out.Name(), err)
			return
		}
		defer out.Close()

		startTime := time.Now()
		log.Printf("Inizio download di `%s`...\n", filename+strconv.Itoa(j.Index)+".mp4")
		if err := downloadFile(c, out, j.Url); err != nil {
			log.Fatalf("Errore durante il download del file `%s`: %s\n", name, err)
		}
		log.Printf("Finito di scaricare `%s` in %s.\n", filename+strconv.Itoa(j.Index)+".mp4", time.Since(startTime).String())
	}
}

func Downloader_m3u8(c *http.Client, path string, filename string, jobs <-chan IndexedUrl) {
	for j := range jobs {
		name := filepath.Join(path, filename+strconv.Itoa(j.Index)+".mp4")
		out, err := os.Create(name)
		if err != nil {
			log.Errorf("La creazione del file `%s` ha prodotto un errore: %s\n", out.Name(), err)
			return
		}
		defer out.Close()

		segs := make(chan *segment)
		startTime := time.Now()
		log.Printf("Inizio download di `%s`...\n", filename+strconv.Itoa(j.Index)+".mp4")
		go getPlaylist(c, j.Url, segs)
		for s := range segs {
			if err := downloadSegment(c, out, s); err != nil {
				log.Fatalf("Errore durante il download del file `%s`: %s\n", name, err)
			}
		}
		log.Printf("Finito di scaricare `%s` in %s.\n", filename+strconv.Itoa(j.Index)+".mp4", time.Since(startTime).String())
	}
}

func getPlaylist(c *http.Client, urlStr string, dlc chan *segment) {
	var segments []segment
	cKey := cache.Key(urlStr)
	if err := cKey.Get(&segments); err == nil {
		log.Info("Usando la cache")
		for _, s := range segments {
			dlc <- &s
		}
		return
	}

	playlistUrl, err := url.Parse(urlStr)
	if err != nil {
		log.Fatal(err)
	}
	req, _ := http.NewRequest("GET", urlStr, nil)
	res, err := c.Do(req)
	if err != nil {
		log.Print(err)
	}
	defer res.Body.Close()
	playlist, err := m3u8.Read(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range playlist.Segments() {
		if v != nil {
			var msURI string
			if strings.HasPrefix(v.Segment, "http") {
				msURI, err = url.QueryUnescape(v.Segment)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				msUrl, err := playlistUrl.Parse(v.Segment)
				if err != nil {
					log.Print(err)
					continue
				}
				msURI, err = url.QueryUnescape(msUrl.String())
				if err != nil {
					log.Fatal(err)
				}
			}
			s := segment{msURI, v.Duration}
			dlc <- &s
			segments = append(segments, s)
		}
	}

	cKey.Set(segments)
	close(dlc)
}
