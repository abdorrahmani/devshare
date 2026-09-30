package detector

import (
	"os"
	"path/filepath"
	"testing"
)

func write(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestDetectGo(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "go.mod", "module x\n")
	if pt, _ := DetectProjectType(dir); pt != "go" {
		t.Fatalf("want go, got %q", pt)
	}
}

func TestDetectLaravel(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "artisan", "#!/usr/bin/env php\n")
	if pt, _ := DetectProjectType(dir); pt != "laravel" {
		t.Fatalf("want laravel, got %q", pt)
	}
}

func TestDetectReactNpm(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "package.json", `{"dependencies":{"react":"^18"}}`)
	write(t, dir, "package-lock.json", "{}")
	if pt, pm := DetectProjectType(dir); pt != "react" || pm != "npm" {
		t.Fatalf("want react/npm, got %q/%q", pt, pm)
	}
}

func TestDetectNextJSYarn(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "next.config.js", "module.exports = {}\n")
	write(t, dir, "package.json", `{"dependencies":{"next":"14"}}`)
	write(t, dir, "yarn.lock", "")
	if pt, pm := DetectProjectType(dir); pt != "nextjs" || pm != "yarn" {
		t.Fatalf("want nextjs/yarn, got %q/%q", pt, pm)
	}
}

// Vue projects also have package.json + a lockfile, so detection must run the
// Vue check before the generic Node.js check or Vue is unreachable.
func TestDetectVuePnpm(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "package.json", `{"dependencies":{"vue":"^3"}}`)
	write(t, dir, "pnpm-lock.yaml", "")
	if pt, pm := DetectProjectType(dir); pt != "vue" || pm != "pnpm" {
		t.Fatalf("want vue/pnpm, got %q/%q", pt, pm)
	}
}

func TestDetectNodeJS(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "package.json", `{"dependencies":{"express":"^4"}}`)
	write(t, dir, "package-lock.json", "{}")
	if pt, pm := DetectProjectType(dir); pt != "nodejs" || pm != "npm" {
		t.Fatalf("want nodejs/npm, got %q/%q", pt, pm)
	}
}

func TestDetectNone(t *testing.T) {
	dir := t.TempDir()
	if pt, _ := DetectProjectType(dir); pt != "" {
		t.Fatalf("want empty, got %q", pt)
	}
}
