package main

import (
	_ "embed" // Import for embedding files
)

//go:embed defaultData/config.toml
var fsDefaultConfig []byte

//go:embed defaultData/start.sh
var fsDefaultStart []byte

//go:generate go install github.com/akavel/rsrc@latest
//go:generate rsrc -manifest assets/manifest.xml -ico assets/icon.ico
