# SocialAPI Go SDK

Official Go client for the [SocialAPI](https://social-api.ai) REST API.

```bash
go get github.com/SocialAPI-AI/socialapi-go
```

```go
client, err := socialapi.NewClientWithResponses(socialapi.DefaultServer, socialapi.WithAPIKey("sapi_key_..."))
if err != nil {
	log.Fatal(err)
}
resp, err := client.ListBrandsWithResponse(ctx)
```

Each response has a typed field per status code (`resp.JSON200`, `resp.JSON404`, ...). The API reference is at [docs.social-api.ai](https://docs.social-api.ai).

## Verifying webhooks

```go
body, _ := io.ReadAll(r.Body)
if err := socialapi.VerifyWebhook("whsec_...", body, r.Header, 0); err != nil {
	http.Error(w, "invalid signature", http.StatusUnauthorized)
	return
}
```

Pass the raw request body. A tolerance of 0 rejects deliveries older than 5 minutes.

## How this package is built

`socialapi.gen.go` is generated from the API's OpenAPI spec by [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen). Do not edit it: every release regenerates it. The other `.go` files are hand-written. Releases are fully automated by `.github/workflows/release.yml`, and versions follow semver from the spec diff.
