package cache

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/MrRainbow0704/animesaturnDownloaderGo/config"
	log "github.com/MrRainbow0704/animesaturnDownloaderGo/internal/logger"
)

var (
	MaxItems int           = config.CacheMaxItems()
	MaxTime  time.Duration = config.CacheMaxTime()
	NoCachce bool          = config.NoCache()
)

var cacheDir string

type cacheKey string

func (c cacheKey) String() string {
	return string(c)
}
func (c cacheKey) Set(v any) error {
	if NoCachce {
		return nil
	}
	return set(c.String(), v)
}
func (c cacheKey) Get(v any) error {
	if NoCachce {
		return errors.New("--no-cache is enabled")
	}
	return get(c.String(), v)
}
func (c cacheKey) Del() error {
	if NoCachce {
		return nil
	}
	return del(c.String())
}

func Key(v ...any) cacheKey {
	if NoCachce {
		return ""
	}
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		return ""
	}
	// return cacheKey(base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s(%#+v)", f.Name(), v))))
	hash := sha256.Sum256(fmt.Appendf([]byte{}, "%s(%#+v)", f.Name(), v))
	return cacheKey(fmt.Sprintf("%x", hash))
}

func Init() {
	if NoCachce {
		return
	}

	if userCache, err := os.UserCacheDir(); err != nil {
		cacheDir = "./.cache"
	} else {
		cacheDir = filepath.Join(userCache, "animesaturn-downloader/.cache")
	}

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		log.Fatalf("Impossibile creare la directory per la cache: %s", err)
	}
}

func cleaner() error {
	fs, err := os.ReadDir(cacheDir)
	if err != nil {
		log.Fatalf("Impossibile leggere la directory per la cache: %s", err)
		return err
	}
	for _, f := range fs {
		if i, err := f.Info(); err == nil && time.Since(i.ModTime()) >= MaxTime {
			del(i.Name())
		}
	}
	
	for len(fs) >= MaxItems {
		if err := os.Remove(oldest()); err != nil {
			log.Fatalf("Impossibile rimuovere il file di cache: %s", err)
			return err
		}
	}
	return nil
}

func oldest() string {
	fs, err := os.ReadDir(cacheDir)
	if err != nil {
		log.Fatalf("Impossibile leggere la directory per la cache: %s", err)
		return ""
	}

	if len(fs) == 0 {
		return ""
	}

	oldestTime := time.Now()
	var oldestFile os.DirEntry
	for _, f := range fs {
		if i, err := f.Info(); err == nil && !f.IsDir() && i.ModTime().Before(oldestTime) {
			oldestFile = f
			oldestTime = i.ModTime()
		}
	}
	return oldestFile.Name()
}

func set(n string, v any) error {
	cleaner()

	f, err := os.Create(file(n))
	if errors.Is(err, fs.ErrExist) {
		f, err = os.Open(file(n))
	}
	if err != nil {
		log.Fatalf("Inpossibile creare file di cache: %s", err)
		return err
	}

	if err := json.NewEncoder(f).Encode(v); err != nil {
		log.Fatalf("Impossibile scrivere nel file di cache: %s", err)
		del(n)
		return err
	}
	return nil
}

func get(n string, v any) error {
	f, err := os.Open(file(n))
	if err != nil {
		return err
	}

	if err := json.NewDecoder(f).Decode(v); err != nil {
		log.Fatalf("Impossibile scrivere nel file di cache: %s", err)
		del(n)
		return err
	}
	return nil
}

func del(n string) error {
	if err := os.Remove(file(n)); err != nil {
		log.Fatalf("Impossibile rimuovere il file di cache: %s", err)
		return err
	}
	return nil
}

func file(n string) string {
	return filepath.Join(cacheDir, n+".CACHE")
}
