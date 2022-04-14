package dist

import "embed"

//go:embed static
var Static embed.FS

//go:embed index.html
var Index []byte
