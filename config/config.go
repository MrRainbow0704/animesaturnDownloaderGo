package config

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"

	log "github.com/MrRainbow0704/animesaturnDownloaderGo/internal/logger"
)

type config struct {
	Verbose bool   `json:"verbose"`
	BaseURL string `json:"base_url"`
	NoCache bool   `json:"no_cache"`
}

var c config

const configPath string = "./config.json"

func Init() {
	f, err := os.Open(configPath)
	if errors.Is(err, fs.ErrNotExist) {
		f, err = os.Create(configPath)
		if err != nil {
			log.Fatalf("Impossible caricare il file di configurazione: %s", err)
		}
		defer f.Close()

		c = config{Verbose: false, BaseURL: "https://www.animesaturn.cx"}
		err = json.NewEncoder(f).Encode(c)
		if err != nil {
			log.Fatalf("Impossible decifrare il file di configurazione: %s", err)
		}
		return
	} else if err != nil {
		log.Fatalf("Impossible caricare il file di configurazione: %s", err)
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&c)
	if err != nil {
		log.Fatalf("Impossible decifrare il file di configurazione: %s", err)
	}
}

func Verbose() bool {
	return c.Verbose
}

func BaseURL() string {
	return c.BaseURL
}

func NoCache() bool {
	return c.NoCache
}

func SetVerbose(v bool) error {
	f, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Impossible caricare il file di configurazione: %s", err)
	}
	defer f.Close()

	c.Verbose = v
	err = json.NewEncoder(f).Encode(c)
	if err != nil {
		log.Fatalf("Impossible decifrare il file di configurazione: %s", err)
	}
	return nil
}

func SetBaseURL(v string) error {
	f, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Impossible caricare il file di configurazione: %s", err)
	}
	defer f.Close()

	c.BaseURL = v
	err = json.NewEncoder(f).Encode(c)
	if err != nil {
		log.Fatalf("Impossible decifrare il file di configurazione: %s", err)
	}
	return nil
}

func SetNoCache(v bool) error {
	f, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Impossible caricare il file di configurazione: %s", err)
	}
	defer f.Close()

	c.NoCache = v
	err = json.NewEncoder(f).Encode(c)
	if err != nil {
		log.Fatalf("Impossible decifrare il file di configurazione: %s", err)
	}
	return nil
}
