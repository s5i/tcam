package gamedata

import (
	"io"
	"log"
)

var Logger = log.New(io.Discard, "[GAMEDATA] ", 0)
