package params

import (
    "fmt"
)

const (
    VersionMajor = 0          // Major version component of the current release
    VersionMinor = 0          // Minor version component of the current release
    VersionPatch = 0         // Patch version component of the current release
)

// Version holds the textual version string.
var ChainVersion = func() string {
    return fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
}()
