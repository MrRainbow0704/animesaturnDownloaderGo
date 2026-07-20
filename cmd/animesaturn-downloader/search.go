package main

import (
	"flag"
	"strings"

	"github.com/MrRainbow0704/animesaturnDownloaderGo/internal/cache"
	"github.com/MrRainbow0704/animesaturnDownloaderGo/internal/helper"
	log "github.com/MrRainbow0704/animesaturnDownloaderGo/internal/logger"
)

var searchCommand = flag.NewFlagSet("search", flag.ExitOnError)

var (
	search string
	page   uint
	base   string
)

func initSearch() {
	searchCommand.Usage = func() {
		log.Print(header + `Schermata di aiuto per il sottocomando "search".

Utilizzo: ` + execName + ` search -s <stringa> [-p <numero>] [--base-url <url>]

Flag per il sottocomando "search":
  -s, --search <stringa>	nome dell'anime da cercare		[obbligatorio]
  -p, --page <numero>		pagina da cercare. 0 => tutte		[default: 0]
  --base-url <url>		url alla home di animesaturn		[obbligatorio]
`)
	}

	searchCommand.BoolVar(&log.Verbose, "verbose", false, "stampa altre informazioni di debug")
	searchCommand.BoolVar(&log.Verbose, "v", false, "stampa altre informazioni di debug")
	searchCommand.BoolVar(&localCache, "local-cache", false, "force the program to use local cache")
	searchCommand.BoolVar(&ver, "version", false, "stampa la versione del programma")
	searchCommand.BoolVar(&ver, "V", false, "stampa la versione del programma")
	searchCommand.BoolVar(&cache.NoCachce, "no-cache", false, "non utilizza la cache")
	searchCommand.IntVar(&helper.MaxRetry, "max-retry", 3, "numero massimo di tentativi per ogni richiesta HTTP")
	searchCommand.StringVar(&search, "search", "", "nome dell'anime da cercare")
	searchCommand.StringVar(&search, "s", "", "nome dell'anime da cercare")
	searchCommand.UintVar(&page, "page", 0, "pagina da cercare")
	searchCommand.UintVar(&page, "p", 0, "pagina da cercare")
	searchCommand.StringVar(&base, "base-url", "", "url alla home di animesaturn")
}

func parseSearch(arguments []string) {
	searchCommand.Parse(arguments)

	search = strings.TrimSpace(search)
	if search == "" {
		log.Fatal("Il flag --search (-s) è obbligatorio")
	}

	if base != "" {
		helper.BaseURL = strings.Trim(base, "/")
	}
	cache.Init(localCache)
}

func runSearch() {
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
		log.Printf("%s\n", a.Title)
		log.Printf("  Url: %s\n", a.Url)
	}
}
