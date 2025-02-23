package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type indexedUrl struct {
	i int
	u []byte
}

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
	const usage = `Flag di AnimesaturnDownloader:
  -h, --help		stampa le informazioni di aiuto
  -w, --worker		quanti worker da utilizzare		[Default: 3]
  -n, --filename	nome del file senza estensione		[Obbligatorio]
  -u, --url		link alla pagina dell'anime		[Obbligatorio]
  -f, --first		primo episodio da scaricare		[Default: 0]
  -l, --last		ultimo episodio da scaricare		[Default: -1]
  -d, --dir		percorso dove salvare i file		[Default: percorso corrente]
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

	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()

	if noFlags() {
		// initializing reader
		reader := bufio.NewReader(os.Stdin)

		// Getting the page link
		fmt.Print("Inserisci il link alla pagina dell'anime: ")
		link, _ = reader.ReadString('\n')

		// Getting the amount of workers to use
		fmt.Print("Inserisci il numero di workers da usare: ")
		workersStr, _ := reader.ReadString('\n')
		if i, err := strconv.Atoi(workersStr); err == nil {
			workers = i
		} else {
			workers = 0
		}

		// Getting the first and last episodes to download
		fmt.Print("Inserisci il primo episodio da scaricare: ")
		primoStr, _ := reader.ReadString('\n')
		if i, err := strconv.Atoi(primoStr); err == nil {
			primo = i
		} else {
			primo = 0
		}
		fmt.Print("Inserisci l'ultimo episodio da scaricare: ")
		ultimoStr, _ := reader.ReadString('\n')
		if i, err := strconv.Atoi(ultimoStr); err == nil {
			ultimo = i
		} else {
			ultimo = -1
		}

		// Getting the path where to store the downloads
		fmt.Printf("Inserisci il percorso dove salvare i file [Vuoto per: \"%s\"]: ", path)
		path, _ = reader.ReadString('\n')

		// Getting filenames
		fmt.Print("Inserisci il nome per i file: ")
		filename, _ = reader.ReadString('\n')
	} else if !noFlags() && (link == "D" || filename == "D") {
		panic("I flag --url e --dir sono obbligatori")
	}

	// Input formatting
	link = strings.TrimSpace(link)
	path = strings.TrimSpace(path)
	if !filepath.IsAbs(path) {
		path = filepath.Join(cwd, path)
	}
	if err := os.MkdirAll(path, 0777); err != nil {
		panic(fmt.Sprintf("Errore durante la creazione della directory `%s`: %s", path, err))
	}
	filename = strings.TrimSpace(filename)

	fmt.Printf("Scaricando da `%s`\n"+
		"da episodio %d a %d\n"+
		"in `%s`\n"+
		"con nome `%s{%d-%d}.mp4`\n"+
		"usando %d workers.\n",
		link, primo, ultimo, path, filename, primo, ultimo, workers,
	)

	var startTime = time.Now()
	run(link, primo, ultimo, path, filename, workers)
	fmt.Println("Download completati!")
	fmt.Printf("Tempo inpiegato: %s\n", time.Since(startTime).String())
}

func run(link string, primo int, ultimo int, path string, filename string, workers int) {
	// Setting up session
	fmt.Println("Inizializzando la sessione...")
	client := &http.Client{Jar: &cookieJar{make(map[string][]*http.Cookie)}}
	fmt.Println("Sessione creata!")

	// getting episode links
	fmt.Println("Cercando i link agli episodi...")
	var episodi [][]byte
	if e, err := getEpisodeLinks(client, link); err == nil {
		episodi = e
	} else {
		panic(fmt.Sprintf("Errore nello scraping dei link agli episodi: %s\n", err))
	}
	if bytes.Equal(episodi[0], []byte("NO EP 0")) && primo == 0 {
		primo = 1
	}
	if ultimo == -1 {
		ultimo = len(episodi) - 1
	}
	fmt.Println("Link agli episodi trovati!")

	// get stream links
	fmt.Println("Cercando i link alle stream...")
	var epLinks = []indexedUrl{}
	for i := primo; i <= ultimo; i++ {
		if indexedLink, err := getStreamLink(client, string(episodi[i]), i); err == nil {
			epLinks = append(epLinks, indexedLink)
		} else {
			panic(fmt.Sprintf("Errore nello scraping dei link alle stream: %s", err))
		}
	}
	fmt.Println("Link alle stream trovati!")

	// get file links
	fmt.Println("Cercando i link ai file...")
	var videoLinks = []indexedUrl{}
	for i := range epLinks {
		if indexedLink, err := getVideoLink(client, string(epLinks[i].u), epLinks[i].i); err == nil {
			videoLinks = append(videoLinks, indexedLink)
		} else {
			panic(fmt.Sprintf("Errore nello scraping dei link ai file: %s", err))
		}
	}
	fmt.Println("Link ai file trovati!")
	fmt.Println("Inizio download...")
	var m3u8Files []indexedUrl
	var mp4Files []indexedUrl
	for _, link := range videoLinks {
		if strings.HasSuffix(string(link.u), ".m3u8") {
			m3u8Files = append(m3u8Files, link)
		} else {
			mp4Files = append(mp4Files, link)
		}
	}
	if len(m3u8Files) > 0 {
		// new style downloads
		fmt.Println("Rilevati file m3u8! Inizializando il download tramite FFMPEG...")
		err := exec.Command("ffmpeg", "-version").Run()
		if err != nil {
			panic("Per questo tipo di file è necessario FFMPEG.")
		}

		jobs := make(chan indexedUrl, len(m3u8Files))
		wg := sync.WaitGroup{}
		wg.Add(len(m3u8Files))
		for range workers {
			go func() {
				defer wg.Done()
				downloader_m3u8(path, filename, jobs)
			}()
		}

		for _, link := range m3u8Files {
			jobs <- link
		}
		close(jobs) // no more urls, so tell workers to stop their loop
		wg.Wait()
	}
	if len(mp4Files) > 0 {
		// old style downloads
		jobs := make(chan indexedUrl, len(mp4Files))
		wg := sync.WaitGroup{}
		wg.Add(len(mp4Files))
		for range workers {
			go func() {
				defer wg.Done()
				downloader_mp4(client, path, filename, jobs)
			}()
		}

		// continually feed in urls to workers
		for _, link := range mp4Files {
			jobs <- link
		}
		close(jobs) // no more urls, so tell workers to stop their loop
		wg.Wait()
	}
}
