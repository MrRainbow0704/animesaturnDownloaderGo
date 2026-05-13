package main

import (
	"flag"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/MrRainbow0704/animesaturnDownloaderGo/internal/cache"
	"github.com/MrRainbow0704/animesaturnDownloaderGo/internal/helper"
	log "github.com/MrRainbow0704/animesaturnDownloaderGo/internal/logger"
	"github.com/MrRainbow0704/animesaturnDownloaderGo/internal/version"
)

var (
	execName   = filepath.Base(os.Args[0])
	localCache bool
	ver        bool
)

func initRoot() {
	flag.Usage = func() {
		log.Print(header + `Questa schermata di aiuto è divisa in più parti, usa "` + execName + ` <sottocomando> -h" per vedere la schermata di aiuto per il sottocomando specifico.

I sottocomandi disponibili sono:
  download		Scarica gli episodi di un anime
  search		Cerca un anime per nome

Utilizzo: ` + execName + ` <sottocomando> [opzioni]

Flag globali:
  -h, --help			stampa le informazioni di aiuto
  -v, --verbose			stampa altre informazioni di debug
  -V, --version			stampa la versione del programma e termina il programma
  --max-retry <numero>		numero massimo di tentativi per ogni richiesta HTTP	[default: 3]
  --no-cache			non utilizza la cache
  --local-cache			forza il programma ad utilizzare una cache locale
`)
	}

	flag.BoolVar(&log.Verbose, "verbose", false, "stampa altre informazioni di debug")
	flag.BoolVar(&log.Verbose, "v", false, "stampa altre informazioni di debug")
	flag.BoolVar(&localCache, "local-cache", false, "force the program to use local cache")
	flag.BoolVar(&ver, "version", false, "stampa la versione del programma")
	flag.BoolVar(&ver, "V", false, "stampa la versione del programma")
	flag.BoolVar(&cache.NoCachce, "no-cache", false, "non utilizza la cache")
	flag.IntVar(&helper.MaxRetry, "max-retry", 3, "numero massimo di tentativi per ogni richiesta HTTP")

	initDownload()
	initSearch()
}

func parseRoot(arguments []string) string {
	flag.Parse()

	if ver {
		return ""
	}

	dashlessArgs := flag.Args()
	if len(dashlessArgs) == 0 {
		log.Fatal("Nessun sottocomando specificato.\nUsa \"" + execName + " -h\" per vedere la schermata di aiuto.")
	}

	subcommand := dashlessArgs[0]
	arguments = slices.DeleteFunc(arguments, func(s string) bool {
		return s == subcommand
	})

	switch subcommand {
	case "download":
		parseDownload(arguments)
	case "search":
		parseSearch(arguments)
	default:
		log.Fatalf("Sottocomando `%s` non riconosciuto.\n", subcommand)
	}
	return subcommand
}

func runRoot(subcommand string) {
	startTime := time.Now()
	if ver {
		log.Printf("AnimesaturnDownloaderGo %s", version.Get())
		return
	}

	switch subcommand {
	case "download":
		runDownload()
	case "search":
		runSearch()
	default:
		log.Fatalf("Sottocomando `%s` non riconosciuto.\n", subcommand)
	}
	log.Infof("Tempo inpiegato: %s\n", time.Since(startTime).String())
}
