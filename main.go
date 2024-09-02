package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
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

type myjar struct {
	jar map[string][]*http.Cookie
}

type IndexedUrl struct {
	i int
	u []byte
}

func (p *myjar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	p.jar[u.Host] = cookies
}

func (p *myjar) Cookies(u *url.URL) []*http.Cookie {
	return p.jar[u.Host]
}

func GetEpisodeLinks(c *http.Client, u string) ([][]byte, error) {
	linkRegexp := regexp.MustCompile("(?i)<a[^>]*class=[\"'][^>]*bottone-ep[^>]*[\"'][^>]*[^>]*>")
	hrefRegexp := regexp2.MustCompile("(?i)(?<=href=[\"']).*?(?=[\"'])", 0)
	zeroEpRegexp := regexp.MustCompile("(?i)[^0-9A-Za-z]*ep-0[^0-9A-Za-z]*")

	req, _ := http.NewRequest("GET", u, nil)
	res, err := c.Do(req)
	if err != nil || res.StatusCode != 200 {
		log.Printf("Errore: %s", err)
		log.Printf("Status: %s", res.Status)
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Errore: %s", err)
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

func GetStreamLink(c *http.Client, u string, i int) (IndexedUrl, error) {
	linkRegexp := regexp.MustCompile("(?i)<a[^>]*href=[\"'][^>]*watch\\?[^>]*[\"'][^>]*>")
	hrefRegexp := regexp2.MustCompile("(?i)(?<=href=[\"']).*?(?=[\"'])", 0)

	req, _ := http.NewRequest("GET", u, nil)
	res, err := c.Do(req)
	if err != nil || res.StatusCode != 200 {
		log.Printf("Errore: %s", err)
		log.Printf("Status: %s", res.Status)
		return IndexedUrl{}, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Errore: %s", err)
		return IndexedUrl{}, err
	}
	content := strings.Replace(string(body), " ", "", -1)
	link := linkRegexp.FindAll([]byte(content), -1)[0]
	var streamLink string
	if match, err := hrefRegexp.FindStringMatch(string(link)); err == nil {
		streamLink = match.String()
	}

	return IndexedUrl{i, []byte(streamLink)}, nil
}

func GetVideoLink(c *http.Client, u string, i int) (IndexedUrl, error) {
	sourceRegexp := regexp.MustCompile("<source[^>]*>")
	srcRegexp := regexp2.MustCompile("(?<=src=[\"']).*?(?=[\"'])", 0)

	req, _ := http.NewRequest("GET", u, nil)
	res, err := c.Do(req)
	if err != nil || res.StatusCode != 200 {
		log.Printf("Errore: %s", err)
		log.Printf("Status: %s", res.Status)
		return IndexedUrl{}, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Errore: %s", err)
		return IndexedUrl{}, err
	}
	content := strings.Replace(string(body), " ", "", -1)
	link := sourceRegexp.FindAll([]byte(content), -1)[0]
	vidLink, err := srcRegexp.FindStringMatch(string(link))
	if err != nil {
		log.Printf("Errore: %s", err)
		return IndexedUrl{}, err
	}

	return IndexedUrl{i, []byte(vidLink.String())}, nil
}

func downloadFile(filepath string, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
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

func Downloader(id int, c *http.Client, path string, filename string, jobs <-chan IndexedUrl, results chan<- int) {
	for j := range jobs {
		name := filepath.Join(path, filename+strconv.Itoa(j.i)+".mp4")
		startTime := time.Now()
		log.Printf("Inizio download di `%s`...", filename+strconv.Itoa(j.i)+".mp4")

		downloadFile(name, string(j.u))

		log.Printf("Finito di scaricare `%s` in %ss", name, time.Since(startTime).String())
		results <- 0 // flag that job is finished
	}
}

func main() {
	// Getting the page link
	fmt.Print("Inserisci il link alla pagina dell'anime: ")
	var link string
	_, err := fmt.Scan(&link)
	if err != nil {
		log.Printf("Errore durante la lettura: %s", err)
		return
	}

	// Getting the first and last episodes to download
	fmt.Print("Inserisci il primo episodio da scaricare: ")
	var primoStr string
	var primo int
	_, err = fmt.Scan(&primoStr)
	if err != nil {
		log.Printf("Errore durante la lettura: %s", err)
		return
	}
	if i, err := strconv.Atoi(primoStr); err == nil {
		primo = i
	} else {
		primo = 0
	}
	fmt.Print("Inserisci l'ultimo episodio da scaricare: ")
	var ultimoStr string
	var ultimo int
	_, err = fmt.Scan(&ultimoStr)
	if err != nil {
		log.Printf("Errore durante la lettura: %s", err)
		return
	}
	if i, err := strconv.Atoi(ultimoStr); err == nil {
		ultimo = i
	} else {
		ultimo = -1
	}

	// Getting the path where to store the downloads
	var path string
	if wd, err := os.Getwd(); err == nil {
		path = wd
	} else {
		log.Printf("Errore durante il controllo della workdir: %s", err)
		return
	}
	fmt.Printf("Inserisci il percorso dove salvare i file [Vuoto per: \"%s\"]: ", path)
	var pathNew string
	_, err = fmt.Scan(&pathNew)
	if err != nil {
		log.Printf("Errore durante la lettura: %s", err)
		return
	}
	if filepath.IsAbs(pathNew) {
		path = pathNew
	} else {
		path = filepath.Join(path, pathNew)
	}
	if err := os.MkdirAll(path, 0777); err != nil {
		log.Printf("Errore durante la creazione della directory: %s", err)
		return
	}

	// Getting filenames
	fmt.Print("Inserisci il nome per i file: ")
	var filename string
	_, err = fmt.Scan(&filename)
	if err != nil {
		log.Printf("Errore durante la lettura: %s", err)
		return
	}

	var startTime = time.Now()
	// Setting up session
	log.Print("Inizializzando la sessione...")
	client := &http.Client{Jar: &myjar{make(map[string][]*http.Cookie)}}
	log.Print("Sessione creata!")

	// getting episode links
	var episodi [][]byte
	if e, err := GetEpisodeLinks(client, link); err == nil {
		episodi = e
	} else {
		log.Printf("Errore nello scraping di link: %s", err)
	}
	if bytes.Equal(episodi[0], []byte("NOT FOUND")) && primo == 0 {
		primo = 1
	}

	// get stream links
	var epLinks = []IndexedUrl{}
	for i := primo; i <= ultimo; i++ {
		if indexedLink, err := GetStreamLink(client, string(episodi[i]), i); err == nil {
			epLinks = append(epLinks, indexedLink)
		} else {
			return
		}
	}

	var videoLinks = []IndexedUrl{}
	for i := 0; i < len(epLinks); i++ {
		if indexedLink, err := GetVideoLink(client, string(epLinks[i].u), epLinks[i].i); err == nil {
			videoLinks = append(videoLinks, indexedLink)
		} else {
			return
		}
	}

	numJobs := len(videoLinks)
	jobs := make(chan IndexedUrl, numJobs)
	results := make(chan int, numJobs)
	for w := 1; w <= 3; w++ { // only 3 workers, all blocked initially
		go Downloader(w, client, path, filename, jobs, results)
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

	log.Printf("Tempo inpiegato: %s", time.Since(startTime).String())
}
