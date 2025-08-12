package config

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	log "github.com/MrRainbow0704/animesaturnDownloaderGo/internal/logger"
)

type config struct {
	Verbose       bool   `json:"verbose"`
	NoCache       bool   `json:"no_cache"`
	BaseURL       string `json:"base_url"`
	CacheMaxItems int    `json:"cache_max_items"`
	CacheMaxTime  int    `json:"cache_max_time"`
}

var c = config{
	Verbose:       false,
	BaseURL:       "https://www.animesaturn.cx",
	NoCache:       false,
	CacheMaxItems: 500,
	CacheMaxTime:  int(time.Hour * 24),
}

var configPath string
var configDir string

func Init(local bool) {
	if userConfig, err := os.UserConfigDir(); err != nil || local {
		configDir = "."
	} else {
		configDir = filepath.Join(userConfig, "animesaturn-downloader")
	}
	configPath = filepath.Join(configDir, "config.json")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Fatalf("Impossibile creare la directory per la cache: %s", err)
	}

	f, err := os.Open(configPath)
	if errors.Is(err, fs.ErrNotExist) {
		f, err = os.Create(configPath)
		if err != nil {
			log.Fatalf("Impossible caricare il file di configurazione: %s", err)
		}
		defer f.Close()

		if err := json.NewEncoder(f).Encode(c); err != nil {
			log.Fatalf("Impossible decifrare il file di configurazione: %s", err)
		}
		return
	} else if err != nil {
		log.Fatalf("Impossible caricare il file di configurazione: %s", err)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&c); err != nil {
		log.Fatalf("Impossible decifrare il file di configurazione: %s. Eliminalo e riprova.", err)
	}
}

func Verbose() bool {
	return c.Verbose
}
func NoCache() bool {
	return c.NoCache
}
func BaseURL() string {
	return c.BaseURL
}
func CacheMaxItems() int {
	return c.CacheMaxItems
}
func CacheMaxTime() time.Duration {
	return time.Duration(c.CacheMaxTime)
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
func SetCacheMaxItems(v int) error {
	f, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Impossible caricare il file di configurazione: %s", err)
	}
	defer f.Close()

	c.CacheMaxItems = v
	err = json.NewEncoder(f).Encode(c)
	if err != nil {
		log.Fatalf("Impossible decifrare il file di configurazione: %s", err)
	}
	return nil
}
func SetCacheMaxTime(v time.Duration) error {
	f, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Impossible caricare il file di configurazione: %s", err)
	}
	defer f.Close()

	c.CacheMaxTime = int(v)
	err = json.NewEncoder(f).Encode(c)
	if err != nil {
		log.Fatalf("Impossible decifrare il file di configurazione: %s", err)
	}
	return nil
}
