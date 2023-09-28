package main

import (
	"log/slog"
	"os"
)

var logger = slog.New(slog.NewTextHandler(os.Stderr, nil))
