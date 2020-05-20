package goose

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestContextQuery(t *testing.T) {
	params := url.Values{}
	params.Set("name", "Zheng Zhou")
	params.Set("email", "go.ose@goose.com")
	req := httptest.NewRequest(http.MethodGet, "/?"+params.Encode(), nil)
	g := New()
	g.context.resetContext(nil, req)
	ctx := g.context

	want := "Zheng Zhou"
	if got := ctx.Query("name"); got != want {
		t.Errorf("Want: %s , Got: %s", want, got)
	}

	want = "go.ose@goose.com"
	if got := ctx.Query("email"); got != want {
		t.Errorf("Want: %s , Got: %s", want, got)
	}
}
