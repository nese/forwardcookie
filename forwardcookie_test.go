package forward_cookie_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	forwardcookie "github.com/nese/forward-cookie"
)

func TestForwardCookie(t *testing.T) {
	cfg := forwardcookie.CreateConfig()
	cfg.Addr = "https://my-domain.com/"
	cfg.Cookies = []string{"jsessionid", "common-web-stuff"}
	cfg.Headers = []string{"X-Forwarded-For"}
	cfg.Parameters = []string{"q-param"}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	plugin, err := forwardcookie.New(ctx, next, cfg, "forwardcookie-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cfg.Addr, nil)
	if err != nil {
		t.Fatal(err)
	}
	plugin.ServeHTTP(recorder, req)
}
