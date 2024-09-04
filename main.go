package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dlclark/regexp2"
)

type indexedUrl struct {
	i int
	u []byte
}

type myjar struct {
	jar map[string][]*http.Cookie
}

func (p *myjar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	p.jar[u.Host] = cookies
}

func (p *myjar) Cookies(u *url.URL) []*http.Cookie {
	return p.jar[u.Host]
}

func getEpisodeLinks(c *http.Client, u string) ([][]byte, error) {
	linkRegexp := regexp.MustCompile("(?i)<a[^>]*class=[\"'][^>]*bottone-ep[^>]*[\"'][^>]*[^>]*>")
	hrefRegexp := regexp2.MustCompile("(?i)(?<=href=[\"']).*?(?=[\"'])", 0)
	zeroEpRegexp := regexp.MustCompile("(?i)[^0-9A-Za-z]*ep-0[^0-9A-Za-z]*")

	req, _ := http.NewRequest("GET", u, nil)
	res, err := c.Do(req)
	if err != nil || res.StatusCode != 200 {
		fmt.Printf("Errore: %s", err)
		fmt.Printf("Status: %s", res.Status)
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Errore: %s", err)
		return nil, err
	}
	content := strings.Replace(string(body), " ", "", -1)
	links := linkRegexp.FindAll([]byte(content), -1)
	linksList := [][]byte{}
	for i := 0; i < len(links); i++ {
		if match, err := hrefRegexp.FindStringMatch(string(links[i])); err == nil {
			linksList = append(linksList, []byte(match.String()))
		}
	}
	if zeroEpRegexp.FindAll(linksList[0], -1) != nil {
		return linksList, nil
	}
	return append([][]byte{[]byte("EP 0 NOT FOUND")}, linksList...), nil
}

func getStreamLink(c *http.Client, u string, i int) (indexedUrl, error) {
	linkRegexp := regexp.MustCompile("(?i)<a[^>]*href=[\"'][^>]*watch\\?[^>]*[\"'][^>]*>")
	hrefRegexp := regexp2.MustCompile("(?i)(?<=href=[\"']).*?(?=[\"'])", 0)

	req, _ := http.NewRequest("GET", u, nil)
	res, err := c.Do(req)
	if err != nil || res.StatusCode != 200 {
		fmt.Printf("Errore: %s", err)
		fmt.Printf("Status: %s", res.Status)
		return indexedUrl{}, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Errore: %s", err)
		return indexedUrl{}, err
	}
	content := strings.Replace(string(body), " ", "", -1)
	link := linkRegexp.FindAll([]byte(content), -1)[0]
	var streamLink string
	if match, err := hrefRegexp.FindStringMatch(string(link)); err == nil {
		streamLink = match.String()
	}

	return indexedUrl{i, []byte(streamLink)}, nil
}

func getVideoLink(c *http.Client, u string, i int) (indexedUrl, error) {
	sourceRegexp := regexp.MustCompile("<source[^>]*>")
	srcRegexp := regexp2.MustCompile("(?<=src=[\"']).*?(?=[\"'])", 0)

	req, _ := http.NewRequest("GET", u, nil)
	res, err := c.Do(req)
	if err != nil || res.StatusCode != 200 {
		fmt.Printf("Errore: %s", err)
		fmt.Printf("Status: %s", res.Status)
		return indexedUrl{}, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Errore: %s", err)
		return indexedUrl{}, err
	}
	content := strings.Replace(string(body), " ", "", -1)
	link := sourceRegexp.FindAll([]byte(content), -1)[0]
	vidLink, err := srcRegexp.FindStringMatch(string(link))
	if err != nil {
		fmt.Printf("Errore: %s", err)
		return indexedUrl{}, err
	}

	return indexedUrl{i, []byte(vidLink.String())}, nil
}

func downloadFile(c *http.Client, filepath string, url string) error {
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
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func downloader(c *http.Client, path string, filename string, jobs <-chan indexedUrl, results chan<- int) {
	for j := range jobs {
		name := filepath.Join(path, filename+strconv.Itoa(j.i)+".mp4")
		startTime := time.Now()
		fmt.Printf("Inizio download di `%s`...", filename+strconv.Itoa(j.i)+".mp4")

		downloadFile(c, name, string(j.u))

		fmt.Printf("Finito di scaricare `%s` in %ss", name, time.Since(startTime).String())
		results <- 0 // flag that job is finished
	}
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
  -u, --url 		link alla pagina dell'anime 	[Obbligatorio]
  -f, --first		primo episodio da scaricare 	[Default: 0]
  -l, --last		ultimo episodio da scaricare 	[Default: -1]
  -d, --dir 		percorso dove salvare i file 	[Default: percorso corrente]
  -n, --filename	nome del file senza estensione	[Obbligatorio]
`

	var link string
	flag.StringVar(&link, "url", "D", "link alla pagina dell'anime")
	flag.StringVar(&link, "u", "D", "link alla pagina dell'anime")
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

	var startTime = time.Now()
	// Setting up session
	fmt.Println("Inizializzando la sessione...")
	client := &http.Client{Jar: &myjar{make(map[string][]*http.Cookie)}}
	fmt.Println("Sessione creata!")

	// getting episode links
	fmt.Println("Cercando i link agli episodi...")
	var episodi [][]byte
	if e, err := getEpisodeLinks(client, link); err == nil {
		episodi = e
	} else {
		panic(fmt.Sprintf("Errore nello scraping dei link agli episodi: %s", err))
	}
	if bytes.Equal(episodi[0], []byte("EP 0 NOT FOUND")) && primo == 0 {
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
	for i := 0; i < len(epLinks); i++ {
		if indexedLink, err := getVideoLink(client, string(epLinks[i].u), epLinks[i].i); err == nil {
			videoLinks = append(videoLinks, indexedLink)
		} else {
			panic(fmt.Sprintf("Errore nello scraping dei link ai file: %s", err))
		}
	}
	fmt.Println("Link ai file trovati!")

	// downloads
	fmt.Println("Inizio download...")
	numJobs := len(videoLinks)
	jobs := make(chan indexedUrl, numJobs)
	results := make(chan int, numJobs)
	for i := 1; i <= 3; i++ { // only 3 workers, all blocked initially
		go downloader(client, path, filename, jobs, results)
	}

	// continually feed in urls to workers
	for _, link := range videoLinks {
		jobs <- link
	}
	close(jobs) // no more urls, so tell workers to stop their loop

	// needed if you want to make sure that workers don't block forever on writing results,
	// remove both this loop and workers writing results if you don't need output from workers
	for a := 1; a <= numJobs; a++ {
		<-results
	}
	fmt.Println("Download completati!")
	fmt.Printf("Tempo inpiegato: %s", time.Since(startTime).String())
}
