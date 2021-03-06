package mtree

import "fmt"

const (
	// AppName is the name ... of this library/application
	AppName = "gomtree"
)

const (
	// VersionMajor is for an API incompatible changes
	VersionMajor = 0
	// VersionMinor is for functionality in a backwards-compatible manner
	VersionMinor = 4
	// VersionPatch is for backwards-compatible bug fixes
	VersionPatch = 1

	// VersionDev indicates development branch. Releases will be empty string.
	VersionDev = ""
)

// Version is the specification version that the package types support.
var Version = fmt.Sprintf("%d.%d.%d%s", VersionMajor, VersionMinor, VersionPatch, VersionDev)
