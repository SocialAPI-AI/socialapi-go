package socialapi

// Release gate for the generated client: base URL, bearer auth, decoding,
// error mapping, and tolerance of values the spec does not list yet.

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testClient(t *testing.T, status int, body string) *ClientWithResponses {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer sapi_key_test" {
			t.Errorf("Authorization = %q", r.Header.Get("Authorization"))
		}
		if r.URL.Path != "/v1/posts/p_1" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	c, err := NewClientWithResponses(srv.URL+"/v1", WithAPIKey("sapi_key_test"))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestForwardCompatibleDecode(t *testing.T) {
	// An unlisted platform, an unknown field, and a null must all decode.
	c := testClient(t, 200, `{"id":"p_1","status":"published","text":null,"hidden":false,"retry_count":0,"queue_id":null,
		"created_at":"2026-03-14T09:00:00Z","updated_at":"2026-03-14T09:00:00Z","field_from_the_future":1,
		"targets":[{"platform":"platform_from_the_future","status":"published"}]}`)
	resp, err := c.GetPostWithResponse(context.Background(), "p_1")
	if err != nil {
		t.Fatal(err)
	}
	if resp.JSON200 == nil || resp.JSON200.Id != "p_1" {
		t.Fatalf("not decoded: %s", resp.Body)
	}
	if got := string((*resp.JSON200.Targets)[0].Platform); got != "platform_from_the_future" {
		t.Fatalf("platform = %q", got)
	}
}

func TestErrorResponse(t *testing.T) {
	c := testClient(t, 404, `{"error":{"code":"post.not_found","message":"Post not found"},"request_id":"req_1"}`)
	resp, err := c.GetPostWithResponse(context.Background(), "p_1")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode() != 404 || resp.JSON404 == nil {
		t.Fatalf("404 not mapped: %d %s", resp.StatusCode(), resp.Body)
	}
}
