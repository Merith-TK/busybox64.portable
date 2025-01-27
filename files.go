package main

import (
	_ "embed" // Import for embedding files
)

//go:embed defaultData/config.toml
var fsDefaultConfig []byte

//go:embed defaultData/start.sh
var fsDefaultStart []byte
