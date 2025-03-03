package helper

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

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
	// resp, err := http.Get(url)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := c.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	if _, err = io.Copy(out, resp.Body); err != nil {
		return err
	}

	return nil
}

func Downloader_mp4(c *http.Client, path string, filename string, jobs <-chan IndexedUrl) {
	for j := range jobs {
		name := filepath.Join(path, filename+strconv.Itoa(j.Index)+".mp4")
		startTime := time.Now()
		fmt.Printf("Inizio download di `%s`...\n", filename+strconv.Itoa(j.Index)+".mp4")

		if err := DownloadFile(c, name, string(j.Url)); err != nil {
			panic(err)
		}

		fmt.Printf("Finito di scaricare `%s` in %s\n", name, time.Since(startTime).String())
	}
}

func Downloader_m3u8(path string, filename string, jobs <-chan IndexedUrl) {
	for j := range jobs {
		outPath := filepath.Join(path, filename+strconv.Itoa(j.Index)+".mp4")
		startTime := time.Now()
		fmt.Printf("Inizio download di `%s`...\n", filename+strconv.Itoa(j.Index)+".mp4")

		if err := ffmpeg.Input(string(j.Url)).Output(
			outPath,
			ffmpeg.KwArgs{"protocol_whitelist": "file,http,https,tcp,tls,crypto", "c": "copy"},
		).Run(); err != nil {
			panic(fmt.Sprintf("FFMPEG failed with error code: %s", err))
		}

		fmt.Printf("Finito di scaricare `%s` in %s\n", filename+strconv.Itoa(j.Index)+".mp4", time.Since(startTime).String())
	}
}
