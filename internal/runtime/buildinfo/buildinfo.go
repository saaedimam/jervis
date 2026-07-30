package buildinfo

import (
	"runtime"

	"github.com/saaedimam/jervis/internal/runtime/version"
)

// Private unexported variables for linker (-ldflags) injection at compile time.
var (
	semVer      = "0.1.0-dev"
	gitCommit   = "unknown"
	buildDate   = "unknown"
	buildTarget = runtime.GOOS + "/" + runtime.GOARCH
)

// Get returns the immutable build version metadata.
func Get() version.Info {
	return version.NewInfo(
		semVer,
		gitCommit,
		buildDate,
		buildTarget,
		runtime.Version(),
	)
}
