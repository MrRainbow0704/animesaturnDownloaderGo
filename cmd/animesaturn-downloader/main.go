package main

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/MrRainbow0704/animesaturnDownloaderGo/internal/helper"
	log "github.com/MrRainbow0704/animesaturnDownloaderGo/internal/logger"
	"github.com/MrRainbow0704/animesaturnDownloaderGo/internal/version"
)

func main() {
	cwd, _ := os.Getwd()
	downloadCommand := flag.NewFlagSet("download", flag.ExitOnError)
	searchCommand := flag.NewFlagSet("search", flag.ExitOnError)

	flag.Usage = func() {
		log.Print(`AnimesaturnDownloader è una utility per scaricare gli anime dal sito AnimeSaturn.
Scritto in Go da Marco Simone.


Questa schermata di aiuto è divisa in più parti, usa "animesaturn-downloader <sottocomando> -h" per vedere la schermata di aiuto per il sottocomando specifico.

I sottocomandi disponibili sono:
  download		Scarica gli episodi di un anime
  search		Cerca un anime per nome

Utilizzo: animesaturn-downloader <sottocomando> [opzioni]

Flag globali:
  -h, --help		stampa le informazioni di aiuto
  -v, --verbose		stampa altre informazioni di debug
  -V, --version		stampa la versione del programma e termina il programma
`)
	}
	downloadCommand.Usage = func() {
		log.Print(`AnimesaturnDownloader è una utility per scaricare gli anime dal sito AnimeSaturn.
Scritto in Go da Marco Simone.


Schermata di aiuto per il sottocomando "download".

Utilizzo: animesaturn-downloader download -u <url> -n <file> [-v] [-d <percorso>] [-f <numero>] [-l <numero>] [-w <numero>]

Flag per il sottocomando "download":
  -u, --url <url>		link alla pagina dell'anime		[obbligatorio]
  -n, --filename <file>		nome del file senza estensione		[obbligatorio]
  -d, --dir <percorso>		percorso dove salvare i file		[default: percorso corrente]
  -f, --first <numero>		primo episodio da scaricare		[default: 0]
  -l, --last <numero>		ultimo episodio da scaricare		[default: -1]
  -w, --worker <numero>		quanti worker da utilizzare		[default: 3]
`)
	}
	searchCommand.Usage = func() {
		log.Print(`AnimesaturnDownloader è una utility per scaricare gli anime dal sito AnimeSaturn.
Scritto in Go da Marco Simone.


Schermata di aiuto per il sottocomando "search".

Utilizzo: animesaturn-downloader search -s <search>

Flag per il sottocomando "search":
  -s, --search <search>		nome dell'anime da cercare		[obbligatorio]
  -p, --page <numero>		pagina da cercare. 0 => tutte		[default: 0]
  -b, --base-url <url>		url alla home di animesaturn		[obbligatorio]
`)
	}

	// Inizializzazione dei flag
	flag.BoolVar(&log.Verbose, "verbose", false, "stampa altre informazioni di debug")
	flag.BoolVar(&log.Verbose, "v", false, "stampa altre informazioni di debug")
	downloadCommand.BoolVar(&log.Verbose, "verbose", false, "stampa altre informazioni di debug")
	downloadCommand.BoolVar(&log.Verbose, "v", false, "stampa altre informazioni di debug")
	searchCommand.BoolVar(&log.Verbose, "verbose", false, "stampa altre informazioni di debug")
	searchCommand.BoolVar(&log.Verbose, "v", false, "stampa altre informazioni di debug")
	var ver bool
	flag.BoolVar(&ver, "version", false, "stampa la versione del programma")
	flag.BoolVar(&ver, "V", false, "stampa la versione del programma")
	downloadCommand.BoolVar(&ver, "version", false, "stampa la versione del programma")
	downloadCommand.BoolVar(&ver, "V", false, "stampa la versione del programma")
	searchCommand.BoolVar(&ver, "version", false, "stampa la versione del programma")
	searchCommand.BoolVar(&ver, "V", false, "stampa la versione del programma")
	var link string
	downloadCommand.StringVar(&link, "url", "", "link alla pagina dell'anime")
	downloadCommand.StringVar(&link, "u", "", "link alla pagina dell'anime")
	var filename string
	downloadCommand.StringVar(&filename, "filename", "", "nome del file, senza numero di episodio e estensione")
	downloadCommand.StringVar(&filename, "n", "", "nome del file, senza numero di episodio e estensione")
	var path string
	downloadCommand.StringVar(&path, "dir", cwd, "percorso dove salvare i file")
	downloadCommand.StringVar(&path, "d", cwd, "percorso dove salvare i file")
	var primo int
	downloadCommand.IntVar(&primo, "first", 0, "primo episodio da scaricare")
	downloadCommand.IntVar(&primo, "f", 0, "primo episodio da scaricare")
	var ultimo int
	downloadCommand.IntVar(&ultimo, "last", -1, "ultimo episodio da scaricare")
	downloadCommand.IntVar(&ultimo, "l", -1, "ultimo episodio da scaricare")
	var workers int
	downloadCommand.IntVar(&workers, "workers", 3, "quanti worker da utilizzare")
	downloadCommand.IntVar(&workers, "w", 3, "quanti worker da utilizzare")
	var search string
	searchCommand.StringVar(&search, "search", "", "nome dell'anime da cercare")
	searchCommand.StringVar(&search, "s", "", "nome dell'anime da cercare")
	var base string
	searchCommand.StringVar(&base, "base-url", "", "url alla home di animesaturn")
	searchCommand.StringVar(&base, "b", "", "url alla home di animesaturn")
	var page uint
	searchCommand.UintVar(&page, "page", 0, "pagina da cercare")
	searchCommand.UintVar(&page, "p", 0, "pagina da cercare")

	flag.Parse()

	if ver {
		log.Printf("AnimesaturnDownloaderGo %s", version.Get())
		return
	}
	if len(os.Args) < 2 {
		log.Fatal("Nessun sottocomando specificato.\nUsa \"animesaturn-downloader -h\" per vedere la schermata di aiuto.")
		return
	}

	switch os.Args[1] {
	case "download":
		downloadCommand.Parse(os.Args[2:])
		if ver {
			log.Printf("AnimesaturnDownloaderGo %s", version.Get())
			return
		}
		link = strings.TrimSpace(link)
		filename = strings.TrimSpace(filename)
		if link == "" || filename == "" {
			log.Fatal("I flag --filename (-n) e --dir (-u) sono obbligatori")
			return
		}

		path = strings.TrimSpace(path)
		if !filepath.IsAbs(path) {
			path = filepath.Join(cwd, path)
		}
		if err := os.MkdirAll(path, 0777); err != nil {
			log.Printf("Errore durante la creazione della directory `%s`: %s", path, err)
		}

		if primo < 0 {
			log.Print("Il primo episodio deve essere maggiore o uguale a 0")
			return
		}
		if ultimo < primo {
			log.Print("L'ultimo episodio deve essere maggiore o uguale al primo")
			return
		}
		if workers < 1 {
			log.Print("Il numero di workers deve essere maggiore di 0")
			return
		}

		log.Infof("Scaricando da `%s`\n"+
			"\tda episodio %d a %d\n"+
			"\tin `%s`\n"+
			"\tcon nome `%s{%d-%d}.mp4`\n"+
			"\tusando %d workers.\n",
			link, primo, ultimo, path, filename, primo, ultimo, workers,
		)

		startTime := time.Now()
		runDownload(link, primo, ultimo, path, filename, workers)
		log.Infof("Tempo inpiegato: %s\n", time.Since(startTime).String())
	case "search":
		searchCommand.Parse(os.Args[2:])
		if ver {
			log.Printf("AnimesaturnDownloaderGo %s", version.Get())
			return
		}
		search = strings.TrimSpace(search)
		if search == "" {
			log.Fatal("Il flag --search (-s) è obbligatorio")
			return
		}
		if base != "" {
			helper.BaseURL = strings.Trim(base, "/")
		}
		runSearch(search, page)
	default:
		log.Fatalf("Sottocomando `%s` non riconosciuto.\n", os.Args[1])
	}
}

func runDownload(link string, primo int, ultimo int, path string, filename string, workers int) {
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

func runSearch(search string, page uint) {
	log.Infoln("Inizializzando la sessione...")
	client := helper.NewClient()
	log.Infoln("Sessione creata!")

	log.Infoln("Cercando gli anime...")
	var anime []helper.Anime
	if a, err := helper.GetSearchResults(client, search, page); err == nil {
		anime = a
	} else {
		log.Fatalf("Errore nello scraping dei link agli anime: %s\n", err)
	}
	log.Infof("Trovati %d anime che corrispondono alla ricerca.\n", len(anime))

	for _, a := range anime {
		log.Printf("Titolo: %s\n", a.Title)
		log.Printf("Url: %s\n", a.Url)
	}
}
