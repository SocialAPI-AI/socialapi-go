package socialapi

// Hand-written: survives regeneration.

import (
	"context"
	"net/http"
)

// DefaultServer is the production API base URL.
const DefaultServer = "https://api.social-api.ai/v1"

// WithAPIKey sends your API key as a bearer token on every request.
//
//	client, err := socialapi.NewClientWithResponses(socialapi.DefaultServer, socialapi.WithAPIKey("sapi_key_..."))
func WithAPIKey(key string) ClientOption {
	return WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+key)
		return nil
	})
}
