package services

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"strings"
)

func ShortenUrl(ctx context.Context, url string) string {
	hasher := md5.New()
	hasher.Write([]byte(url))
	hash := hasher.Sum(nil)

	encoded := base64.StdEncoding.EncodeToString(hash)
	cleaned := strings.ReplaceAll(strings.ReplaceAll(encoded, "=", ""), "/", "_")
	return cleaned[0:8]
}
