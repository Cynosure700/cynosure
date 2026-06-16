package tools

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTMLToText(t *testing.T) {
	in := `<html><head><style>.a{}</style><script>var x=1;</script></head>` +
		`<body><h1>Title</h1><p>Hello&nbsp;&amp; world</p></body></html>`
	out := htmlToText(in)
	if strings.Contains(out, "var x") || strings.Contains(out, ".a{}") {
		t.Fatalf("script/style not stripped: %q", out)
	}
	if !strings.Contains(out, "Title") || !strings.Contains(out, "Hello") || !strings.Contains(out, "world") {
		t.Fatalf("text not extracted: %q", out)
	}
}

func TestRunWebFetchWithoutProcessor(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><body><p>page body</p></body></html>"))
	}))
	defer srv.Close()

	ctx := context.Background()
	// Use the test server's client to trust its TLS cert by temporarily
	// swapping the default client transport.
	old := http.DefaultClient
	http.DefaultClient = srv.Client()
	defer func() { http.DefaultClient = old }()

	// httptest TLS server URL is https already; pass through.
	out, err := RunWebFetch(ctx, srv.URL, "summarize")
	if err != nil {
		t.Fatalf("web_fetch: %v", err)
	}
	if !strings.Contains(out, "page body") {
		t.Fatalf("expected page body, got %q", out)
	}
}

func TestRunWebFetchWithProcessor(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><body>secret content</body></html>"))
	}))
	defer srv.Close()

	old := http.DefaultClient
	http.DefaultClient = srv.Client()
	defer func() { http.DefaultClient = old }()

	ctx := WithWebProcessor(context.Background(), func(ctx context.Context, prompt, content string) (string, error) {
		if !strings.Contains(content, "secret content") {
			t.Fatalf("processor did not receive cleaned content: %q", content)
		}
		return "PROCESSED:" + prompt, nil
	})

	out, err := RunWebFetch(ctx, srv.URL, "extract")
	if err != nil {
		t.Fatalf("web_fetch: %v", err)
	}
	if out != "PROCESSED:extract" {
		t.Fatalf("expected processed output, got %q", out)
	}
}

func TestRunWebSearchPlaceholder(t *testing.T) {
	out, err := RunWebSearch(context.Background(), "golang")
	if err != nil {
		t.Fatalf("web_search: %v", err)
	}
	if !strings.Contains(out, "not configured") {
		t.Fatalf("expected placeholder message, got %q", out)
	}
}
