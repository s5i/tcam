package parser

import (
	"io"
	"log"
)

var Logger = log.New(io.Discard, "[PARSER] ", 0)
