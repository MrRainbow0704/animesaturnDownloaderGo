package main

import (
	"bufio"
	"flag"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/MrRainbow0704/animesaturnDownloaderGo/internal/helper"
	log "github.com/MrRainbow0704/animesaturnDownloaderGo/internal/logger"
)

const VERSION = "0.1.0"

type cookieJar struct {
	jar map[string][]*http.Cookie
}

func (p *cookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	p.jar[u.Host] = cookies
}

func (p *cookieJar) Cookies(u *url.URL) []*http.Cookie {
	return p.jar[u.Host]
}

func noFlags() bool {
	found := true
	flag.Visit(func(f *flag.Flag) {
		found = false
	})
	return found
}

func main() {
	cwd, _ := os.Getwd()

	// initialize flags arguments
	const usage = `AnimesaturnDownloader è una utility per scaricare gli anime dal sito AnimeSaturn.
Scritto in Go da Marco Simone.

utilizzo: animesaturn-downloader -u <link> -n <filename> [-v] [-d <dir>] [-f <first>] [-l <last>] [-w <workers>]
  -h, --help		stampa le informazioni di aiuto
  -v, --verbose		stampa altre informazioni di debug
  -V, --version		stampa la versione del programma e termina il programma
  -u, --url		link alla pagina dell'anime		[obbligatorio]
  -n, --filename	nome del file senza estensione		[obbligatorio]
  -d, --dir		percorso dove salvare i file		[default: percorso corrente]
  -f, --first		primo episodio da scaricare		[default: 0]
  -l, --last		ultimo episodio da scaricare		[default: -1]
  -w, --worker		quanti worker da utilizzare		[default: 3]
`

	var link string
	flag.StringVar(&link, "url", "D", "link alla pagina dell'anime")
	flag.StringVar(&link, "u", "D", "link alla pagina dell'anime")
	var workers int
	flag.IntVar(&workers, "workers", 3, "quanti worker da utilizzare")
	flag.IntVar(&workers, "w", 3, "quanti worker da utilizzare")
	var primo int
	flag.IntVar(&primo, "first", 0, "primo episodio da scaricare")
	flag.IntVar(&primo, "f", 0, "primo episodio da scaricare")
	var ultimo int
	flag.IntVar(&ultimo, "last", -1, "ultimo episodio da scaricare")
	flag.IntVar(&ultimo, "l", -1, "ultimo episodio da scaricare")
	var path string
	flag.StringVar(&path, "dir", cwd, "percorso dove salvare i file")
	flag.StringVar(&path, "d", cwd, "percorso dove salvare i file")
	var filename string
	flag.StringVar(&filename, "filename", "D", "nome del file, senza numero di episodio e estensione")
	flag.StringVar(&filename, "n", "D", "nome del file, senza numero di episodio e estensione")
	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "stampa altre informazioni di debug")
	flag.BoolVar(&verbose, "v", false, "stampa altre informazioni di debug")
	var version bool
	flag.BoolVar(&version, "version", false, "stampa la versione del programma")
	flag.BoolVar(&version, "V", false, "stampa la versione del programma")

	flag.Usage = func() { log.Print(usage) }
	flag.Parse()

	if noFlags() {
		// initializing reader
		reader := bufio.NewReader(os.Stdin)

		// Getting the page link
		log.Print("Inserisci il link alla pagina dell'anime: ")
		link, _ = reader.ReadString('\n')

		// Getting the amount of workers to use
		log.Print("Inserisci il numero di workers da usare: ")
		workersStr, _ := reader.ReadString('\n')
		if i, err := strconv.Atoi(workersStr); err == nil {
			workers = i
		} else {
			workers = 0
		}

		// Getting the first and last episodes to download
		log.Print("Inserisci il primo episodio da scaricare: ")
		primoStr, _ := reader.ReadString('\n')
		if i, err := strconv.Atoi(primoStr); err == nil {
			primo = i
		} else {
			primo = 0
		}
		log.Print("Inserisci l'ultimo episodio da scaricare: ")
		ultimoStr, _ := reader.ReadString('\n')
		if i, err := strconv.Atoi(ultimoStr); err == nil {
			ultimo = i
		} else {
			ultimo = -1
		}

		// Getting the path where to store the downloads
		log.Printf("Inserisci il percorso dove salvare i file [Vuoto per: \"%s\"]: ", path)
		path, _ = reader.ReadString('\n')

		// Getting filenames
		log.Print("Inserisci il nome per i file: ")
		filename, _ = reader.ReadString('\n')
	} else if (link == "D" || filename == "D") && !verbose && !version {
		panic("I flag --url e --dir sono obbligatori")
	} else if version {
		log.Printf("AnimesaturnDownloaderGo %s", VERSION)
		return
	}

	if verbose {
		log.Verbose = true
	}

	// Input formatting
	link = strings.TrimSpace(link)
	path = strings.TrimSpace(path)
	if !filepath.IsAbs(path) {
		path = filepath.Join(cwd, path)
	}
	if err := os.MkdirAll(path, 0777); err != nil {
		log.Fatalf("Errore durante la creazione della directory `%s`: %s", path, err)
	}
	filename = strings.TrimSpace(filename)

	log.Infof("Scaricando da `%s`\n"+
		"\tda episodio %d a %d\n"+
		"\tin `%s`\n"+
		"\tcon nome `%s{%d-%d}.mp4`\n"+
		"\tusando %d workers.\n",
		link, primo, ultimo, path, filename, primo, ultimo, workers,
	)

	var startTime = time.Now()
	run(link, primo, ultimo, path, filename, workers)
	log.Infof("Tempo inpiegato: %s\n", time.Since(startTime).String())
}

func run(link string, primo int, ultimo int, path string, filename string, workers int) {
	// Setting up session
	log.Println("Ottenimento dei file...")
	log.Infoln("Inizializzando la sessione...")
	client := &http.Client{Jar: &cookieJar{make(map[string][]*http.Cookie)}}
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
		log.Infoln("Rilevati file m3u8! Inizializando il download tramite FFMPEG...")
		err := exec.Command("ffmpeg", "-version").Run()
		if err != nil {
			log.Fatalln("Per questo tipo di file è necessario FFMPEG. Per favore installalo e riprova.")
		}

		jobs := make(chan helper.IndexedUrl, len(m3u8Files))
		wg := sync.WaitGroup{}
		for range workers {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				helper.Downloader_m3u8(path, filename, jobs)
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
		jobs := make(chan helper.IndexedUrl, len(mp4Files))
		wg := sync.WaitGroup{}
		for range workers {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				helper.Downloader_mp4(client, path, filename, jobs)
			}(&wg)
		}

		// continually feed in urls to workers
		for _, link := range mp4Files {
			jobs <- link
		}
		close(jobs) // no more urls, so tell workers to stop their loop
		wg.Wait()
	}

	log.Println("Download completati.")
}
