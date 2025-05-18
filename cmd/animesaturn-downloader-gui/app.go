package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/MrRainbow0704/animesaturnDownloaderGo/config"
	"github.com/MrRainbow0704/animesaturnDownloaderGo/internal/helper"
	"github.com/MrRainbow0704/animesaturnDownloaderGo/internal/logger"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx    context.Context
	client *http.Client
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	appLogger.SetContext(&ctx)
	logger.Verbose = config.Verbose()
	helper.BaseURL = config.BaseURL()

	wails.LogInfo(a.ctx, "Inizializzando la sessione...\n")
	a.client = helper.NewClient()
	wails.LogInfo(a.ctx, "Sessione creata!\n")
}

func (a *App) SearchAnime(s string, p uint) []helper.Anime {
	if p == 0 {
		wails.LogErrorf(a.ctx, "Pagina non valida: %d\n", p)
		return []helper.Anime{}
	}
	animes, err := helper.GetSearchResults(a.client, s, p)
	if err != nil {
		wails.LogErrorf(a.ctx, "Errore durante la ricerca: %s\n", err)
		return []helper.Anime{}
	}
	wails.LogInfof(a.ctx, "Trovati %d anime\n", len(animes))
	return animes
}

func (a *App) DownloadAnime(link string, primo int, ultimo int, filename string, workers int) bool {
	done := make(chan bool)
	defer func() { done <- true }()
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				wails.WindowExecJS(
					a.ctx,
					fmt.Sprintf(
						"window.progressBarProgress = %f; window.progressBarTotal = %f",
						helper.Progress,
						helper.Total,
					),
				)
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
	path, err := wails.OpenDirectoryDialog(
		a.ctx,
		wails.OpenDialogOptions{
			Title:                "Scegli dove salvare i file:",
			CanCreateDirectories: true,
			ResolvesAliases:      true,
		},
	)
	if err != nil {
		return false
	}
	wails.LogPrint(a.ctx, "Iniziando il download...")

	wails.LogInfo(a.ctx, "Cercando i link agli episodi...\n")
	var episodi []string
	if e, err := helper.GetEpisodeLinks(a.client, link); err == nil {
		episodi = e
	} else {
		wails.LogErrorf(a.ctx, "Errore nello scraping dei link agli episodi: %s\n", err)
		return false
	}
	if episodi[0] == "NO EP 0" && primo == 0 {
		primo = 1
	}
	if ultimo == -1 {
		ultimo = len(episodi) - 1
	}

	wails.LogInfof(a.ctx, "Trovati %d episodi.\n", len(episodi))

	// get stream links
	wails.LogInfo(a.ctx, "Cercando i link alle stream...\n")
	var epLinks = []helper.IndexedUrl{}
	for i := primo; i <= ultimo; i++ {
		if indexedLink, err := helper.GetStreamLink(a.client, episodi[i], i); err == nil {
			epLinks = append(epLinks, indexedLink)
		} else {
			wails.LogErrorf(a.ctx, "Errore nello scraping dei link alle stream: %s", err)
			return false
		}
	}
	wails.LogInfof(a.ctx, "Trovate %d stream.\n", len(epLinks))

	// get file links
	wails.LogInfo(a.ctx, "Cercando i link ai file...\n")
	var videoLinks = []helper.IndexedUrl{}
	for i := range epLinks {
		if indexedLink, err := helper.GetVideoLink(a.client, epLinks[i].Url, epLinks[i].Index); err == nil {
			videoLinks = append(videoLinks, indexedLink)
		} else {
			wails.LogErrorf(a.ctx, "Errore nello scraping dei link ai file: %s", err)
			return false
		}
	}
	wails.LogInfo(a.ctx, "Link ai file trovati.\n")
	wails.LogPrint(a.ctx, "URL ai video ottenuti.\n")
	wails.LogPrint(a.ctx, "Inizio download...\n")
	var m3u8Files []helper.IndexedUrl
	var mp4Files []helper.IndexedUrl
	for _, link := range videoLinks {
		if strings.HasSuffix(link.Url, ".m3u8") {
			m3u8Files = append(m3u8Files, link)
		} else {
			mp4Files = append(mp4Files, link)
		}
	}
	if len(m3u8Files) > 0 {
		// new style downloads
		helper.ProgressStart_m3u8(a.client, m3u8Files)

		jobs := make(chan helper.IndexedUrl, len(m3u8Files))
		wg := sync.WaitGroup{}
		for range workers {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				helper.Downloader_m3u8(a.client, path, filename, jobs)
			}(&wg)
		}

		for _, link := range m3u8Files {
			jobs <- link
		}
		close(jobs) // no more urls, so tell workers to stop their loop
		wg.Wait()
	}
	if len(mp4Files) > 0 {
		// old style downloads
		helper.ProgressStart_mp4(a.client, mp4Files)

		jobs := make(chan helper.IndexedUrl, len(mp4Files))
		wg := sync.WaitGroup{}
		for range workers {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				helper.Downloader_mp4(a.client, path, filename, jobs)
			}(&wg)
		}

		// continually feed in urls to workers
		for _, link := range mp4Files {
			jobs <- link
		}
		close(jobs) // no more urls, so tell workers to stop their loop
		wg.Wait()
	}

	wails.LogPrint(a.ctx, "Download completati.\n")
	return true
}

func (a *App) GetDefaultAnime() []helper.Anime {
	animes, err := helper.GetDefaultAnime(a.client)
	if err != nil {
		wails.LogErrorf(a.ctx, "Errore durante la ricerca: %s\n", err)
		return []helper.Anime{}
	}
	wails.LogInfof(a.ctx, "Trovati %d anime\n", len(animes))
	return animes
}

func (a *App) GetAnimeInfo(link string) helper.AnimeInfo {
	info, err := helper.GetAnimeInfo(a.client, link)
	if err != nil {
		wails.LogErrorf(a.ctx, "Errore durante l'ottenimento delle informazioni: %s\n", err)
		return helper.AnimeInfo{}
	}
	return info
}

func (a *App) SetBaseUrl(u string) {
	helper.BaseURL = strings.Trim(u, "/")
}

func (a *App) GetBaseUrl() string {
	return helper.BaseURL
}

func (a *App) GetPageNumber(s string) uint {
	p, err := helper.GetPageNumber(a.client, s)
	if err != nil {
		wails.LogErrorf(a.ctx, "Errore durante la ricerca delle pagine: %s", err)
		return 0
	}
	return p
}
