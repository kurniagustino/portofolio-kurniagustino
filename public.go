package main

import (
	"embed"
)

//go:embed public/*
var publicFS embed.FS
