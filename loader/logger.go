package loader

import (
	"io"
	"log"
)

var Logger = log.New(io.Discard, "[LOADER] ", 0)
