package dbfetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"whereyouat/internal/env"
)

type GitHubRelease struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type DBFetcher struct {
	dbDir   string
	dbPath  string
	tagPath string
}

func New() *DBFetcher {
	cfg := env.Get()
	return &DBFetcher{
		dbDir:   cfg.DB_DIR,
		dbPath:  filepath.Join(cfg.DB_DIR, cfg.DB_FILENAME),
		tagPath: filepath.Join(cfg.DB_DIR, cfg.TAG_FILE),
	}
}

func (f *DBFetcher) EnsureDatabase() (string, error) {
	if err := os.MkdirAll(f.dbDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create db directory: %w", err)
	}

	release, err := f.getLatestRelease()
	if err != nil {
		log.Printf("Warning: failed to check for updates: %v", err)
		if _, err := os.Stat(f.dbPath); err == nil {
			log.Println("Using existing database")
			return f.dbPath, nil
		}
		return "", fmt.Errorf("no existing database and failed to fetch release info: %w", err)
	}

	needsDownload := f.needsDownload(release.TagName)

	if needsDownload {
		cfg := env.Get()
		downloadURL := ""
		for _, asset := range release.Assets {
			if asset.Name == cfg.DB_FILENAME {
				downloadURL = asset.BrowserDownloadURL
				break
			}
		}

		if downloadURL == "" {
			return "", fmt.Errorf("GeoLite2-Country.mmdb not found in release assets")
		}

		log.Printf("Downloading database from release %s...", release.TagName)
		if err := f.downloadDatabase(downloadURL, release.TagName); err != nil {
			return "", fmt.Errorf("failed to download database: %w", err)
		}
	} else {
		log.Println("Database is up-to-date")
	}

	return f.dbPath, nil
}

func (f *DBFetcher) getLatestRelease() (*GitHubRelease, error) {
	cfg := env.Get()
	resp, err := http.Get(cfg.MMDB_GITHUB_API_URL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse release info: %w", err)
	}

	return &release, nil
}

func (f *DBFetcher) needsDownload(latestTag string) bool {
	if _, err := os.Stat(f.dbPath); os.IsNotExist(err) {
		return true
	}

	storedTag, err := f.getStoredTag()
	if err != nil {
		return true
	}

	if storedTag != latestTag {
		log.Printf("New version available: %s (current: %s)", latestTag, storedTag)
		return true
	}

	return false
}

func (f *DBFetcher) downloadDatabase(url string, tag string) error {
	if _, err := os.Stat(f.dbPath); err == nil {
		log.Println("Removing old database...")
		if err := os.Remove(f.dbPath); err != nil {
			log.Printf("Warning: failed to remove old database: %v", err)
		}
	}

	log.Printf("Downloading from %s...", url)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download database: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with HTTP %d", resp.StatusCode)
	}

	tempFile := f.dbPath + ".tmp"
	out, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	_, err = io.Copy(out, resp.Body)
	out.Close()

	if err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to write database file: %w", err)
	}

	if err := os.Rename(tempFile, f.dbPath); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	if err := f.storeTag(tag); err != nil {
		log.Printf("Warning: failed to store tag: %v", err)
	}

	log.Printf("Database downloaded successfully (version: %s)", tag)
	return nil
}

func (f *DBFetcher) getStoredTag() (string, error) {
	data, err := os.ReadFile(f.tagPath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func (f *DBFetcher) storeTag(tag string) error {
	return os.WriteFile(f.tagPath, []byte(tag), 0644)
}
