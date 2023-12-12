package url

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testUrl = "https://www.villaleobos.com/"
)

func TestShorten(t *testing.T) {
	ctx := context.Background()
	params := ShortenParams{Url: testUrl}
	url, err := Shorten(ctx, &params)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, url)
	assert.NotEmpty(t, url.Id)
}

func TestRemove(t *testing.T) {
	ctx := context.Background()
	params := ShortenParams{Url: testUrl}
	// Add a URL to the database
	url, err := Shorten(ctx, &params)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, url)
	// Remove the same URL
	err = Remove(ctx, url.Id)
	assert.Nil(t, err)
}

func TestFetch(t *testing.T) {
	ctx := context.Background()
	params := ShortenParams{Url: testUrl}
	// Add a URL to the database
	url, err := Shorten(ctx, &params)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, url)
	// Fetch the same URL
	target := fmt.Sprintf("/url/%s", url.Id)
	req := httptest.NewRequest(http.MethodGet, target, nil)
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(FetchInternal)
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusPermanentRedirect, w.Result().StatusCode)
	location := w.Result().Header.Get("Location")
	assert.Equal(t, url.Url, location)
}
