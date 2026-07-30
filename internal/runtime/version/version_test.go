package version_test

import (
	"testing"

	"github.com/saaedimam/jervis/internal/runtime/version"
)

func TestVersionInfo(t *testing.T) {
	info := version.NewInfo("1.2.3", "commit123", "2026-07-29", "darwin/arm64", "go1.22.0")

	if info.SemVer() != "1.2.3" {
		t.Errorf("expected 1.2.3, got %s", info.SemVer())
	}
	if info.GitCommit() != "commit123" {
		t.Errorf("expected commit123, got %s", info.GitCommit())
	}
	if info.BuildDate() != "2026-07-29" {
		t.Errorf("expected 2026-07-29, got %s", info.BuildDate())
	}
	if info.BuildTarget() != "darwin/arm64" {
		t.Errorf("expected darwin/arm64, got %s", info.BuildTarget())
	}
	if info.GoVersion() != "go1.22.0" {
		t.Errorf("expected go1.22.0, got %s", info.GoVersion())
	}

	expectedString := "SemVer: 1.2.3, Commit: commit123, BuildDate: 2026-07-29, Target: darwin/arm64, GoVersion: go1.22.0"
	if info.String() != expectedString {
		t.Errorf("expected string %q, got %q", expectedString, info.String())
	}
}
