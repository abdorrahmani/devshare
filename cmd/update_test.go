package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestCompareVersions(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"1.9.0", "1.10.0", -1}, // string compare gets this backwards ("1.9.0" > "1.10.0")
		{"1.10.0", "1.9.0", 1},
		{"1.1.0", "1.1.0", 0},
		{"v1.2.3", "1.2.3", 0}, // leading v ignored
		{"2.0.0", "1.9.9", 1},
		{"1.0.0", "1.0.1", -1},
		{"1.2.0-rc1", "1.2.0", 0}, // pre-release suffix ignored
	}
	for _, c := range cases {
		if got := compareVersions(c.a, c.b); got != c.want {
			t.Errorf("compareVersions(%q,%q)=%d want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestVerifyChecksum(t *testing.T) {
	dir := t.TempDir()
	name := "DevShare_Linux_x86_64.tar.gz"
	archive := filepath.Join(dir, name)
	data := []byte("pretend archive bytes")
	if err := os.WriteFile(archive, data, 0644); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(data)
	good := fmt.Sprintf("%s  %s\n", hex.EncodeToString(sum[:]), name)

	if err := verifyChecksum(archive, name, []byte(good)); err != nil {
		t.Fatalf("valid checksum rejected: %v", err)
	}
	bad := "0000000000000000000000000000000000000000000000000000000000000000  " + name + "\n"
	if err := verifyChecksum(archive, name, []byte(bad)); err == nil {
		t.Fatal("mismatched checksum accepted")
	}
	if err := verifyChecksum(archive, "other.tar.gz", []byte(good)); err == nil {
		t.Fatal("missing checksum entry accepted")
	}
}
