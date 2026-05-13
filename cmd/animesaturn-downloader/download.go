package main

import (
	"flag"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/MrRainbow0704/animesaturnDownloaderGo/internal/cache"
	"github.com/MrRainbow0704/animesaturnDownloaderGo/internal/helper"
	log "github.com/MrRainbow0704/animesaturnDownloaderGo/internal/logger"
)

var downloadCommand = flag.NewFlagSet("download", flag.ExitOnError)

var (
	cwd      string
	link     string
	filename string
	path     string
	primo    int
	ultimo   int
	workers  int
)

func initDownload() {
	cwd, _ = os.Getwd()

	downloadCommand.Usage = func() {
		log.Print(header + `Schermata di aiuto per il sottocomando "download".

Utilizzo: ` + execName + ` download -u <url> -n <file> [-v] [-d <percorso>] [-f <numero>] [-l <numero>] [-w <numero>]

Flag per il sottocomando "download":
  -u, --url <url>		link alla pagina dell'anime		[obbligatorio]
  -n, --filename <file>		nome del file senza estensione		[obbligatorio]
  -d, --dir <percorso>		percorso dove salvare i file		[default: percorso corrente]
  -f, --first <numero>		primo episodio da scaricare		[default: 0]
  -l, --last <numero>		ultimo episodio da scaricare		[default: -1]
  -w, --worker <numero>		quanti worker da utilizzare		[default: 3]
`)
	}

	downloadCommand.BoolVar(&log.Verbose, "verbose", false, "stampa altre informazioni di debug")
	downloadCommand.BoolVar(&log.Verbose, "v", false, "stampa altre informazioni di debug")
	downloadCommand.BoolVar(&localCache, "local-cache", false, "force the program to use local cache")
	downloadCommand.BoolVar(&ver, "version", false, "stampa la versione del programma")
	downloadCommand.BoolVar(&ver, "V", false, "stampa la versione del programma")
	downloadCommand.BoolVar(&cache.NoCachce, "no-cache", false, "non utilizza la cache")
	downloadCommand.IntVar(&helper.MaxRetry, "max-retry", 3, "numero massimo di tentativi per ogni richiesta HTTP")
	downloadCommand.StringVar(&link, "url", "", "link alla pagina dell'anime")
	downloadCommand.StringVar(&link, "u", "", "link alla pagina dell'anime")
	downloadCommand.StringVar(&filename, "filename", "", "nome del file, senza numero di episodio e estensione")
	downloadCommand.StringVar(&filename, "n", "", "nome del file, senza numero di episodio e estensione")
	downloadCommand.StringVar(&path, "dir", cwd, "percorso dove salvare i file")
	downloadCommand.StringVar(&path, "d", cwd, "percorso dove salvare i file")
	downloadCommand.IntVar(&primo, "first", 0, "primo episodio da scaricare")
	downloadCommand.IntVar(&primo, "f", 0, "primo episodio da scaricare")
	downloadCommand.IntVar(&ultimo, "last", -1, "ultimo episodio da scaricare")
	downloadCommand.IntVar(&ultimo, "l", -1, "ultimo episodio da scaricare")
	downloadCommand.IntVar(&workers, "workers", 3, "quanti worker da utilizzare")
	downloadCommand.IntVar(&workers, "w", 3, "quanti worker da utilizzare")
}

func parseDownload(arguments []string) {
	downloadCommand.Parse(arguments)

	link = strings.TrimSpace(link)
	filename = strings.TrimSpace(filename)
	if link == "" || filename == "" {
		log.Fatal("I flag --filename (-n) e --dir (-u) sono obbligatori")
	}

	path = strings.TrimSpace(path)
	if !filepath.IsAbs(path) {
		path = filepath.Join(cwd, path)
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		log.Fatalf("Errore durante la creazione della directory `%s`: %s", path, err)
	}

	if primo < 0 {
		log.Fatalf("Il primo episodio deve essere maggiore o uguale a 0")
	}

	if ultimo < primo {
		log.Fatalf("L'ultimo episodio deve essere maggiore o uguale al primo")
	}

	if workers < 1 {
		log.Fatalf("Il numero di workers deve essere maggiore di 0")
	}

	log.Infof("Scaricando da `%s`\n"+
		"\tda episodio %d a %d\n"+
		"\tin `%s`\n"+
		"\tcon nome `%s{%d-%d}.mp4`\n"+
		"\tusando %d workers.\n",
		link, primo, ultimo, path, filename, primo, ultimo, workers,
	)
	cache.Init(localCache)
}

func runDownload() {
	// Setting up session
	log.Println("Ottenimento dei file...")
	log.Infoln("Inizializzando la sessione...")
	client := helper.NewClient()
	log.Infoln("Sessione creata!")

	// getting episode links
	log.Infoln("Cercando i link agli episodi...")
	var episodi []string
	if e, err := helper.GetEpisodeLinks(client, link); err == nil {
		episodi = e
	} else {
		log.Fatalf("Errore nello scraping dei link agli episodi: %s\n", err)
	}

	if episodi[0] == "NO EP 0" && primo == 0 {
		primo = 1
	}

	if ultimo == -1 {
		ultimo = len(episodi) - 1
	}

	log.Infof("Trovati %d episodi.\n", len(episodi))

	// get stream links
	log.Infoln("Cercando i link alle stream...")
	var epLinks = []helper.IndexedUrl{}
	for i := primo; i <= ultimo; i++ {
		if indexedLink, err := helper.GetStreamLink(client, episodi[i], i); err == nil {
			epLinks = append(epLinks, indexedLink)
		} else {
			log.Fatalf("Errore nello scraping dei link alle stream: %s", err)
		}
	}
	log.Infof("Trovate %d stream.\n", len(epLinks))

	// get file links
	log.Infoln("Cercando i link ai file...")
	var videoLinks = []helper.IndexedUrl{}
	for i := range epLinks {
		if indexedLink, err := helper.GetVideoLink(client, epLinks[i].Url, epLinks[i].Index); err == nil {
			videoLinks = append(videoLinks, indexedLink)
		} else {
			log.Fatalf("Errore nello scraping dei link ai file: %s", err)
		}
	}
	log.Infoln("Link ai file trovati.")
	log.Println("URL ai video ottenuti.")

	log.Println("Inizio download...")
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
		downloadM3U8(m3u8Files, client)
	}

	if len(mp4Files) > 0 {
		// old style downloads
		downloadMP4(mp4Files, client)
	}

	log.Println("Download completati.")
}

func downloadM3U8(files []helper.IndexedUrl, client *http.Client) {
	// new style downloads
	helper.ProgressStartM3U8(client, files)

	jobs := make(chan helper.IndexedUrl, len(files))
	wg := sync.WaitGroup{}
	for range workers {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			helper.DownloaderM3U8(client, path, filename, jobs)
		}(&wg)
	}

	for _, link := range files {
		jobs <- link
	}
	close(jobs) // no more urls, so tell workers to stop their loop
	wg.Wait()
}

func downloadMP4(files []helper.IndexedUrl, client *http.Client) {
	helper.ProgressStartMP4(client, files)

	jobs := make(chan helper.IndexedUrl, len(files))
	wg := sync.WaitGroup{}
	for range workers {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			helper.DownloaderMP4(client, path, filename, jobs)
		}(&wg)
	}

	// continually feed in urls to workers
	for _, link := range files {
		jobs <- link
	}
	close(jobs) // no more urls, so tell workers to stop their loop
	wg.Wait()
}
