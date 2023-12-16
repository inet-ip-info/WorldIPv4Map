package main

import (
	"crypto/sha1"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// cacheDir はキャッシュファイルを保存するディレクトリを指定します。
const cacheDir = "./cache"
const defaultTTL = 12 * time.Hour

// cacheEntry stores information about a cached response.
type cacheEntry struct {
	Expiration time.Time
	Path       string
}

// getCachePath calculates a file path for the cache of the given URL.
func getCachePath(url string) string {
	hash := sha1.Sum([]byte(url))
	return fmt.Sprintf("%s/%x", cacheDir, hash)
}

// saveCacheEntry saves a cache entry.
func saveCacheEntry(url string, expiration time.Time) error {
	entry := cacheEntry{Expiration: expiration, Path: getCachePath(url)}
	file, err := os.Create(getCacheEntryPath(url))
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(entry)
}

// loadCacheEntry loads a cache entry.
func loadCacheEntry(url string) (*cacheEntry, error) {
	file, err := os.Open(getCacheEntryPath(url))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entry cacheEntry
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&entry)
	return &entry, err
}

// cacheIsValid checks if the cache for the given URL is still valid.
func cacheIsValid(url string) bool {
	entry, err := loadCacheEntry(url)
	if err != nil {
		return false
	}
	return time.Now().Before(entry.Expiration)
}

// FetchData fetches data from the URL or the cache.
func OpenURLFile(url string) (io.ReadCloser, error) {
	cachePath := getCachePath(url)

	// Check cache validity
	if cacheIsValid(url) {
		return os.Open(cachePath)
	}

	// Fetch from the URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	expiration := parseCacheHeaders(resp.Header)
	if err := saveToCache(resp.Body, cachePath); err != nil {
		return nil, err
	}
	if err := saveCacheEntry(url, expiration); err != nil {
		return nil, err
	}

	return os.Open(cachePath)
}

func parseCacheHeaders(headers http.Header) time.Time {
	// Cache-Control ヘッダーを確認
	if cacheControl := headers.Get("Cache-Control"); cacheControl != "" {
		parts := strings.Split(cacheControl, ",")
		for _, part := range parts {
			if strings.HasPrefix(part, "max-age=") {
				ageStr := strings.TrimSpace(strings.TrimPrefix(part, "max-age="))
				if maxAge, err := strconv.Atoi(ageStr); err == nil {
					log.Printf("max-age=%d", maxAge)
					return time.Now().Add(time.Duration(maxAge) * time.Second)
				}
			}
		}
	}

	// Expires ヘッダーを確認
	if expires := headers.Get("Expires"); expires != "" {
		if expTime, err := http.ParseTime(expires); err == nil {
			log.Printf("Expires=%v", expTime)
			return expTime
		}
	}

	// ヘッダーがない場合のデフォルトの有効期限
	log.Printf("default Expire=%v", defaultTTL)
	return time.Now().Add(24 * time.Hour)
}

func saveToCache(body io.Reader, filePath string) error {
	// キャッシュディレクトリが存在しない場合は作成
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	// ファイルを開く（なければ作成）
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// データをファイルに書き込む
	_, err = io.Copy(file, body)
	return err
}

// getCacheEntryPath returns a file path for the cache entry of the given URL.
func getCacheEntryPath(url string) string {
	hash := sha1.Sum([]byte(url))
	cacheEntryFile := fmt.Sprintf("%x.cache", hash[:]) // Convert the hash to a hex string and append .cache
	return fmt.Sprintf("%s/%s", cacheDir, cacheEntryFile)
}
