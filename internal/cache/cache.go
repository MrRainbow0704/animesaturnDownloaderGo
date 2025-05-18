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

	log "github.com/MrRainbow0704/animesaturnDownloaderGo/internal/logger"
)

var (
	MaxItems = 500
	MaxTime  = time.Hour * 24
	cacheDir = "./.cache"
)

type cacheKey string

func (c cacheKey) String() string {
	return string(c)
}
func (c cacheKey) Set(v any) error {
	return set(c.String(), v)
}
func (c cacheKey) Get(v any) error {
	return get(c.String(), v)
}
func (c cacheKey) Del() error {
	return del(c.String())
}

func Key(v ...any) cacheKey {
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

func init() {
	if err := os.MkdirAll(cacheDir, 0x666); err != nil {
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
