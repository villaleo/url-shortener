package url

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"encore.dev/beta/errs"
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
//encore:api public raw method=GET path=/url/:id
func Fetch(w http.ResponseWriter, req *http.Request) {
	FetchInternal(w, req)
}

// Internal implementation for Fetch. Exposed for testing purposes.
func FetchInternal(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := req.URL.Path[len("/url/"):]
	url, err := findUrlById(req.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		body := errs.Error{
			Code:    errs.NotFound,
			Message: "No URL was found for the specified ID. Check the ID and try again.",
		}
		encodedBody, _ := json.Marshal(body)
		w.Write(encodedBody)
		return
	}
	// url will never be nil
	http.Redirect(w, req, *url, http.StatusPermanentRedirect)
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

// Removes the URL from the database with the specified ID.
//
//encore:api public method=DELETE path=/url/:id
func Remove(ctx context.Context, id string) error {
	if _, err := findUrlById(ctx, id); err != nil {
		return &errs.Error{
			Code:    errs.NotFound,
			Message: "No URL was found for the specified ID. Check the ID and try again.",
		}
	}
	if _, err := db.Query(ctx, "DELETE FROM url WHERE id = $1", id); err != nil {
		return &errs.Error{
			Code:    errs.Internal,
			Message: "Failed to remove URL.",
			Details: errs.Details(err),
		}
	}
	return nil
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

func findUrlById(ctx context.Context, id string) (url *string, err error) {
	err = db.QueryRow(ctx, "SELECT original_url FROM url WHERE id = $1", id).Scan(&url)
	return
}
