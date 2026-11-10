package models

import (
	"regexp"
	"strings"
)

var nonSlugChars = regexp.MustCompile(`[^a-z0-9]+`)

// Slugify mengubah judul menjadi slug URL yang aman.
// Contoh: "Sistem Informasi Pegawai" -> "sistem-informasi-pegawai"
func Slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = nonSlugChars.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "untitled"
	}
	return s
}
