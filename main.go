package main

import (
	"bufio"
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

type IndexedUrl struct {
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
	// initializing reader
	reader := bufio.NewReader(os.Stdin)

	// Getting the page link
	fmt.Print("Inserisci il link alla pagina dell'anime: ")
	link, _ := reader.ReadString('\n')
	link = strings.TrimSpace(link)

	// Getting the first and last episodes to download
	fmt.Print("Inserisci il primo episodio da scaricare: ")
	primoStr, _ := reader.ReadString('\n')
	var primo int
	if i, err := strconv.Atoi(primoStr); err == nil {
		primo = i
	} else {
		primo = 0
	}
	fmt.Print("Inserisci l'ultimo episodio da scaricare: ")
	ultimoStr, _ := reader.ReadString('\n')
	var ultimo int
	if i, err := strconv.Atoi(ultimoStr); err == nil {
		ultimo = i
	} else {
		ultimo = -1
	}

	// Getting the path where to store the downloads
	path, _ := os.Getwd()
	fmt.Printf("Inserisci il percorso dove salvare i file [Vuoto per: \"%s\"]: ", path)
	pathNew, _ := reader.ReadString('\n')
	pathNew = strings.TrimSpace(pathNew)
	if filepath.IsAbs(pathNew) {
		path = pathNew
	} else {
		path = filepath.Join(path, pathNew)
	}
	if err := os.MkdirAll(path, 0777); err != nil {
		log.Panicf("Errore durante la creazione della directory: %s", err)
	}

	// Getting filenames
	fmt.Print("Inserisci il nome per i file: ")
	filename, _ := reader.ReadString('\n')
	filename = strings.TrimSpace(filename)

	var startTime = time.Now()
	// Setting up session
	log.Print("Inizializzando la sessione...")
	client := &http.Client{Jar: &myjar{make(map[string][]*http.Cookie)}}
	log.Print("Sessione creata!")

	// getting episode links
	log.Print("Cercando i link agli episodi...")
	var episodi [][]byte
	if e, err := GetEpisodeLinks(client, link); err == nil {
		episodi = e
	} else {
		log.Panicf("Errore nello scraping dei link agli episodi: %s", err)
	}
	if bytes.Equal(episodi[0], []byte("NOT FOUND")) && primo == 0 {
		primo = 1
	}
	log.Print("Link agli episodi trovati!")

	// get stream links
	log.Print("Cercando i link alle stream...")
	var epLinks = []IndexedUrl{}
	for i := primo; i <= ultimo; i++ {
		if indexedLink, err := GetStreamLink(client, string(episodi[i]), i); err == nil {
			epLinks = append(epLinks, indexedLink)
		} else {
			log.Panicf("Errore nello scraping dei link alle stream: %s", err)
		}
	}
	log.Print("Link alle stream trovati!")

	// get file links
	log.Print("Cercando i link ai file...")
	var videoLinks = []IndexedUrl{}
	for i := 0; i < len(epLinks); i++ {
		if indexedLink, err := GetVideoLink(client, string(epLinks[i].u), epLinks[i].i); err == nil {
			videoLinks = append(videoLinks, indexedLink)
		} else {
			log.Panicf("Errore nello scraping dei link ai file: %s", err)
		}
	}
	log.Print("Link ai file trovati!")

	// downloads
	log.Print("Inizio download...")
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
	log.Print("Download completati!")
	log.Printf("Tempo inpiegato: %s", time.Since(startTime).String())
}
