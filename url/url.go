package url

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"encore.dev/storage/sqldb"
)

type URL struct {
	// An identifier generated for the Url.
	Id string
	// The Url to shorten.
	Url string
}

// Fetches an existing URL by the associated Id.
//
//encore:api public method=GET path=/url/:id
func Fetch(ctx context.Context, id string) (*URL, error) {
	u := &URL{Id: id}
	err := db.QueryRow(ctx, `
		SELECT original_url FROM url
		WHERE id = $1
	`, id).Scan(&u.Url)
	return u, err
}

type ShortenParams struct {
	// The Url to shorten.
	Url string
}

// Generates an Id for the URL and saves the mapping to a database.
//
//encore:api public method=POST path=/url
func Shorten(ctx context.Context, params *ShortenParams) (*URL, error) {
	id, err := generateId()
	if err != nil {
		return nil, err
	} else if err := insert(ctx, id, params.Url); err != nil {
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

var db = sqldb.NewDatabase("url", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})

func insert(ctx context.Context, id, url string) error {
	_, err := db.Exec(ctx, `
		INSERT INTO url (id, original_url)
		VALUES ($1, $2)
	`, id, url)
	return err
}
