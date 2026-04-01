package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveFontPathUsesConfigDirForNamedFonts(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	fontDir := filepath.Join(home, ".config", "go-figlet")
	if err := os.MkdirAll(fontDir, 0o755); err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}
	fontPath := filepath.Join(fontDir, "efont-b24-bloody-double.flf")
	if err := os.WriteFile(fontPath, []byte("test"), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	got, err := resolveFontPath("efont-b24-bloody-double")
	if err != nil {
		t.Fatalf("resolveFontPath failed: %v", err)
	}
	if got != fontPath {
		t.Fatalf("path = %q, want %q", got, fontPath)
	}
}

func TestResolveFontPathPreservesExplicitPaths(t *testing.T) {
	got, err := resolveFontPath("./fonts/custom.flf")
	if err != nil {
		t.Fatalf("resolveFontPath failed: %v", err)
	}
	if got != "./fonts/custom.flf" {
		t.Fatalf("path = %q, want %q", got, "./fonts/custom.flf")
	}
}

func TestListNamedFonts(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	fontDir := filepath.Join(home, ".config", "go-figlet")
	if err := os.MkdirAll(fontDir, 0o755); err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}
	for _, name := range []string{"zeta.flf", "alpha.flf", "note.txt"} {
		if err := os.WriteFile(filepath.Join(fontDir, name), []byte("test"), 0o644); err != nil {
			t.Fatalf("WriteFile failed: %v", err)
		}
	}

	got, err := listNamedFonts()
	if err != nil {
		t.Fatalf("listNamedFonts failed: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0] != "alpha" || got[1] != "zeta" {
		t.Fatalf("got %v, want [alpha zeta]", got)
	}
}

func TestListColors(t *testing.T) {
	got := listColors()
	if len(got) == 0 {
		t.Fatal("expected non-empty color list")
	}
	if got[0] != "preset:cool" {
		t.Fatalf("first entry = %q, want %q", got[0], "preset:cool")
	}
	if got[len(got)-1] != "name:yellowgreen" {
		t.Fatalf("last entry = %q, want %q", got[len(got)-1], "name:yellowgreen")
	}
}
