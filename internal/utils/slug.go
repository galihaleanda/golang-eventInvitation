package utils

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gosimple/slug"
)

func GenerateSlug(title string) string {
	base := slug.Make(title)
	if base == "" {
		base = "event"
	}
	return fmt.Sprintf("%s-%s", base, randomSuffix(6))
}

func randomSuffix(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}
