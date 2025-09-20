package gui

import "akashic_scribe/core"

// ScribeOptions is now defined in the core package to avoid import cycles.
// We create a type alias here for backward compatibility within the GUI package.
type ScribeOptions = core.ScribeOptions
