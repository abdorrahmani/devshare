package cmd

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
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
		var release ghRelease
		if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
			cmd.Println("Failed to parse release info:", err)
			return
		}
		current := strings.TrimPrefix(Version, "v")
		latest := strings.TrimPrefix(release.TagName, "v")
		switch compareVersions(current, latest) {
		case 0:
			cmd.Println("You are already running the latest version (", Version, ")!")
			return
		case 1:
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

		// Verify the archive against the release checksums before extracting or
		// executing anything from it.
		checksumsURL := findAssetURL(release.Assets, "checksums.txt")
		if checksumsURL == "" {
			cmd.Println("❌ No checksums file found in the release; cannot verify the download. Aborting for safety.")
			return
		}
		sums, err := httpGetBytes(client, checksumsURL)
		if err != nil {
			cmd.Println("Failed to download checksums:", err)
			return
		}
		if err := verifyChecksum(tmpArchive.Name(), assetName, sums); err != nil {
			cmd.Println("❌ Checksum verification failed:", err)
			return
		}
		cmd.Println("✅ Checksum verified.")

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

type ghAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type ghRelease struct {
	TagName string    `json:"tag_name"`
	Assets  []ghAsset `json:"assets"`
}

// compareVersions compares dotted numeric versions ("1.9.0", "1.10.0").
// Returns -1 if a<b, 0 if equal, 1 if a>b. Any pre-release/build suffix is ignored.
// ponytail: numeric-only; switch to golang.org/x/mod/semver if pre-release ordering ever matters.
func compareVersions(a, b string) int {
	na, nb := parseVersion(a), parseVersion(b)
	for i := 0; i < 3; i++ {
		switch {
		case na[i] < nb[i]:
			return -1
		case na[i] > nb[i]:
			return 1
		}
	}
	return 0
}

func parseVersion(v string) [3]int {
	v = strings.TrimPrefix(v, "v")
	if i := strings.IndexAny(v, "-+"); i >= 0 {
		v = v[:i]
	}
	var out [3]int
	for i, part := range strings.SplitN(v, ".", 3) {
		out[i], _ = strconv.Atoi(strings.TrimSpace(part))
	}
	return out
}

// findAssetURL returns the download URL of the first asset whose name ends with suffix.
func findAssetURL(assets []ghAsset, suffix string) string {
	for _, a := range assets {
		if strings.HasSuffix(strings.ToLower(a.Name), suffix) {
			return a.BrowserDownloadURL
		}
	}
	return ""
}

// verifyChecksum checks archivePath's SHA-256 against the entry for assetName in
// a goreleaser-style checksums file ("<sha256>  <filename>" per line).
func verifyChecksum(archivePath, assetName string, checksums []byte) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	got := hex.EncodeToString(h.Sum(nil))

	sc := bufio.NewScanner(bytes.NewReader(checksums))
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) != 2 {
			continue
		}
		if filepath.Base(fields[1]) == assetName {
			if strings.EqualFold(fields[0], got) {
				return nil
			}
			return fmt.Errorf("sha256 mismatch for %s", assetName)
		}
	}
	return fmt.Errorf("no checksum entry for %s", assetName)
}

func httpGetBytes(client *http.Client, url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status %s", resp.Status)
	}
	return io.ReadAll(resp.Body)
}
