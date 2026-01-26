package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// GitHubReleasesURL is the API endpoint for releases
	GitHubReleasesURL = "https://api.github.com/repos/nerveband/beeper-api-cli/releases/latest"
	// CacheFileName is the name of the update cache file
	CacheFileName = "update-cache.json"
	// CacheDuration is how long to cache update checks
	CacheDuration = 24 * time.Hour
)

// GitHubRelease represents a GitHub release response
type GitHubRelease struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	PublishedAt time.Time `json:"published_at"`
	HTMLURL     string    `json:"html_url"`
}

// UpdateCache stores cached update information
type UpdateCache struct {
	LastCheck      time.Time `json:"last_check"`
	LatestVersion  string    `json:"latest_version"`
	ReleaseURL     string    `json:"release_url"`
	CurrentVersion string    `json:"current_version"`
}

// UpdateInfo contains information about an available update
type UpdateInfo struct {
	CurrentVersion string
	LatestVersion  string
	ReleaseURL     string
	UpdateAvailable bool
}

// getCacheDir returns the cache directory path
func getCacheDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".beeper-api-cli")
}

// getCachePath returns the full path to the cache file
func getCachePath() string {
	dir := getCacheDir()
	if dir == "" {
		return ""
	}
	return filepath.Join(dir, CacheFileName)
}

// loadCache loads the update cache from disk
func loadCache() (*UpdateCache, error) {
	cachePath := getCachePath()
	if cachePath == "" {
		return nil, fmt.Errorf("could not determine cache path")
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}

	var cache UpdateCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}

	return &cache, nil
}

// saveCache saves the update cache to disk
func saveCache(cache *UpdateCache) error {
	cachePath := getCachePath()
	if cachePath == "" {
		return fmt.Errorf("could not determine cache path")
	}

	// Ensure directory exists
	dir := filepath.Dir(cachePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0644)
}

// fetchLatestRelease fetches the latest release from GitHub
func fetchLatestRelease() (*GitHubRelease, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest("GET", GitHubReleasesURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "beeper-api-cli")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

// normalizeVersion removes 'v' prefix from version string
func normalizeVersion(version string) string {
	return strings.TrimPrefix(version, "v")
}

// compareVersions returns true if latest > current
func compareVersions(current, latest string) bool {
	current = normalizeVersion(current)
	latest = normalizeVersion(latest)

	// Skip comparison for dev versions
	if current == "dev" || current == "" {
		return false
	}

	// Simple string comparison works for semver
	return latest > current
}

// Check checks for updates and returns info about available updates
// Uses cached data if available and not expired
func Check(currentVersion string) (*UpdateInfo, error) {
	info := &UpdateInfo{
		CurrentVersion: currentVersion,
	}

	// Try to load cache first
	cache, err := loadCache()
	if err == nil && cache != nil {
		// Check if cache is still valid
		if time.Since(cache.LastCheck) < CacheDuration && cache.CurrentVersion == currentVersion {
			info.LatestVersion = cache.LatestVersion
			info.ReleaseURL = cache.ReleaseURL
			info.UpdateAvailable = compareVersions(currentVersion, cache.LatestVersion)
			return info, nil
		}
	}

	// Fetch latest release
	release, err := fetchLatestRelease()
	if err != nil {
		return info, err
	}

	info.LatestVersion = normalizeVersion(release.TagName)
	info.ReleaseURL = release.HTMLURL
	info.UpdateAvailable = compareVersions(currentVersion, info.LatestVersion)

	// Save to cache
	newCache := &UpdateCache{
		LastCheck:      time.Now(),
		LatestVersion:  info.LatestVersion,
		ReleaseURL:     info.ReleaseURL,
		CurrentVersion: currentVersion,
	}
	_ = saveCache(newCache) // Ignore save errors

	return info, nil
}

// CheckAsync performs an update check in a goroutine and returns results via channel
func CheckAsync(currentVersion string) <-chan *UpdateInfo {
	ch := make(chan *UpdateInfo, 1)
	go func() {
		info, _ := Check(currentVersion)
		ch <- info
		close(ch)
	}()
	return ch
}

// FormatUpdateNotice formats an update notification message
func FormatUpdateNotice(info *UpdateInfo) string {
	if info == nil || !info.UpdateAvailable {
		return ""
	}
	return fmt.Sprintf(
		"\nUpdate available: %s -> %s\nRun 'beeper upgrade' to update, or visit:\n%s\n",
		info.CurrentVersion,
		info.LatestVersion,
		info.ReleaseURL,
	)
}
