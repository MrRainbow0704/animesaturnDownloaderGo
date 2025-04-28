package helper

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	log "github.com/MrRainbow0704/animesaturnDownloaderGo/internal/logger"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func DownloadFile(c *http.Client, filepath string, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := c.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Errorf("La richiesta HTTP ha prodotto un errore: %s\n", err)
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	log.Infof("Scrivendo il file `%s`...\n", filepath)
	if _, err = io.Copy(out, resp.Body); err != nil {
		log.Errorf("La scrittura del file `%s` ha prodotto un errore: %s\n",filepath, err)
		return err
	}
	log.Infof("Terminata la scrittura del file `%s`.\n", filepath)
	return nil
}

func Downloader_mp4(c *http.Client, path string, filename string, jobs <-chan IndexedUrl) {
	for j := range jobs {
		name := filepath.Join(path, filename+strconv.Itoa(j.Index)+".mp4")
		startTime := time.Now()
		log.Printf("Inizio download di `%s`...\n", filename+strconv.Itoa(j.Index)+".mp4")

		if err := DownloadFile(c, name, j.Url); err != nil {
			log.Fatalf("Errore durante il download del file `%s`: %s\n", name, err)
		}

		log.Printf("Finito di scaricare `%s` in %s.\n", filename+strconv.Itoa(j.Index)+".mp4", time.Since(startTime).String())
	}
}

func Downloader_m3u8(path string, filename string, jobs <-chan IndexedUrl) {
	for j := range jobs {
		outPath := filepath.Join(path, filename+strconv.Itoa(j.Index)+".mp4")
		startTime := time.Now()
		log.Printf("Inizio download di `%s`...\n", filename+strconv.Itoa(j.Index)+".mp4")

		if err := ffmpeg.Input(j.Url).Output(
			outPath,
			ffmpeg.KwArgs{"protocol_whitelist": "file,http,https,tcp,tls,crypto", "c": "copy"},
		).Run(); err != nil {
			log.Fatalf("Errore durante il download del file `%s` con FFMPEG: %s\n", outPath, err)
		}

		log.Printf("Finito di scaricare `%s` in %s.\n", filename+strconv.Itoa(j.Index)+".mp4", time.Since(startTime).String())
	}
}
