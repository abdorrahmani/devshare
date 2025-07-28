package cmd

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update DevShare to the latest version",
	Long:  `Update DevShare to the latest version available.`,
	Run: func(cmd *cobra.Command, args []string) {
		const repo = "abdorrahmani/devshare"
		const apiURL = "https://api.github.com/repos/" + repo + "/releases/latest"
		client := &http.Client{Timeout: 15 * time.Second}
		resp, err := client.Get(apiURL)
		if err != nil {
			cmd.Println("Failed to check for updates:", err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			cmd.Println("Failed to fetch release info. Status:", resp.Status)
			return
		}
		var release struct {
			TagName string `json:"tag_name"`
			Assets  []struct {
				Name               string `json:"name"`
				BrowserDownloadURL string `json:"browser_download_url"`
			} `json:"assets"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
			cmd.Println("Failed to parse release info:", err)
			return
		}
		current := Version
		latest := release.TagName
		latest = strings.TrimPrefix(latest, "v")
		current = strings.TrimPrefix(current, "v")
		if current == latest {
			cmd.Println("You are already running the latest version (", Version, ")!")
			return
		}

		if current > latest {
			cmd.Println("🎉 You are running a newer version than the latest release. This is unexpected.")
			cmd.Printf("✅ Current version: %s, Latest version: %s\n", current, latest)
			return
		}

		cmd.Printf("New version available: %s (current: %s)\n", release.TagName, Version)
		osName := runtime.GOOS
		arch := runtime.GOARCH
		var archStr string
		switch arch {
		case "amd64":
			archStr = "x86_64"
		case "386":
			archStr = "i386"
		default:
			archStr = arch
		}
		var archiveExt, archiveFormat string
		var assetName string
		projectName := "DevShare"
		osTitle := cases.Title(language.English).String(osName)
		if osName == "windows" {
			archiveExt = ".zip"
			archiveFormat = "zip"
			assetName = fmt.Sprintf("%s_%s_%s%s", projectName, osTitle, archStr, archiveExt)
		} else {
			archiveExt = ".tar.gz"
			archiveFormat = "tar.gz"
			assetName = fmt.Sprintf("%s_%s_%s%s", projectName, osTitle, archStr, archiveExt)
		}
		var downloadURL string
		for _, asset := range release.Assets {
			if asset.Name == assetName {
				downloadURL = asset.BrowserDownloadURL
				break
			}
		}
		if downloadURL == "" {
			cmd.Printf("No archive found for your OS/arch: %s/%s\n", osName, arch)
			return
		}
		tmpArchive, err := os.CreateTemp("", assetName)
		if err != nil {
			cmd.Println("Failed to create temp file:", err)
			return
		}
		defer os.Remove(tmpArchive.Name())
		cmd.Println("Downloading:", downloadURL)
		resp, err = client.Get(downloadURL)
		if err != nil {
			cmd.Println("Download failed:", err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			cmd.Println("Failed to download archive. Status:", resp.Status)
			return
		}
		if _, err := io.Copy(tmpArchive, resp.Body); err != nil {
			cmd.Println("Failed to save archive:", err)
			return
		}
		if err := tmpArchive.Close(); err != nil {
			cmd.Println("Failed to close temp file:", err)
			return
		}

		var binName string
		if osName == "windows" {
			binName = "devshare.exe"
		} else {
			binName = "devshare"
		}
		tmpBin, err := os.CreateTemp("", binName)
		if err != nil {
			cmd.Println("Failed to create temp binary file:", err)
			return
		}
		defer os.Remove(tmpBin.Name())
		if archiveFormat == "zip" {
			if err := extractFromZip(tmpArchive.Name(), binName, tmpBin); err != nil {
				cmd.Println("Failed to extract binary from zip:", err)
				return
			}
		} else {
			if err := extractFromTarGz(tmpArchive.Name(), binName, tmpBin); err != nil {
				cmd.Println("Failed to extract binary from tar.gz:", err)
				return
			}
		}
		if err := tmpBin.Close(); err != nil {
			cmd.Println("Failed to close extracted binary file:", err)
			return
		}
		if osName != "windows" {
			if err := os.Chmod(tmpBin.Name(), 0755); err != nil {
				cmd.Println("Failed to set permissions:", err)
				return
			}
		}
		if osName == "windows" {
			// Run install.bat in the directory containing the extracted binary
			binDir := filepath.Dir(tmpBin.Name())
			installScript := filepath.Join(binDir, "install.bat")
			if _, err := os.Stat(installScript); os.IsNotExist(err) {
				cmd.Println("install.bat not found in extracted archive. Please install manually.")
				return
			}
			updateCmd := exec.Command("cmd", "/C", installScript)
			updateCmd.Dir = binDir
			updateCmd.Stdout = os.Stdout
			updateCmd.Stderr = os.Stderr
			cmd.Println("Running installer...")
			if err := updateCmd.Run(); err != nil {
				cmd.Println("Failed to run install.bat. Try running as Administrator or install manually.")
				return
			}
		} else {
			// Run install.sh in the directory containing the extracted binary
			binDir := filepath.Dir(tmpBin.Name())
			installScript := filepath.Join(binDir, "install.sh")
			if _, err := os.Stat(installScript); os.IsNotExist(err) {
				cmd.Println("install.sh not found in extracted archive. Please install manually.")
				return
			}
			updateCmd := exec.Command("sh", installScript)
			updateCmd.Dir = binDir
			updateCmd.Stdout = os.Stdout
			updateCmd.Stderr = os.Stderr
			cmd.Println("Running installer...")
			if err := updateCmd.Run(); err != nil {
				cmd.Println("Failed to run install.sh. Try running with sudo or install manually.")
				return
			}
		}
		cmd.Println("✅ Updated to version", release.TagName)
		cmd.Println("Please restart DevShare.")
	},
}

// extractFromZip extracts the specified binary from a zip archive.
// It returns an error if the binary is not found or if any other issue occurs.
func extractFromZip(zipPath, binName string, outFile *os.File) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		if f.Name == binName || strings.HasSuffix(f.Name, "/"+binName) {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()
			_, err = io.Copy(outFile, rc)
			return err
		}
	}
	return fmt.Errorf("binary %s not found in zip", binName)
}

// extractFromTarGz extracts the specified binary from a tar.gz archive.
// It returns an error if the binary is not found or if any other issue occurs.
func extractFromTarGz(tarGzPath, binName string, outFile *os.File) error {
	f, err := os.Open(tarGzPath)
	if err != nil {
		return err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()
	tarReader := tar.NewReader(gz)
	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if hdr.Typeflag == tar.TypeReg && (hdr.Name == binName || strings.HasSuffix(hdr.Name, "/"+binName)) {
			_, err = io.Copy(outFile, tarReader)
			return err
		}
	}
	return fmt.Errorf("binary %s not found in tar.gz", binName)
}
