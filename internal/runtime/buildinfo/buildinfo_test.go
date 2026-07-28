package buildinfo_test

import (
	"runtime"
	"testing"

	"github.com/ioriimasu/jervis/internal/runtime/buildinfo"
)

func TestBuildInfoGet(t *testing.T) {
	info := buildinfo.Get()

	if info.SemVer() == "" {
		t.Errorf("expected non-empty SemVer")
	}
	if info.GitCommit() == "" {
		t.Errorf("expected non-empty GitCommit")
	}
	if info.BuildDate() == "" {
		t.Errorf("expected non-empty BuildDate")
	}
	if info.BuildTarget() == "" {
		t.Errorf("expected non-empty BuildTarget")
	}
	if info.GoVersion() != runtime.Version() {
		t.Errorf("expected %s, got %s", runtime.Version(), info.GoVersion())
	}
}
