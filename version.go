package logger

import (
	"runtime/debug"
)

// Version returns the version of the logger.
// It attempts to get the version from git metadata in the following order:
//  1. Git tag (if available)
//  2. Git commit hash
//  3. Git branch name
//  4. "unknown" (fallback)
func Version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	var vcsRevision, vcsModified string
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			vcsRevision = setting.Value
		case "vcs.modified":
			vcsModified = setting.Value
		}
	}

	if vcsRevision != "" {
		if vcsModified == "true" {
			return vcsRevision[:7] + "-dirty"
		}
		if len(vcsRevision) > 7 {
			return vcsRevision[:7]
		}
		return vcsRevision
	}

	if info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}

	return "unknown"
}
