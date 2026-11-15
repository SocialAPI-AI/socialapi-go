package socialapi

import (
	"net/http"
	"testing"
	"time"
)

// Reference vector shared by every SDK: HMAC-SHA256("whsec_test", "1700000000." + body).
const (
	testSecret = "whsec_test"
	testTS     = "1700000000"
	testBody   = `{"event":"comment.received"}`
	testSig    = "sha256=a1c084394bcdc1743e02d34447378108cde7b6ac0c963bd211e97ba76bb3f3bf"
)

func TestVerifyWebhook(t *testing.T) {
	h := http.Header{}
	h.Set("X-SocialAPI-Timestamp", testTS)
	h.Set("X-SocialAPI-Signature-V2", testSig)
	now := time.Unix(1700000000, 0)
	body := []byte(testBody)

	if err := verifyWebhookAt(testSecret, body, h, 0, now.Add(time.Minute)); err != nil {
		t.Fatalf("valid delivery rejected: %v", err)
	}
	if verifyWebhookAt("wrong", body, h, 0, now) == nil {
		t.Fatal("wrong secret accepted")
	}
	if verifyWebhookAt(testSecret, []byte(`{}`), h, 0, now) == nil {
		t.Fatal("tampered body accepted")
	}
	if verifyWebhookAt(testSecret, body, h, 0, now.Add(6*time.Minute)) == nil {
		t.Fatal("stale timestamp accepted")
	}
	h.Del("X-SocialAPI-Timestamp")
	if verifyWebhookAt(testSecret, body, h, 0, now) == nil {
		t.Fatal("missing timestamp accepted")
	}
}
