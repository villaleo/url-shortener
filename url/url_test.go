package url

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortenAndFetch(t *testing.T) {
	const testUrl = "https://github.com/encoredev/encore"
	ctx := context.Background()

	params := ShortenParams{Url: testUrl}
	shortenResp, err := Shorten(ctx, &params)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, shortenResp.Url, testUrl)

	fetchResp, err := Fetch(ctx, shortenResp.Id)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, shortenResp.Url, fetchResp.Url)
}
