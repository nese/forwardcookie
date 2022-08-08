// Package forwardcookie this is the package.
package forward_cookie

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/textproto"
	"strings"
)

// Config the plugin configuration.
type Config struct {
	Addr       string   `json:"addr,omitempty"`
	Cookies    []string `json:"cookies,omitempty"`
	Headers    []string `json:"headers,omitempty"`
	Parameters []string `json:"parameters,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Cookies:    make([]string, 0),
		Headers:    make([]string, 0),
		Parameters: make([]string, 0),
	}
}

// Cookies holder.
type Cookies struct {
	cookieHeaders map[string]string
}

// ForwardCookie a ForwardCookie plugin.
type ForwardCookie struct {
	next       http.Handler
	addr       string
	name       string
	cookies    []string
	headers    []string
	parameters []string
}

// New created a new ForwardCookie plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &ForwardCookie{
		next:       next,
		name:       name,
		addr:       config.Addr,
		cookies:    config.Cookies,
		headers:    config.Headers,
		parameters: config.Parameters,
	}, nil
}

func (e *ForwardCookie) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	isFetchReq := true
	// if the incoming request contains cookies
	if len(req.Cookies()) > 0 {
		// make sure that those cookies won't be added two-fold
		for _, name := range e.cookies {
			cookie, err := req.Cookie(name)
			if err != nil {
				isFetchReq = false
			}
			if cookie != nil {
				isFetchReq = false
			}
		}
	}
	if isFetchReq {
		fetchReq, err := http.NewRequest(http.MethodGet, e.addr, nil)
		if err != nil {
			log.Fatalf(fmt.Sprintf("%s", err))
			return
		}

		AddHeaders(fetchReq, req, e)
		AddParameters(fetchReq, req, e)

		client := http.Client{}
		forwardResponse, err := client.Do(fetchReq)

		if err != nil {
			log.Fatalf(fmt.Sprintf("%s", err))
			return
		}

		cookies := Cookies{
			cookieHeaders: make(map[string]string),
		}

		responseCookies := forwardResponse.Header["Set-Cookie"]

		for _, wantedCookie := range e.cookies {
			for _, header := range responseCookies {
				cookieName := getCookieName(header)
				if cookieName == wantedCookie {
					cookies.cookieHeaders[cookieName] = header
				}
			}
		}

		for _, cookieHeader := range cookies.cookieHeaders {
			rw.Header().Add("Set-Cookie", cookieHeader)
		}
	}

	e.next.ServeHTTP(rw, req)
}

// Extract cookie name
func getCookieName(setCookieString string) string {
	parts := strings.Split(textproto.TrimString(setCookieString), ";")
	if len(parts) == 1 && parts[0] == "" {
		return ""
	}
	cookiePair := textproto.TrimString(parts[0])
	j := strings.Index(cookiePair, "=")
	if j < 0 {
		return ""
	}
	cookieName := cookiePair[:j]
	return cookieName
}

// AddHeaders to fetchReq.
func AddHeaders(fetchReq, req *http.Request, config *ForwardCookie) {
	for _, wantedHeader := range config.headers {
		value := req.Header.Get(wantedHeader)
		if value != "" {
			fetchReq.Header.Add(wantedHeader, value)
		}
	}
}

// AddParameters to fetchReq.
func AddParameters(fetchReq, req *http.Request, config *ForwardCookie) {
	for _, wantedParam := range config.parameters {
		value := req.URL.Query().Get(wantedParam)
		if value != "" {
			q := req.URL.Query()
			q.Add(wantedParam, value)
			fetchReq.URL.RawQuery = q.Encode()
		}
	}
}
