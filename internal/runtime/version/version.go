package version

import "fmt"

// Info represents immutable runtime version details.
type Info struct {
	semVer      string
	gitCommit   string
	buildDate   string
	buildTarget string
	goVersion   string
}

// NewInfo constructs an immutable Info container.
func NewInfo(semVer, gitCommit, buildDate, buildTarget, goVersion string) Info {
	return Info{
		semVer:      semVer,
		gitCommit:   gitCommit,
		buildDate:   buildDate,
		buildTarget: buildTarget,
		goVersion:   goVersion,
	}
}

// SemVer returns the semantic version string (e.g. "1.0.0").
func (i Info) SemVer() string {
	return i.semVer
}

// GitCommit returns the Git commit hash string.
func (i Info) GitCommit() string {
	return i.gitCommit
}

// BuildDate returns the build date string.
func (i Info) BuildDate() string {
	return i.buildDate
}

// BuildTarget returns the target platform string (e.g. "darwin/arm64").
func (i Info) BuildTarget() string {
	return i.buildTarget
}

// GoVersion returns the Go runtime compiler version.
func (i Info) GoVersion() string {
	return i.goVersion
}

// String produces a formatted string representation of the version info.
func (i Info) String() string {
	return fmt.Sprintf("SemVer: %s, Commit: %s, BuildDate: %s, Target: %s, GoVersion: %s",
		i.semVer, i.gitCommit, i.buildDate, i.buildTarget, i.goVersion)
}
