package url

import (
	"context"
	"crypto/rand"
	"encoding/base64"
)

type URL struct {
	Id  string
	Url string
}

type ShortenedParams struct {
	Url string
}

//encore:api public method=POST path=/url
func Shorten(ctx context.Context, params *ShortenedParams) (*URL, error) {
	id, err := generateId()
	if err != nil {
		return nil, err
	}
	return &URL{Id: id, Url: params.Url}, nil
}

func generateId() (string, error) {
	var data [6]byte
	if _, err := rand.Read(data[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data[:]), nil
}
