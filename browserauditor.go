package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

const (
	toolName    = "BrowserAuditor"
	toolVersion = "1.0"
)

type RecoverySnapshot struct {
	ToolName     string             `json:"tool_name"`
	ToolVersion  string             `json:"tool_version"`
	GeneratedAt  string             `json:"generated_at"`
	MachineID    string             `json:"machine_id"`
	SystemInfo   SystemInfo         `json:"system_info"`
	BrowserAudit []BrowserInventory `json:"browser_audit"`
	Notes        []string           `json:"notes"`
}

type SystemInfo struct {
	OS       string   `json:"os"`
	Arch     string   `json:"arch"`
	Hostname string   `json:"hostname"`
	Username string   `json:"username"`
	HomeDir  string   `json:"home_dir"`
	LocalIPs []string `json:"local_ips"`
}

func main() {
	outputPath := flag.String("out", "", "path to the output JSON file")
	flag.Parse()

	now := time.Now().UTC()
	resolvedOutputPath := *outputPath
	if resolvedOutputPath == "" {
		resolvedOutputPath = defaultOutputPath(now)
	}

	snapshot := RecoverySnapshot{
		ToolName:     toolName,
		ToolVersion:  toolVersion,
		GeneratedAt:  now.Format(time.RFC3339),
		MachineID:    newMachineID(),
		SystemInfo:   collectSystemInfo(),
		BrowserAudit: collectBrowserInventories(),
		Notes: []string{
			"Local-only recovery snapshot.",
			"No browser secrets, cookies, history contents, or passwords are collected.",
			"Browser entries only report profile directories and known artifact file locations.",
		},
	}

	if err := saveSnapshot(resolvedOutputPath, snapshot); err != nil {
		fmt.Fprintf(os.Stderr, "error saving snapshot: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s %s wrote %s with %d browser entries\n", toolName, toolVersion, resolvedOutputPath, len(snapshot.BrowserAudit))
}

func newMachineID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "unknown-machine"
	}

	return hex.EncodeToString(buf)
}

func saveSnapshot(path string, data RecoverySnapshot) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, jsonData, 0644)
}

func defaultOutputPath(now time.Time) string {
	return fmt.Sprintf("browser-audit-report-%s.json", now.Format("20060102-150405Z"))
}
