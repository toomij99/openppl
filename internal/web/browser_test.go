package web

import "testing"

func TestBrowserURLLocalhostMapping(t *testing.T) {
	if got := browserURL("0.0.0.0", 5016); got != "http://localhost:5016" {
		t.Fatalf("expected localhost mapping, got %q", got)
	}

	if got := browserURL("::", 5016); got != "http://localhost:5016" {
		t.Fatalf("expected localhost mapping for ipv6 wildcard, got %q", got)
	}

	if got := browserURL("[::]", 5016); got != "http://localhost:5016" {
		t.Fatalf("expected localhost mapping for bracketed ipv6 wildcard, got %q", got)
	}

	if got := browserURL("127.0.0.1", 5016); got != "http://127.0.0.1:5016" {
		t.Fatalf("expected host passthrough, got %q", got)
	}
}
