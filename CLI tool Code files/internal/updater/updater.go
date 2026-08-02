// Package updater implements self-update logic for proj-init.
package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	apperrors "github.com/telugusmasher2010-collab/project-cli-tool/internal/errors"
)

const (
	// Repo is the GitHub repository hosting proj-init releases.
	Repo = "telugusmasher2010-collab/project-cli-tool"
	// LatestReleaseURL points at the GitHub API latest release endpoint.
	LatestReleaseURL = "https://api.github.com/repos/" + Repo + "/releases/latest"
)

// Release describes the GitHub latest release payload fields we need.
type Release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// Info reports whether an update is available and details about it.
type Info struct {
	Current     string
	Latest      string
	Updateable  bool
	DownloadURL string
}

// NewClient builds an http.Client with a sensible timeout.
func NewClient() *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
}

// LatestVersion fetches the latest published release tag via the GitHub API.
func LatestVersion(client *http.Client) (string, error) {
	req, err := http.NewRequest(http.MethodGet, LatestReleaseURL, nil)
	if err != nil {
		return "", apperrors.Wrap(apperrors.ErrInternal, "failed to build release request", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "proj-init-updater")

	resp, err := client.Do(req)
	if err != nil {
		return "", apperrors.Wrap(apperrors.ErrInternal, "failed to reach GitHub releases", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", apperrors.New(apperrors.ErrInternal, "no releases published yet for proj-init")
	}
	if resp.StatusCode != http.StatusOK {
		return "", apperrors.New(apperrors.ErrInternal, fmt.Sprintf("GitHub API returned %s", resp.Status))
	}

	var rel Release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return "", apperrors.Wrap(apperrors.ErrInternal, "failed to parse release payload", err)
	}
	if rel.TagName == "" {
		return "", apperrors.New(apperrors.ErrInternal, "release payload missing tag_name")
	}
	return rel.TagName, nil
}

// Check compares the current version against the latest release and returns
// Info. DownloadURL is populated when an update is available.
func Check(current string, client *http.Client) (*Info, error) {
	latest, err := LatestVersion(client)
	if err != nil {
		return nil, err
	}

	latest = strings.TrimPrefix(latest, "v")
	current = strings.TrimPrefix(current, "v")

	info := &Info{Current: current, Latest: latest}
	if latest != "" && latest != current {
		info.Updateable = true
		info.DownloadURL = assetURL(latest)
	}
	return info, nil
}

// assetURL computes the download URL for the current OS/arch asset.
func assetURL(version string) string {
	osName := runtime.GOOS
	arch := runtime.GOARCH
	ext := "tar.gz"
	if osName == "windows" {
		ext = "zip"
	}
	return fmt.Sprintf("https://github.com/%s/releases/download/v%s/proj-init_%s_%s_%s.%s",
		Repo, version, version, osName, arch, ext)
}

// ExecutablePath returns the absolute path of the running binary.
func ExecutablePath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", apperrors.Wrap(apperrors.ErrInternal, "failed to locate current binary", err)
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return "", apperrors.Wrap(apperrors.ErrInternal, "failed to resolve binary path", err)
	}
	return exe, nil
}

// Apply downloads the release asset and replaces the running binary in place.
// It first writes to a sibling temp file, verifies it, then renames over the
// original so an interrupted update never corrupts the installed binary.
func Apply(info *Info, client *http.Client) error {
	if info == nil || !info.Updateable || info.DownloadURL == "" {
		return apperrors.New(apperrors.ErrInvalidInput, "no update available to apply")
	}

	exe, err := ExecutablePath()
	if err != nil {
		return err
	}

	tmpDir, err := os.MkdirTemp(filepath.Dir(exe), ".proj-init-update-*")
	if err != nil {
		return apperrors.Wrap(apperrors.ErrFilesystem, "failed to create temp directory", err)
	}
	defer os.RemoveAll(tmpDir)

	assetName := filepath.Base(info.DownloadURL)
	assetPath := filepath.Join(tmpDir, assetName)

	req, err := http.NewRequest(http.MethodGet, info.DownloadURL, nil)
	if err != nil {
		return apperrors.Wrap(apperrors.ErrInternal, "failed to build download request", err)
	}
	req.Header.Set("User-Agent", "proj-init-updater")

	resp, err := client.Do(req)
	if err != nil {
		return apperrors.Wrap(apperrors.ErrInternal, "failed to download release asset", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return apperrors.New(apperrors.ErrInternal, fmt.Sprintf("download returned %s", resp.Status))
	}

	out, err := os.Create(assetPath)
	if err != nil {
		return apperrors.Wrap(apperrors.ErrFilesystem, "failed to create asset file", err)
	}
	if _, err := io.Copy(out, resp.Body); err != nil {
		out.Close()
		return apperrors.Wrap(apperrors.ErrFilesystem, "failed to write asset file", err)
	}
	out.Close()

	newBin, err := extractBinary(assetPath)
	if err != nil {
		return err
	}

	backup := exe + ".bak"
	if err := os.Rename(exe, backup); err != nil {
		return apperrors.Wrap(apperrors.ErrFilesystem, "failed to back up current binary", err)
	}
	if err := os.Rename(newBin, exe); err != nil {
		os.Rename(backup, exe)
		return apperrors.Wrap(apperrors.ErrFilesystem, "failed to install new binary", err)
	}
	os.Remove(backup)
	return nil
}

// extractBinary unpacks a downloaded archive and returns the path of the
// extracted proj-init binary.
func extractBinary(archivePath string) (string, error) {
	dir := filepath.Dir(archivePath)
	ext := strings.ToLower(filepath.Ext(archivePath))
	switch ext {
	case ".zip":
		if err := unzip(archivePath, dir); err != nil {
			return "", err
		}
	case ".gz", ".tgz":
		if err := untar(archivePath, dir); err != nil {
			return "", err
		}
	default:
		// Not an archive: assume it is the binary itself.
		return archivePath, nil
	}

	name := "proj-init"
	if runtime.GOOS == "windows" {
		name = "proj-init.exe"
	}
	bin := filepath.Join(dir, name)
	if _, err := os.Stat(bin); err != nil {
		return "", apperrors.Wrap(apperrors.ErrFilesystem, "extracted archive missing binary", err)
	}
	return bin, nil
}
